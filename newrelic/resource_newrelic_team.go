package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	newrelic "github.com/newrelic/newrelic-client-go/v2/newrelic"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/servicearchintelligence"
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
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the team. Must be unique within the organization.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A free-text description of the team.",
			},
			"aliases": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Additional searchable names for the team. Each alias must be unique within the organization.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Tags to assign to this resource. Each tag has a key and one or more values.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The tag key.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The tag values.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["key"].(string))
				},
			},
			// ── Hierarchy ───────────────────────────────────────────────────
			"parent_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Entity management GUID of the parent team. Cannot be set to the team's own GUID.",
			},
			// ── Resources (links / docs) ─────────────────────────────────────
			// resources is a TypeSet (not TypeList) because order is not meaningful —
			// the API stores them as an unordered collection and returns them in
			// arbitrary order, which would cause phantom ordering diffs with TypeList.
			// The set hash is based on type+content, the natural unique key.
			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Supplemental resource links attached to this team (e.g. runbooks, chat channels, on-call schedules). " +
					"These are displayed in the team's Settings page in the New Relic UI under Contacts or Links, " +
					"automatically categorised by type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Resource category. Valid values are defined in TeamResourceTypeValues() in the go-client and stay in sync with the API automatically.",
							// ValidateFunc delegates to the go-client's TeamResourceTypeValues() so that
							// adding a new resource type to the go-client automatically makes it valid
							// in the provider without any code change here.
							ValidateFunc: validation.StringInSlice(servicearchintelligence.TeamResourceTypeValues(), false),
						},
						"content": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The URL or contact address for this resource.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"title": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "A short label shown in the UI for this resource.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			// ── Membership ──────────────────────────────────────────────────
			// user_id is the integer NR user ID. Terraform resolves it to the
			// EntityManagementUserEntity GUID before calling the NGEP API.
			"members": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of New Relic user IDs (integers) to add as members of this team. Managers must be a subset of members.",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			"managers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Set of New Relic user IDs (integers) to designate as team managers. Each ID must also appear in members.",
				Elem: &schema.Schema{
					Type:         schema.TypeInt,
					ValidateFunc: validation.IntAtLeast(1),
				},
			},
			// ── Ownership ───────────────────────────────────────────────────
			"entity_management_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "managed",
				Description: "Controls how Terraform manages this team's entity ownership collection.\n\n" +
					"  • `managed` (default): Terraform tracks the `entities` block and reconciles it with the " +
					"    ownership collection. Out-of-band additions appear as drift. All entity warnings are shown.\n" +
					"  • `unmanaged`: Terraform does not control entity ownership. The ownership collection is left " +
					"    entirely to tag-based discovery and/or manual UI management. No drift or warnings " +
					"    are shown for entities, and the `entities` attribute cannot be set in this mode.",
				ValidateFunc: validation.StringInSlice([]string{"managed", "unmanaged"}, false),
			},
			"entities": {
				Type:     schema.TypeSet,
				Optional: true,
				// Computed: true is required so CustomizeDiff.SetNew can suppress the
				// entities diff during entity_management_mode transitions. It also means
				// that when entities is absent from config (unmanaged mode), Terraform
				// uses the provider's computed value rather than planning a removal.
				Computed:    true,
				Description: "Set of entity GUIDs that this team owns. Added to the team's auto-created ownership collection.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
			// ── Computed / infrastructure ────────────────────────────────────
			// organization_id is resolved automatically from the provider configuration.
			// Customers should never supply it — it is fetched and stored as Computed.
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization UUID. Resolved automatically from the provider account.",
			},
			"membership_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created team membership collection. Read-only.",
			},
			"ownership_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created team ownership collection. Read-only.",
			},
			// hierarchy_level_id is set automatically by the platform based on the
			// team's position in the organisational hierarchy (determined by parent_id).
			// It matches one of the level GUIDs in newrelic_teams_organization_settings
			// hierarchy_levels. Cannot be set directly.
			"hierarchy_level_id": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "GUID of the hierarchy level this team is assigned to. Set automatically " +
					"by the platform based on the team's parentId chain. Matches a level GUID " +
					"from the newrelic_teams_organization_settings hierarchy_levels list.",
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

	// organization_id is always resolved automatically — customers never supply it.
	orgID, err := getOrganizationID(ctx, providerConfig, "")
	if err != nil {
		return diag.FromErr(err)
	}

	// ── Build create input ────────────────────────────────────────────────
	input := servicearchintelligence.EntityManagementTeamEntityCreateInput{
		Name: d.Get("name").(string),
		Scope: servicearchintelligence.EntityManagementScopedReferenceInput{
			ID:   orgID,
			Type: servicearchintelligence.EntityManagementEntityScopeTypes.ORGANIZATION,
		},
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("aliases"); ok {
		for _, a := range v.(*schema.Set).List() {
			input.Aliases = append(input.Aliases, a.(string))
		}
	}
	if v, ok := d.GetOk("tags"); ok {
		input.Tags = expandSAITags(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("resources"); ok {
		input.Resources = expandTeamResources(v.(*schema.Set).List())
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

	// ── Sync collections ───────────────────────────────────────────────────
	// Membership, managers, and entities all need the team to exist first
	// (collections are auto-created by NGEP on team creation).
	// applyTeamCollections orchestrates all three in the right dependency order.
	// Old slices are nil because nothing existed before this create call.
	//
	// In unmanaged mode we skip entity sync — the collection is left to
	// tag-based discovery and/or the Teams UI.
	// In managed mode we sync whatever the entities block declares. If it is
	// absent from config, d.Get returns an empty set and the sync is a no-op.
	// If a declared GUID is already in the collection (e.g. a tag-discovered
	// entity the user is now taking declarative control of), NGEP returns
	// "already belongs to collection" which syncTeamOwnership treats as success
	// and the GUID is saved to state.
	mode := d.Get("entity_management_mode").(string)
	var newEntityGUIDs []string
	if mode == "managed" {
		newEntityGUIDs = expandEntityGUIDsFromSet(d.Get("entities").(*schema.Set))
	}
	if err := applyTeamCollections(ctx, client, teamID, membershipColID, ownershipColID,
		nil, expandUserIDsFromSet(d.Get("members").(*schema.Set)),
		nil, expandUserIDsFromSet(d.Get("managers").(*schema.Set)), // oldManagers=nil: fresh create
		nil, newEntityGUIDs,
	); err != nil {
		return diag.FromErr(err)
	}

	// In unmanaged mode, clear entities from state — Terraform does not track
	// the ownership collection when NGEP/UI manages it exclusively.
	if mode == "unmanaged" {
		_ = d.Set("entities", schema.NewSet(schema.HashString, []interface{}{}))
	}

	// Set remaining state from input — no Read round-trip needed.
	_ = d.Set("name", input.Name)
	_ = d.Set("description", input.Description)
	_ = d.Set("parent_id", input.ParentId)
	if len(input.Aliases) > 0 {
		_ = d.Set("aliases", input.Aliases)
	}
	if v := tagsInputToFlattenedSet(input.Tags); v != nil {
		_ = d.Set("tags", v)
	}

	// Indexing gate: block until the entity is visible in entitySearch so that
	// subsequent Read calls find it immediately.
	log.Printf("[INFO] Waiting for team %s to appear in entityManagement index", teamID)
	if err := waitForNGEPEntityIndexed(ctx, &client.Scorecards, teamID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return nil
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
		// If entity not found (deleted outside Terraform), remove from state.
		if errors.As(err, &notFound) {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	if entityIface == nil {
		d.SetId("")
		return nil
	}

	team, ok := (*entityIface).(*servicearchintelligence.EntityManagementTeamEntity)
	if !ok {
		return diag.Errorf("entity %s is not a TeamEntity", d.Id())
	}

	_ = d.Set("name", team.Name)
	_ = d.Set("description", team.Description)
	_ = d.Set("parent_id", team.ParentId)
	_ = d.Set("organization_id", team.Scope.ID)
	_ = d.Set("membership_collection_id", team.Membership.ID)
	_ = d.Set("ownership_collection_id", team.Ownership.ID)
	_ = d.Set("hierarchy_level_id", team.HierarchyLevelId)
	_ = d.Set("resources", flattenTeamResources(team.Resources))
	_ = d.Set("tags", flattenSAITags(team.Tags))

	// aliases: the entity read path has an eventual-consistency lag — the API
	// may return nil for a freshly created team even though aliases were
	// stored. Only overwrite state when the API returns a non-nil list.
	if team.Aliases != nil {
		_ = d.Set("aliases", team.Aliases)
	}

	// Read the membership collection as a GUID→userID map. This serves double
	// duty: it populates the members list AND lets us decode manager GUIDs back
	// to integer userIds for idempotent state storage.
	guidToUserID, err := readTeamMembershipMap(ctx, &client.Scorecards, team.Membership.ID)
	if err != nil {
		log.Printf("[WARN] Could not read membership collection for team %s: %v", d.Id(), err)
	} else {
		memberUserIDs := make([]int, 0, len(guidToUserID))
		for _, uid := range guidToUserID {
			memberUserIDs = append(memberUserIDs, uid)
		}
		_ = d.Set("members", flattenMemberUserIDs(memberUserIDs))
		_ = d.Set("managers", decodeManagerGUIDsToUserIDs(team.Managers, guidToUserID))
	}

	// ── Ownership collection — non-authoritative split ──────────────────────────
	// In unmanaged mode, skip all entity operations entirely and clear entities
	// from state. The ownership collection is controlled by NGEP (tag discovery
	// or manual UI); Terraform does not track or show drift for it.
	mode := d.Get("entity_management_mode").(string)
	// entity_management_mode is never stored by the API — it is Terraform-only
	// state. On terraform import, d.Get() returns "" (the zero value). Default
	// to "managed" so the imported state matches the schema Default and the
	// ImportStateVerify round-trip passes.
	if mode == "" {
		mode = "managed"
		_ = d.Set("entity_management_mode", "managed")
	}
	if mode == "unmanaged" {
		_ = d.Set("entities", schema.NewSet(schema.HashString, []interface{}{}))
		return nil
	}

	// NGEP places both manually-added and tag-matched entities in the same
	// ownership collection. readStaticOwnershipGUIDs separates them:
	//   - staticGUIDs  → entities tracked in state; drift shown if absent from config
	//   - discoveryGUIDs → entities managed by NGEP tag discovery; warning only
	//
	// declaredGUIDs: entities the user has explicitly put in the entities block.
	// These are ALWAYS treated as static regardless of their tags — a Terraform
	// declaration trumps tag-based discovery classification, preventing a perpetual
	// diff if a user adds a discovery tag to an entity they also declare explicitly.
	declaredGUIDs := expandEntityGUIDsFromSet(d.Get("entities").(*schema.Set))

	orgSettings, orgErr := client.Scorecards.GetTeamsOrganizationSettings()
	if orgErr != nil {
		log.Printf("[WARN] Could not read TeamsOrganizationSettings for team %s: %v — treating all ownership entities as static", d.Id(), orgErr)
	}

	staticGUIDs, discoveryGUIDs, ownerErr := readStaticOwnershipGUIDs(
		ctx,
		&client.Scorecards,
		&client.Entities,
		team.Ownership.ID,
		team.Name,
		team.Aliases,
		declaredGUIDs,
		orgSettings,
	)
	if ownerErr != nil {
		log.Printf("[WARN] Could not read ownership collection for team %s: %v", d.Id(), ownerErr)
	} else {
		// Always update state so the Computed value stays fresh. When entities is
		// not declared in config (Computed path), this prevents spurious drift by
		// keeping state in sync with what NGEP currently holds.
		_ = d.Set("entities", flattenEntityGUIDs(staticGUIDs))
	}

	// Warnings below are informational — they fire whenever there are entities
	// in the collection that Terraform should tell the user about. The Computed:true
	// schema attribute handles drift suppression when entities is absent from
	// config. Warnings always fire when there is something to warn about.

	// Build two diagnostic warnings:
	//  1. Tag-discovery warning — entities auto-assigned by NGEP, not tracked by Terraform
	//  2. Out-of-band addition warning — manually added outside Terraform; plan will remove them
	var diags diag.Diagnostics

	if len(discoveryGUIDs) > 0 {
		tagKeyHint := "team"
		if orgSettings != nil && len(orgSettings.Discovery.TagKeys) > 0 {
			tagKeyHint = strings.Join(orgSettings.Discovery.TagKeys, "/")
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf(
				"Team %q: %d entity/entities in collection are tag-discovered — excluded from drift detection",
				team.Name, len(discoveryGUIDs),
			),
			Detail: fmt.Sprintf(
				"The following entity GUIDs are present in the ownership collection because "+
					"their `tags.%s` value matches the team name or an alias. "+
					"Terraform excludes these from drift — they are currently managed by "+
					"tag-based discovery, not by this resource.\n\n"+
					"Discovery GUIDs:\n  %s\n\n"+
					"To take Terraform control of one of these entities (recommended):\n"+
					"  1. Add its GUID to the `entities` block and run apply.\n"+
					"     The provider treats the NGEP 'already in collection' response as success\n"+
					"     and begins tracking it — no manual collection changes needed.\n"+
					"  2. Optionally, remove the team tag from the entity afterwards so that\n"+
					"     tag-based discovery no longer manages it. This is cleaner but not required.\n\n"+
					"To keep ALL ownership governed by tag-based discovery only:\n"+
					"  Omit the `entities` block entirely — Terraform will not show drift.\n\n"+
					"To disable entity tracking and warnings completely:\n"+
					"  Set `entity_management_mode = \"unmanaged\"` on this resource.",
				tagKeyHint,
				strings.Join(discoveryGUIDs, "\n  "),
			),
		})
	}

	// Identify static entities in the collection that are not in the user's config.
	// Use the config-declared set (not the state-declared set) so the warning is
	// suppressed as soon as the user adds the entity to their entities block —
	// even on the apply that clears the drift, not just on the following plan.
	configSet := declaredSetFromConfig(d, declaredGUIDs)
	var outOfBandGUIDs []string
	for _, g := range staticGUIDs {
		if !configSet[g] {
			outOfBandGUIDs = append(outOfBandGUIDs, g)
		}
	}
	if len(outOfBandGUIDs) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary: fmt.Sprintf(
				"Team %q: %d entity/entities were added to the collection outside Terraform",
				team.Name, len(outOfBandGUIDs),
			),
			Detail: fmt.Sprintf(
				"The following entity GUIDs exist in the ownership collection but are NOT "+
					"declared in the `entities` block of this resource. The plan above shows "+
					"them as '-' (removal) — running apply will remove them from the collection "+
					"to restore the declared configuration.\n\n"+
					"Out-of-band GUIDs (will be removed on apply):\n  %s\n\n"+
					"To KEEP them: add their GUIDs to the `entities` block before applying.\n"+
					"To REMOVE them: run terraform apply (this is what the plan will do).",
				strings.Join(outOfBandGUIDs, "\n  "),
			),
		})
	}

	return diags
}

// buildTeamEntityUpdateInput assembles the team update mutation payload from
// ResourceData changes. Pointer fields on EntityManagementTeamEntityUpdateInput:
//   - nil                → field omitted → API leaves the attribute unchanged
//   - pointer to ""/ []  → field sent → API clears the attribute
//
// Only attributes that d.HasChange() reports as changed are populated;
// unchanged attributes stay nil and are omitted from the JSON payload,
// matching the update pattern used throughout this provider.
func buildTeamEntityUpdateInput(ctx context.Context, d *schema.ResourceData, client *newrelic.NewRelic) (servicearchintelligence.EntityManagementTeamEntityUpdateInput, diag.Diagnostics) {
	upd := servicearchintelligence.EntityManagementTeamEntityUpdateInput{}

	if d.HasChange("name") {
		upd.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		v := d.Get("description").(string)
		upd.Description = &v
	}
	if d.HasChange("aliases") {
		aliases := make([]string, 0)
		for _, a := range d.Get("aliases").(*schema.Set).List() {
			aliases = append(aliases, a.(string))
		}
		upd.Aliases = &aliases
	}
	if d.HasChange("tags") {
		merged := mergeWithSystemTags(
			expandSAITags(d.Get("tags").(*schema.Set).List()),
			fetchEntitySystemTags(ctx, &client.Scorecards, d.Id()),
		)
		upd.Tags = &merged
	}
	if d.HasChange("resources") {
		resources := expandTeamResourcesUpdate(d.Get("resources").(*schema.Set).List())
		upd.Resources = &resources
	}
	if d.HasChange("parent_id") {
		v := d.Get("parent_id").(string)
		upd.ParentId = &v
	}
	return upd, nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicTeamUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// ── Team entity fields ─────────────────────────────────────────────────
	// Only invoke the mutation when a relevant field actually changed.
	// buildTeamEntityUpdateInput populates only changed fields (pointer
	// semantics: nil = no change, ptr-to-zero = clear); unchanged fields
	// are omitted from the JSON payload and the API leaves them intact.
	if d.HasChangesExcept("members", "entities", "managers", "entity_management_mode") {
		upd, diags := buildTeamEntityUpdateInput(ctx, d, client)
		if diags != nil {
			return diags
		}
		if _, err := client.Scorecards.EntityManagementUpdateTeam(d.Id(), upd); err != nil {
			return diag.FromErr(err)
		}
	}

	// ── Collections ────────────────────────────────────────────────────────
	// Reconcile members, managers, and owned entities whenever any of them
	// change. applyTeamCollections handles all three in the correct dependency
	// order (membership must be committed before managers are set).
	//
	// Entity sync fires when entities changed AND the mode is not "unmanaged".
	// When entities is absent from config, Computed:true means the planned
	// value equals the prior state value, so d.HasChange returns false and the
	// collection is never touched — no extra guard needed.
	// When entities is present and a declared GUID is already in the collection
	// (e.g. a tag-discovered entity being taken into declarative control),
	// syncTeamOwnership treats the "already belongs" response as success and
	// the GUID is saved to state.
	mode := d.Get("entity_management_mode").(string)
	entitiesChanged := mode == "managed" && d.HasChange("entities")
	if d.HasChange("members") || d.HasChange("managers") || entitiesChanged {
		oldMembersRaw, newMembersRaw := d.GetChange("members")
		oldManagersRaw, newManagersRaw := d.GetChange("managers")

		var oldEntities, newEntities []string
		if entitiesChanged {
			oldEntitiesRaw, newEntitiesRaw := d.GetChange("entities")
			oldEntities = expandEntityGUIDsFromSet(oldEntitiesRaw.(*schema.Set))
			newEntities = expandEntityGUIDsFromSet(newEntitiesRaw.(*schema.Set))
		}

		if err := applyTeamCollections(ctx, client, d.Id(),
			d.Get("membership_collection_id").(string),
			d.Get("ownership_collection_id").(string),
			expandUserIDsFromSet(oldMembersRaw.(*schema.Set)),
			expandUserIDsFromSet(newMembersRaw.(*schema.Set)),
			expandUserIDsFromSet(oldManagersRaw.(*schema.Set)),
			expandUserIDsFromSet(newManagersRaw.(*schema.Set)),
			oldEntities,
			newEntities,
		); err != nil {
			return diag.FromErr(err)
		}
	}

	// If switching to unmanaged mode, clear entities from state but DO NOT touch
	// the collection. NGEP continues to manage what's in the collection; Terraform
	// just stops tracking it.
	if d.HasChange("entity_management_mode") {
		newMode := d.Get("entity_management_mode").(string)
		if newMode == "unmanaged" {
			_ = d.Set("entities", schema.NewSet(schema.HashString, []interface{}{}))
		}
	}

	// Terraform's ResourceData already tracks all changed values; no Read
	// round-trip is needed.
	return nil
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
