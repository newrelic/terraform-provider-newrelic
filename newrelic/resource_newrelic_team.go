package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		CustomizeDiff: resourceNewRelicTeamCustomizeDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// ── Core identity ───────────────────────────────────────────────
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the team. Must be unique within the organization.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free-text description of the team.",
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
			// ── Hierarchy ───────────────────────────────────────────────────
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NGEP GUID of the parent team. Cannot be set to the team's own GUID.",
			},
			// ── Resources (links / docs) ─────────────────────────────────────
			"resources": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Supplemental resources (e.g. runbooks, wikis). Each item requires type and content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type (e.g. 'link').",
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
			// ── Membership ──────────────────────────────────────────────────
			// user_id is the integer NR user ID. Terraform resolves it to the
			// EntityManagementUserEntity GUID before calling the NGEP API.
			"members": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of New Relic user IDs (integers) to add as members of this team. " +
					"These are stored in the team's auto-created membership collection. " +
					"Managers must be a subset of members.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The integer New Relic user ID. Obtainable from the newrelic_user data source.",
						},
					},
				},
			},
			// ── Managers ────────────────────────────────────────────────────
			"managers": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of New Relic user IDs (integers) who are managers of this team. " +
					"Every manager must also be present in the members block. " +
					"Validated at plan time via CustomizeDiff.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The integer New Relic user ID of the manager.",
						},
					},
				},
			},
			// ── Ownership ───────────────────────────────────────────────────
			"entities": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of entity GUIDs that this team owns. " +
					"Added to the team's auto-created ownership collection. " +
					"Entities from any account in the organization are accepted. " +
					"Deleted entities are removed from the collection automatically by NGEP.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The GUID of the entity to assign ownership to.",
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
				Description: "The NGEP organization UUID. Auto-fetched from the account if omitted.",
			},
			"membership_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created membership collection. Read-only.",
			},
			"ownership_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created ownership collection. Read-only.",
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

	// ── Build minimal create input ─────────────────────────────────────────
	// tags and resources are deliberately excluded from the create payload.
	// The NGEP validation service intermittently returns HTTP 500 when
	// resources is present in entityManagementCreateTeam. Both fields are
	// reliable when sent via the update mutation, so we always apply them in
	// the post-create update step below.
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

	// ── Create the team entity ─────────────────────────────────────────────
	result, err := client.Scorecards.EntityManagementCreateTeam(input)
	if err != nil {
		return diag.FromErr(err)
	}

	teamID := result.Entity.ID
	membershipColID := result.Entity.Membership.ID
	ownershipColID := result.Entity.Ownership.ID

	log.Printf("[INFO] Created NGEP team %s (membership=%s ownership=%s)", teamID, membershipColID, ownershipColID)

	d.SetId(teamID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set("membership_collection_id", membershipColID)
	_ = d.Set("ownership_collection_id", ownershipColID)

	// ── Post-create update: tags + resources ───────────────────────────────
	// Applied unconditionally so the resource is fully configured in one
	// terraform apply without requiring a second iteration. Mirrors the
	// filter_current_dashboard pattern in newrelic_one_dashboard.
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
			return diag.Errorf("post-create update (tags/resources) on team %s: %v", teamID, err)
		}
	}

	// ── Post-create update: members ────────────────────────────────────────
	// Members and managers cannot be set at create time (NGEP requires the
	// membership collection to exist first, which only happens after create).
	// We apply them here so a single terraform apply fully configures the team.
	if wantedMembers := expandMemberUserIDsFromSet(d.Get("members").(*schema.Set)); len(wantedMembers) > 0 {
		if err := syncMembers(ctx, client, membershipColID, nil, wantedMembers); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Post-create update: managers ──────────────────────────────────────
	// Managers are set after members — NGEP validates managers ⊆ members.
	if wantedManagers := expandManagerUserIDsFromSet(d.Get("managers").(*schema.Set)); len(wantedManagers) > 0 {
		if err := syncManagers(ctx, client, teamID, wantedManagers); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Post-create update: owned entities ───────────────────────────────
	if wantedEntities := expandEntityGUIDsFromSet(d.Get("entities").(*schema.Set)); len(wantedEntities) > 0 {
		if err := syncOwnedEntities(ctx, &client.Scorecards, ownershipColID, nil, wantedEntities); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Indexing gate ─────────────────────────────────────────────────────
	// NGEP entities are asynchronously synced to the standard entitySearch
	// index. The Teams UI polls for this after create; we do the same so
	// Terraform doesn't return until the entity is queryable.
	log.Printf("[INFO] Waiting for team %s to appear in entitySearch index", teamID)
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := client.Entities.GetEntitySearchByQueryWithContext(
			ctx,
			entities.EntitySearchOptions{},
			fmt.Sprintf("id IN ('%s')", teamID),
			[]entities.EntitySearchSortCriteria{},
		)
		if err != nil {
			return resource.NonRetryableError(
				fmt.Errorf("entitySearch failed while waiting for team %s: %w", teamID, err))
		}
		if res == nil || len(res.Results.Entities) == 0 {
			return resource.RetryableError(
				fmt.Errorf("team %s not yet visible in entitySearch — retrying", teamID))
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
		if errors.As(err, &notFound) || isNGEPGhostNotFound(err) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if entityIface == nil {
		d.SetId("")
		return nil
	}

	team, ok := (*entityIface).(*scorecards.EntityManagementTeamEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamEntity", d.Id())
	}

	_ = d.Set("name", team.Name)
	_ = d.Set("description", team.Description)
	_ = d.Set("parent_id", team.ParentId)
	_ = d.Set("organization_id", team.Scope.ID)
	_ = d.Set("membership_collection_id", team.Membership.ID)
	_ = d.Set("ownership_collection_id", team.Ownership.ID)
	_ = d.Set("resources", flattenTeamResources(team.Resources))
	_ = d.Set("tags", flattenTeamTags(team.Tags))

	// aliases: the entity read path has an eventual-consistency lag — the API
	// may return nil for a freshly created team even though aliases were
	// stored. Only overwrite state when the API returns a non-nil list.
	if team.Aliases != nil {
		_ = d.Set("aliases", team.Aliases)
	}

	// Read the membership collection as a GUID→userID map. This serves double
	// duty: it populates the members list AND lets us decode manager GUIDs back
	// to integer userIds for idempotent state storage.
	guidToUserID, err := readCollectionMembersMap(ctx, &client.Scorecards, team.Membership.ID)
	if err != nil {
		log.Printf("[WARN] Could not read membership collection for team %s: %v", d.Id(), err)
	} else {
		memberUserIDs := make([]int, 0, len(guidToUserID))
		for _, uid := range guidToUserID {
			memberUserIDs = append(memberUserIDs, uid)
		}
		_ = d.Set("members", flattenMemberUserIDs(memberUserIDs))
		_ = d.Set("managers", flattenManagerUserIDs(team.Managers, guidToUserID))
	}

	ownerGUIDs, err := readCollectionEntityGUIDs(ctx, &client.Scorecards, team.Ownership.ID)
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

	// ── Team entity fields ─────────────────────────────────────────────────
	teamFieldsChanged := d.HasChangesExcept("members", "entities", "managers")
	if teamFieldsChanged {
		upd := scorecards.EntityManagementTeamEntityUpdateInput{}
		if d.HasChange("name") {
			upd.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			upd.Description = d.Get("description").(string)
		}
		if d.HasChange("aliases") {
			for _, a := range d.Get("aliases").([]interface{}) {
				upd.Aliases = append(upd.Aliases, a.(string))
			}
		}
		if d.HasChange("tags") {
			upd.Tags = expandTeamTags(d.Get("tags").([]interface{}))
		}
		if d.HasChange("resources") {
			upd.Resources = expandTeamResourcesUpdate(d.Get("resources").([]interface{}))
		}
		if d.HasChange("parent_id") {
			if newID := d.Get("parent_id").(string); newID != "" {
				upd.ParentId = newID
			} else {
				// Clearing parent_id requires explicit null — the struct uses
				// omitempty so an empty string would be silently dropped.
				if err := clearTeamParentID(ctx, client, d.Id()); err != nil {
					return diag.FromErr(err)
				}
				// parent_id is now cleared; other changed fields still need
				// the main update call below.
			}
		}

		if _, err := client.Scorecards.EntityManagementUpdateTeam(d.Id(), upd); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Members ────────────────────────────────────────────────────────────
	if d.HasChange("members") {
		membershipColID := d.Get("membership_collection_id").(string)
		oldRaw, newRaw := d.GetChange("members")
		oldIDs := expandMemberUserIDsFromSet(oldRaw.(*schema.Set))
		newIDs := expandMemberUserIDsFromSet(newRaw.(*schema.Set))
		if err := syncMembers(ctx, client, membershipColID, oldIDs, newIDs); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Managers (after members are up-to-date) ────────────────────────────
	if d.HasChange("managers") || d.HasChange("members") {
		if err := syncManagers(ctx, client, d.Id(),
			expandManagerUserIDsFromSet(d.Get("managers").(*schema.Set))); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Owned entities ─────────────────────────────────────────────────────
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
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Deleting NGEP team %s", d.Id())
	if _, err := client.Scorecards.EntityManagementDelete(d.Id()); err != nil {
		var notFound *nrErrors.NotFound
		if errors.As(err, &notFound) {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}
