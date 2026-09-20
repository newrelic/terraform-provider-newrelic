package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Additional searchable names for the team. Each alias must be unique within the organization.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:        schema.TypeSet,
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

	// ── Build create input ────────────────────────────────────────────────
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
	if v, ok := d.GetOk("tags"); ok {
		input.Tags = expandNGEPTags(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("resources"); ok {
		input.Resources = expandTeamResources(v.([]interface{}))
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
	// Membership, managers, and entities all require the team to exist first
	// (the backing collections are auto-created by NGEP on team creation).
	// applyTeamCollections handles all three in the correct dependency order.
	// Old slices are nil because nothing existed before this create call.
	if err := applyTeamCollections(ctx, client, teamID, membershipColID, ownershipColID,
		nil, expandUserIDsFromSet(d.Get("members").(*schema.Set)),
		expandUserIDsFromSet(d.Get("managers").(*schema.Set)),
		nil, expandEntityGUIDsFromSet(d.Get("entities").(*schema.Set)),
	); err != nil {
		return diag.FromErr(err)
	}

	// ── Indexing gate ─────────────────────────────────────────────────────
	// NGEP syncs entities to the standard entitySearch index asynchronously.
	// Wait until the team is visible before returning from Create.
	log.Printf("[INFO] Waiting for team %s to appear in entitySearch index", teamID)
	if err := waitForNGEPEntityIndexed(ctx, &client.Entities, teamID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
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
	_ = d.Set("tags", flattenNGEPTags(team.Tags))

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

	// ── Ownership collection — non-authoritative split ───────────────────────
	// Fetch org-level discovery settings to separate statically-managed entities
	// from those auto-assigned via tag-based rules. Terraform only tracks the
	// static subset in state; dynamic entities are left untouched.
	orgSettings, orgErr := client.Scorecards.GetTeamsOrganizationSettings()
	if orgErr != nil {
		log.Printf("[WARN] Could not read TeamsOrganizationSettings for team %s: %v — treating all ownership entities as static", d.Id(), orgErr)
	}

	staticGUIDs, dynamicCount, ownerErr := readStaticOwnershipGUIDs(
		ctx,
		&client.Scorecards,
		&client.Entities,
		team.Ownership.ID,
		team.Name,
		team.Aliases,
		orgSettings,
	)
	if ownerErr != nil {
		log.Printf("[WARN] Could not read ownership collection for team %s: %v", d.Id(), ownerErr)
	} else {
		_ = d.Set("entities", flattenEntityGUIDs(staticGUIDs))
	}

	// Warn when dynamic (tag-auto-assigned) entities are present in the
	// ownership collection and the user has declared an entities block.
	// We do NOT show a diff for dynamic entities — they are outside Terraform's
	// management scope. The warning nudges users to consolidate ownership under
	// Terraform if they want full declarative control.
	if dynamicCount > 0 && len(d.Get("entities").(*schema.Set).List()) > 0 {
		return diag.Diagnostics{
			{
				Severity: diag.Warning,
				Summary:  "Team ownership contains tag-auto-assigned entities not managed by Terraform",
				Detail: fmt.Sprintf(
					"%d entity/entities in team %q's ownership collection were assigned automatically "+
						"via tag-based discovery rules and are not tracked in this resource's `entities` block. "+
						"Terraform will only manage the %d entity/entities you have explicitly declared. "+
						"If you want full declarative control over team ownership, add all owned entities to "+
						"the `entities` block and disable automatic tag-based assignment for this team.",
					dynamicCount, team.Name, len(staticGUIDs),
				),
			},
		}
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
			newDesc := d.Get("description").(string)
			if newDesc != "" {
				upd.Description = newDesc
			} else {
				// Explicit clear: omitempty drops "", so use a raw call.
				if err := clearTeamDescriptionRaw(ctx, client, d.Id()); err != nil {
					return diag.Errorf("clearing description on team %s: %v", d.Id(), err)
				}
			}
		}
		if d.HasChange("aliases") {
			newAliases := d.Get("aliases").([]interface{})
			if len(newAliases) > 0 {
				for _, a := range newAliases {
					upd.Aliases = append(upd.Aliases, a.(string))
				}
			} else {
				// Explicit clear: omitempty would drop nil, so use raw call.
				if err := clearTeamAliasesRaw(ctx, client, d.Id()); err != nil {
					return diag.Errorf("clearing aliases on team %s: %v", d.Id(), err)
				}
			}
		}
		if d.HasChange("tags") {
			userTags := expandNGEPTags(d.Get("tags").(*schema.Set).List())
			// Always merge with current system tags (e.g. nr.hierarchy.level) —
			// NGEP rejects any update that would remove tags prefixed with "nr.".
			sysTags := fetchEntitySystemTags(ctx, &client.Scorecards, d.Id())
			mergedTags := mergeWithSystemTags(userTags, sysTags)

			if len(mergedTags) > 0 {
				// Normal path: send user tags + preserved system tags.
				upd.Tags = mergedTags
			} else {
				// Edge case: user clearing all tags on a team with no system tags.
				// omitempty would silently drop an empty slice, so use a raw call.
				if err := clearTeamTagsRaw(ctx, client, d.Id()); err != nil {
					return diag.Errorf("clearing tags on team %s: %v", d.Id(), err)
				}
			}
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

	// ── Collections ────────────────────────────────────────────────────────
	// Reconcile members, managers, and owned entities whenever any of them
	// change. All three are handled by applyTeamCollections in the correct
	// dependency order (membership must be settled before managers are set).
	if d.HasChange("members") || d.HasChange("managers") || d.HasChange("entities") {
		oldMembersRaw, newMembersRaw := d.GetChange("members")
		_, newManagersRaw := d.GetChange("managers")
		oldEntitiesRaw, newEntitiesRaw := d.GetChange("entities")
		if err := applyTeamCollections(ctx, client, d.Id(),
			d.Get("membership_collection_id").(string),
			d.Get("ownership_collection_id").(string),
			expandUserIDsFromSet(oldMembersRaw.(*schema.Set)),
			expandUserIDsFromSet(newMembersRaw.(*schema.Set)),
			expandUserIDsFromSet(newManagersRaw.(*schema.Set)),
			expandEntityGUIDsFromSet(oldEntitiesRaw.(*schema.Set)),
			expandEntityGUIDsFromSet(newEntitiesRaw.(*schema.Set)),
		); err != nil {
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
