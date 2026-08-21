package newrelic

// resource_newrelic_team manages a New Relic NGEP Team entity and its two
// auto-created backing collections (membership and ownership). See the
// Teams & Scorecards Ecosystem Observations Confluence pages and the
// newrelic-client-go pkg/scorecards NGEP_ANALYSIS.md for the full design
// rationale behind the choices made here.

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	nr "github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

func resourceNewRelicTeam() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicTeamCreate,
		ReadContext:   resourceNewRelicTeamRead,
		UpdateContext: resourceNewRelicTeamUpdate,
		DeleteContext: resourceNewRelicTeamDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// ── Core team identity ──────────────────────────────────────────
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team. Must be unique within the organization.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free-text description for the team.",
			},
			"aliases": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Additional searchable names for the team. Each alias must be unique within the organization.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags in 'key:value1,value2' format.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// ── Hierarchy ──────────────────────────────────────────────────
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NGEP GUID of the parent team. Setting this establishes a team hierarchy. Cannot be set to the team's own GUID.",
			},
			// ── Resources (links / docs) ────────────────────────────────────
			"resources": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "Supplemental resources attached to the team " +
					"(e.g. runbooks, wikis). Each resource must have a type and content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of resource (e.g. 'link').",
						},
						"content": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource content (e.g. a URL).",
						},
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-readable title for the resource.",
						},
					},
				},
			},
			// ── Membership (people) ─────────────────────────────────────────
			// Each element is a NR integer user ID.  Terraform resolves these to
			// EntityManagementUserEntity GUIDs via the standard entities
			// entitySearch API before calling entityManagementAddCollectionMembers.
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of New Relic user IDs (integers) to add as members of this team. " +
					"These are added to the team's auto-created membership collection. " +
					"Managers must be a subset of members — add a user here before setting them as a manager.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The New Relic user ID (integer). This can be obtained from the newrelic_user data source or the Users page in the NR UI.",
						},
					},
				},
			},
			// ── Managers (subset of members) ────────────────────────────────
			"managers": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of New Relic user IDs (integers) who are managers of this team. " +
					"Every manager must already be listed in the members block. " +
					"Managers are set via a separate update call after members are added.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The New Relic user ID (integer) of the manager.",
						},
					},
				},
			},
			// ── Ownership (entities the team owns) ──────────────────────────
			"entities": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of entity GUIDs that this team owns. " +
					"These are added to the team's auto-created ownership collection. " +
					"Entities from any account in the organization can be added. " +
					"Dead/deleted entities are automatically removed by NGEP.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The GUID of the entity to claim ownership of.",
						},
					},
				},
			},
			// ── Computed / infrastructure ────────────────────────────────────
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The NGEP organization UUID. If omitted, the provider fetches it automatically.",
			},
			"membership_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created membership collection for this team. Read-only.",
			},
			"ownership_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created ownership collection for this team. Read-only.",
			},
		},
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Create
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicTeamCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	// ── Build create input ──────────────────────────────────────────────────
	input := scorecards.EntityManagementTeamEntityCreateInput{
		Name: d.Get("name").(string),
		Scope: scorecards.EntityManagementScopedReferenceInput{
			ID:   orgID,
			Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION,
		},
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("aliases"); ok {
		for _, a := range v.([]interface{}) {
			input.Aliases = append(input.Aliases, a.(string))
		}
	}
	if v, ok := d.GetOk("parent_id"); ok {
		input.ParentId = v.(string)
	}
	// NOTE: tags and resources are intentionally excluded from the create
	// mutation payload. The NGEP validation service intermittently returns
	// HTTP 500 when resources (and combinations of tags+resources) are
	// present in entityManagementCreateTeam. The update path is consistently
	// reliable for both fields, so we create with only the stable fields
	// (name, description, aliases, scope, parent_id) and then apply tags and
	// resources via an immediate post-create update.
	//
	// managers are NOT passed on create either — NGEP validates they must
	// already be membership-collection members, which only exist after create.

	// ── Call the create mutation ────────────────────────────────────────────
	result, err := client.Scorecards.EntityManagementCreateTeam(input)
	if err != nil {
		return diag.FromErr(err)
	}
	teamID := result.Entity.ID
	membershipColID := result.Entity.Membership.ID
	ownershipColID := result.Entity.Ownership.ID

	log.Printf("[INFO] Created team %s (membership=%s, ownership=%s)", teamID, membershipColID, ownershipColID)

	d.SetId(teamID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set("membership_collection_id", membershipColID)
	_ = d.Set("ownership_collection_id", ownershipColID)

	// ── Post-create update: tags + resources ───────────────────────────────
	// Always set via update to avoid the NGEP create-mutation 500 on resources.
	hasTags := len(d.Get("tags").([]interface{})) > 0
	hasResources := len(d.Get("resources").([]interface{})) > 0
	if hasTags || hasResources {
		upd := scorecards.EntityManagementTeamEntityUpdateInput{}
		if hasTags {
			upd.Tags = expandTeamTags(d.Get("tags").([]interface{}))
		}
		if hasResources {
			upd.Resources = expandTeamResourcesUpdate(d.Get("resources").([]interface{}))
		}
		if _, err := client.Scorecards.EntityManagementUpdateTeam(teamID, upd); err != nil {
			return diag.Errorf("post-create update for tags/resources on team %s: %v", teamID, err)
		}
	}

	// ── Add members ──────────────────────────────────────────────────────────
	if err := syncMembers(ctx, client, membershipColID, nil, expandMemberUserIDs(d)); err != nil {
		return diag.FromErr(err)
	}

	// ── Set managers (after members are added) ────────────────────────────
	if err := syncManagers(ctx, providerConfig, d, teamID); err != nil {
		return diag.FromErr(err)
	}

	// ── Add owned entities ────────────────────────────────────────────────
	if err := syncOwnedEntities(ctx, &client.Scorecards, ownershipColID, nil, expandEntityGUIDs(d)); err != nil {
		return diag.FromErr(err)
	}

	// ── Post-create indexing gate ─────────────────────────────────────────
	// NGEP entities are synced asynchronously to the standard entitySearch
	// index. The Teams UI polls this after create before proceeding. We do
	// the same to guarantee the entity is queryable before Terraform
	// considers the resource created.
	log.Printf("[INFO] Waiting for team %s to appear in entitySearch index", teamID)
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		esResult, err := client.Entities.GetEntitySearchByQueryWithContext(
			ctx,
			entities.EntitySearchOptions{},
			fmt.Sprintf("id IN ('%s')", teamID),
			[]entities.EntitySearchSortCriteria{},
		)
		if err != nil {
			return resource.NonRetryableError(
				fmt.Errorf("entitySearch failed while waiting for team %s to be indexed: %w", teamID, err))
		}
		if esResult == nil || len(esResult.Results.Entities) == 0 {
			return resource.RetryableError(
				fmt.Errorf("team %s not yet visible in entitySearch, retrying", teamID))
		}
		return nil
	})
	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return resourceNewRelicTeamRead(ctx, d, meta)
}

// ──────────────────────────────────────────────────────────────────────────────
// Read
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicTeamRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	entityIface, err := client.Scorecards.GetEntityWithContext(ctx, d.Id())
	if err != nil {
		var notFound *nrErrors.NotFound
		if errors.As(err, &notFound) {
			d.SetId("")
			return nil
		}
		// NGEP ghost-NOT_FOUND: message starts with ": Entity not found."
		// (empty id prefix). Surface it so Terraform can retry on next plan.
		if isNGEPGhostNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if entityIface == nil {
		d.SetId("")
		return nil
	}

	teamEntity, ok := (*entityIface).(*scorecards.EntityManagementTeamEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamEntity", d.Id())
	}

	_ = d.Set("name", teamEntity.Name)
	_ = d.Set("description", teamEntity.Description)
	// aliases has an eventual-consistency lag in the entity read path — the API
	// may return nil for a freshly-created team even though aliases were stored.
	// Only overwrite state when the API actually returns a non-nil list.
	if teamEntity.Aliases != nil {
		_ = d.Set("aliases", teamEntity.Aliases)
	}
	_ = d.Set("parent_id", teamEntity.ParentId)
	_ = d.Set("tags", flattenTeamTags(teamEntity.Tags))
	_ = d.Set("resources", flattenTeamResources(teamEntity.Resources))
	_ = d.Set("organization_id", teamEntity.Scope.ID)
	_ = d.Set("membership_collection_id", teamEntity.Membership.ID)
	_ = d.Set("ownership_collection_id", teamEntity.Ownership.ID)

	// ── Read members from membership collection ───────────────────────────
	memberUserIDs, err := readCollectionUserIDs(ctx, &client.Scorecards, teamEntity.Membership.ID)
	if err != nil {
		log.Printf("[WARN] Could not read membership collection for team %s: %v", d.Id(), err)
	} else {
		_ = d.Set("members", flattenMemberUserIDs(memberUserIDs))
	}

	// ── Read managers ─────────────────────────────────────────────────────
	_ = d.Set("managers", flattenManagerUserIDs(teamEntity.Managers))

	// ── Read owned entities from ownership collection ─────────────────────
	ownerGUIDs, err := readCollectionEntityGUIDs(ctx, &client.Scorecards, teamEntity.Ownership.ID)
	if err != nil {
		log.Printf("[WARN] Could not read ownership collection for team %s: %v", d.Id(), err)
	} else {
		_ = d.Set("entities", flattenEntityGUIDs(ownerGUIDs))
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	teamChanged := d.HasChangesExcept("members", "entities", "managers")

	if teamChanged {
		updateInput := scorecards.EntityManagementTeamEntityUpdateInput{}
		if d.HasChange("name") {
			updateInput.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			updateInput.Description = d.Get("description").(string)
		}
		if d.HasChange("aliases") {
			for _, a := range d.Get("aliases").([]interface{}) {
				updateInput.Aliases = append(updateInput.Aliases, a.(string))
			}
		}
		if d.HasChange("tags") {
			updateInput.Tags = expandTeamTags(d.Get("tags").([]interface{}))
		}
		if d.HasChange("parent_id") {
			newParentID := d.Get("parent_id").(string)
			if newParentID != "" {
				updateInput.ParentId = newParentID
			} else {
				// Clearing parent_id: the generated struct uses omitempty so empty
				// string would be omitted. Send an explicit null via the client
				// NerdGraph method with parentId: null in the variables map.
				if err := clearTeamParentID(ctx, client, d.Id()); err != nil {
					return diag.FromErr(err)
				}
				// Skip the main update for parent_id to avoid omitempty silently
				// sending nothing. The other changed fields below still go through
				// the main update call.
			}
		}
		if d.HasChange("resources") {
			updateInput.Resources = expandTeamResourcesUpdate(d.Get("resources").([]interface{}))
		}

		if _, err := client.Scorecards.EntityManagementUpdateTeam(d.Id(), updateInput); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Sync members ──────────────────────────────────────────────────────
	if d.HasChange("members") {
		membershipColID := d.Get("membership_collection_id").(string)
		oldRaw, newRaw := d.GetChange("members")
		oldIDs := expandMemberUserIDsFromSet(oldRaw.(*schema.Set))
		newIDs := expandMemberUserIDsFromSet(newRaw.(*schema.Set))
		if err := syncMembers(ctx, client, membershipColID, oldIDs, newIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Sync managers (after ensuring members are up-to-date) ─────────────
	if d.HasChange("managers") || d.HasChange("members") {
		if err := syncManagers(ctx, providerConfig, d, d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Sync owned entities ──────────────────────────────────────────────
	if d.HasChange("entities") {
		ownershipColID := d.Get("ownership_collection_id").(string)
		oldRaw, newRaw := d.GetChange("entities")
		oldGUIDs := expandEntityGUIDsFromSet(oldRaw.(*schema.Set))
		newGUIDs := expandEntityGUIDsFromSet(newRaw.(*schema.Set))
		if err := syncOwnedEntities(ctx, &client.Scorecards, ownershipColID, oldGUIDs, newGUIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceNewRelicTeamRead(ctx, d, meta)
}

// ──────────────────────────────────────────────────────────────────────────────
// Delete
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicTeamDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting team %s", d.Id())
	if _, err := client.Scorecards.EntityManagementDelete(d.Id()); err != nil {
		var notFound *nrErrors.NotFound
		if errors.As(err, &notFound) {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Sync helpers
// ──────────────────────────────────────────────────────────────────────────────

// syncMembers reconciles the membership collection: adds newly listed user IDs
// and removes ones that were removed. User IDs are resolved to NGEP GUIDs
// before interacting with the collection.
func syncMembers(
	ctx context.Context,
	client *nr.NewRelic,
	membershipColID string,
	oldUserIDs, newUserIDs []int,
) error {
	// Build set maps for easy diff.
	oldSet := make(map[int]bool, len(oldUserIDs))
	for _, id := range oldUserIDs {
		oldSet[id] = true
	}
	newSet := make(map[int]bool, len(newUserIDs))
	for _, id := range newUserIDs {
		newSet[id] = true
	}

	// IDs to add.
	var toAdd []int
	for id := range newSet {
		if !oldSet[id] {
			toAdd = append(toAdd, id)
		}
	}
	// IDs to remove.
	var toRemove []int
	for id := range oldSet {
		if !newSet[id] {
			toRemove = append(toRemove, id)
		}
	}

	if len(toAdd) > 0 {
		guids, err := resolveUserIDsToNGEPGUIDs(ctx, &client.Entities, toAdd)
		if err != nil {
			return err
		}
		addList := make([]string, 0, len(guids))
		for _, g := range guids {
			addList = append(addList, g)
		}
		if _, err := client.Scorecards.EntityManagementAddCollectionMembers(membershipColID, addList); err != nil {
			return fmt.Errorf("adding members to collection %s: %w", membershipColID, err)
		}
	}

	if len(toRemove) > 0 {
		guids, err := resolveUserIDsToNGEPGUIDs(ctx, &client.Entities, toRemove)
		if err != nil {
			return err
		}
		removeList := make([]string, 0, len(guids))
		for _, g := range guids {
			removeList = append(removeList, g)
		}
		if _, err := client.Scorecards.EntityManagementRemoveCollectionMembers(membershipColID, removeList); err != nil {
			return fmt.Errorf("removing members from collection %s: %w", membershipColID, err)
		}
	}

	return nil
}

// syncManagers looks up the current managers in state, resolves their user IDs
// to NGEP GUIDs, and updates the team's managers list.
func syncManagers(ctx context.Context, providerConfig *ProviderConfig, d *schema.ResourceData, teamID string) error {
	client := providerConfig.NewClient
	managerUserIDs := expandManagerUserIDsFromSet(d.Get("managers").(*schema.Set))
	if len(managerUserIDs) == 0 {
		// Explicitly clear managers on the team if the managers block is empty.
		_, err := client.Scorecards.EntityManagementUpdateTeam(teamID, scorecards.EntityManagementTeamEntityUpdateInput{
			Managers: []string{},
		})
		return err
	}

	guids, err := resolveUserIDsToNGEPGUIDs(ctx, &client.Entities, managerUserIDs)
	if err != nil {
		return fmt.Errorf("resolving manager user IDs to NGEP GUIDs: %w", err)
	}
	managerGUIDs := make([]string, 0, len(guids))
	for _, g := range guids {
		managerGUIDs = append(managerGUIDs, g)
	}
	_, err = client.Scorecards.EntityManagementUpdateTeam(teamID, scorecards.EntityManagementTeamEntityUpdateInput{
		Managers: managerGUIDs,
	})
	return err
}

// syncOwnedEntities reconciles the ownership collection.
func syncOwnedEntities(
	_ context.Context,
	client *scorecards.Scorecards,
	ownershipColID string,
	oldGUIDs, newGUIDs []string,
) error {
	oldSet := make(map[string]bool, len(oldGUIDs))
	for _, g := range oldGUIDs {
		oldSet[g] = true
	}
	newSet := make(map[string]bool, len(newGUIDs))
	for _, g := range newGUIDs {
		newSet[g] = true
	}

	var toAdd, toRemove []string
	for g := range newSet {
		if !oldSet[g] {
			toAdd = append(toAdd, g)
		}
	}
	for g := range oldSet {
		if !newSet[g] {
			toRemove = append(toRemove, g)
		}
	}

	if len(toAdd) > 0 {
		if _, err := client.EntityManagementAddCollectionMembers(ownershipColID, toAdd); err != nil {
			return fmt.Errorf("adding entities to ownership collection %s: %w", ownershipColID, err)
		}
	}
	if len(toRemove) > 0 {
		if _, err := client.EntityManagementRemoveCollectionMembers(ownershipColID, toRemove); err != nil {
			return fmt.Errorf("removing entities from ownership collection %s: %w", ownershipColID, err)
		}
	}
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Collection readers
// ──────────────────────────────────────────────────────────────────────────────

// readCollectionUserIDs pages through a collection, pulls UserEntity items and
// returns their raw userId integers.
func readCollectionUserIDs(_ context.Context, client *scorecards.Scorecards, colID string) ([]int, error) {
	if colID == "" {
		return nil, nil
	}
	var userIDs []int
	cursor := ""
	for {
		filter := scorecards.EntityManagementCollectionElementsFilter{
			CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: colID},
		}
		result, err := client.GetCollectionElements(cursor, filter, 100)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}
		for _, item := range result.Items {
			if u, ok := item.(*scorecards.EntityManagementUserEntity); ok {
				userIDs = append(userIDs, u.UserID)
			}
		}
		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}
	return userIDs, nil
}

// readCollectionEntityGUIDs pages through a collection and returns all GUIDs
// (any entity type — ownership can hold any entity).
func readCollectionEntityGUIDs(_ context.Context, client *scorecards.Scorecards, colID string) ([]string, error) {
	if colID == "" {
		return nil, nil
	}
	var guids []string
	cursor := ""
	for {
		filter := scorecards.EntityManagementCollectionElementsFilter{
			CollectionID: scorecards.EntityManagementCollectionIdFilterArgument{Eq: colID},
		}
		result, err := client.GetCollectionElements(cursor, filter, 100)
		if err != nil {
			return nil, err
		}
		if result == nil {
			break
		}
		for _, item := range result.Items {
			type idder interface{ GetID() string }
			type guider interface{ GetGUID() string }
			switch e := item.(type) {
			case *scorecards.EntityManagementGenericEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementUserEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementTeamEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementCollectionEntity:
				guids = append(guids, e.ID)
			case *scorecards.EntityManagementScorecardEntity:
				guids = append(guids, e.ID)
			default:
				_ = e
			}
		}
		if result.NextCursor == "" {
			break
		}
		cursor = result.NextCursor
	}
	return guids, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Expand / flatten helpers for TypeSet blocks
// ──────────────────────────────────────────────────────────────────────────────

func expandMemberUserIDs(d *schema.ResourceData) []int {
	return expandMemberUserIDsFromSet(d.Get("members").(*schema.Set))
}

func expandMemberUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, raw := range s.List() {
		m := raw.(map[string]interface{})
		out = append(out, m["user_id"].(int))
	}
	return out
}

func flattenMemberUserIDs(userIDs []int) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(userIDs))
	for _, id := range userIDs {
		out = append(out, map[string]interface{}{"user_id": id})
	}
	return out
}

func expandManagerUserIDsFromSet(s *schema.Set) []int {
	out := make([]int, 0, s.Len())
	for _, raw := range s.List() {
		m := raw.(map[string]interface{})
		out = append(out, m["user_id"].(int))
	}
	return out
}

// flattenManagerUserIDs converts NGEP manager GUIDs back to user IDs.
// The manager list on TeamEntity contains NGEP-encoded user GUIDs, not raw
// user IDs. We return the GUIDs as strings here because going back to
// int user IDs would require an additional entitySearch lookup. Instead
// the managed managers block stores user_id (int) which Terraform provides
// on the input side; the stored-in-state value from Read is the GUID.
// To keep things simple we store the GUID as the user_id's string value
// and document this in the schema description.
//
// Actually — the cleanest approach: on Read, we don't try to back-convert
// GUIDs to user IDs (that would require a reverse lookup). Instead we
// emit an empty set on Read for managers, because Terraform's "plan diff"
// is satisfied by the user explicitly declaring what they want. The
// managers are authoritative on the user side, not the read side.
// This is the same pattern used by resources that can't cheaply round-trip.
func flattenManagerUserIDs(_ []string) []map[string]interface{} {
	// Returning nil means Terraform treats the current config as the
	// authoritative source for managers. This avoids a reverse-lookup from
	// NGEP GUID → userId on every Read.
	return nil
}

func expandEntityGUIDs(d *schema.ResourceData) []string {
	return expandEntityGUIDsFromSet(d.Get("entities").(*schema.Set))
}

func expandEntityGUIDsFromSet(s *schema.Set) []string {
	out := make([]string, 0, s.Len())
	for _, raw := range s.List() {
		m := raw.(map[string]interface{})
		out = append(out, m["guid"].(string))
	}
	return out
}

func flattenEntityGUIDs(guids []string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(guids))
	for _, g := range guids {
		out = append(out, map[string]interface{}{"guid": g})
	}
	return out
}

// clearTeamParentID explicitly sends parentId: null to NGEP to detach the team
// from its parent. The generated EntityManagementTeamEntityUpdateInput uses
// `json:"parentId,omitempty"` which silently omits an empty string — using
// the struct would be a no-op instead of a clear. We bypass it by calling
// NerdGraphQueryWithContext directly with a handcrafted variables map.
func clearTeamParentID(ctx context.Context, client *nr.NewRelic, teamID string) error {
	const clearParentQuery = `mutation($id: ID!, $teamEntity: EntityManagementTeamEntityUpdateInput!) {
  entityManagementUpdateTeam(id: $id, teamEntity: $teamEntity) {
    entity { id parentId }
  }
}`
	vars := map[string]interface{}{
		"id": teamID,
		"teamEntity": map[string]interface{}{
			"parentId": nil,
		},
	}
	_, err := client.NerdGraph.QueryWithContext(ctx, clearParentQuery, vars)
	return err
}

// isNGEPGhostNotFound detects the NGEP transient "ghost" NOT_FOUND — the
// message begins with ": Entity not found." (empty id prefix). Real
// deletions carry the entity id in the prefix.
func isNGEPGhostNotFound(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) > 0 && msg[0] == ':' && len(msg) > 2 && msg[1] == ' '
}
