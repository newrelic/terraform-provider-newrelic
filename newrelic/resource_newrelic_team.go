package newrelic

import (
	"context"
	"errors"
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
							Type:     schema.TypeString,
							Required: true,
							Description: "The resource type. Must be one of: ATLASSIAN_CONFLUENCE, ATLASSIAN_JIRA, " +
								"ATLASSIAN_JIRA_SCORECARDS, BASECAMP, BLAMELESS, EMAIL, FACEBOOK_WORKPLACE, " +
								"GITHUB, GITLAB, GOOGLE_CHAT, GOOGLE_CLOUD_PLATFORM, GOOGLE_DRIVE, " +
								"MICROSOFT_AZURE, MICROSOFT_SHAREPOINT, MICROSOFT_TEAMS, OPSGENIE, " +
								"OTHER_CONTACT, OTHER_LINK, PAGERDUTY, ROCKET_CHAT, SERVICENOW, " +
								"SKYPE, SLACK, ZENDESK.",
							ValidateFunc: validation.StringInSlice([]string{
								"ATLASSIAN_CONFLUENCE",
								"ATLASSIAN_JIRA",
								"ATLASSIAN_JIRA_SCORECARDS",
								"BASECAMP",
								"BLAMELESS",
								"EMAIL",
								"FACEBOOK_WORKPLACE",
								"GITHUB",
								"GITLAB",
								"GOOGLE_CHAT",
								"GOOGLE_CLOUD_PLATFORM",
								"GOOGLE_DRIVE",
								"MICROSOFT_AZURE",
								"MICROSOFT_SHAREPOINT",
								"MICROSOFT_TEAMS",
								"OPSGENIE",
								"OTHER_CONTACT",
								"OTHER_LINK",
								"PAGERDUTY",
								"ROCKET_CHAT",
								"SERVICENOW",
								"SKYPE",
								"SLACK",
								"ZENDESK",
							}, false),
						},
						"content": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The resource content (e.g. a URL). Must not be empty.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"title": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "A human-readable title for the resource.",
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
			// organization_id is resolved automatically from the provider configuration.
			// Customers should never supply it — it is fetched and stored as Computed.
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NGEP organization UUID. Resolved automatically from the provider account.",
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

	// organization_id is always resolved automatically — customers never supply it.
	orgID, err := getOrganizationID(ctx, providerConfig, "")
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
		for _, a := range v.(*schema.Set).List() {
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
	log.Printf("[INFO] Waiting for team %s to appear in entitySearch index", teamID)
	if err := waitForNGEPEntityIndexed(ctx, &client.Entities, teamID, d.Timeout(schema.TimeoutCreate)); err != nil {
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

	// ── Ownership collection (non-authoritative) ────────────────────────────────
	// The ownership collection is the authoritative source for manually-added
	// entities. Tag-discovery-placed entities DO appear in the collectionElements
	// API response but their concrete types (e.g. EntityManagementFleetEntity)
	// are not registered in UnmarshalEntityManagementEntityInterface, so they are
	// silently dropped before reaching readTeamOwnedEntityGUIDs and never cause
	// phantom drift. Any entity whose type IS known and is absent from the config
	// will show as drift and be removed on the next apply.
	entityGUIDs, ownerErr := readTeamOwnedEntityGUIDs(ctx, &client.Scorecards, team.Ownership.ID)
	if ownerErr != nil {
		log.Printf("[WARN] Could not read ownership collection for team %s: %v", d.Id(), ownerErr)
	} else {
		_ = d.Set("entities", flattenEntityGUIDs(entityGUIDs))
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
		// updHasFields tracks whether upd has any non-zero fields that need
		// to be sent to EntityManagementUpdateTeam. Fields cleared via raw
		// calls (description="", aliases=[], tags=[], parentId=nil) do NOT
		// populate upd — without this guard we would make a redundant API
		// call with an empty struct ({}) when only clear operations occurred.
		updHasFields := false

		if d.HasChange("name") {
			upd.Name = d.Get("name").(string)
			updHasFields = true
		}
		if d.HasChange("description") {
			newDesc := d.Get("description").(string)
			if newDesc != "" {
				// Non-empty: include in the main update mutation below.
				upd.Description = newDesc
				updHasFields = true
			} else {
				// Explicit clear: omitempty drops "", so use a raw call.
				if err := clearTeamDescriptionRaw(ctx, client, d.Id()); err != nil {
					return diag.Errorf("clearing description on team %s: %v", d.Id(), err)
				}
			}
		}
		if d.HasChange("aliases") {
			newAliases := d.Get("aliases").(*schema.Set).List()
			if len(newAliases) > 0 {
				for _, a := range newAliases {
					upd.Aliases = append(upd.Aliases, a.(string))
				}
				updHasFields = true
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
				updHasFields = true
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
			updHasFields = true
		}
		if d.HasChange("parent_id") {
			if newID := d.Get("parent_id").(string); newID != "" {
				upd.ParentId = newID
				updHasFields = true
			} else {
				// Clearing parent_id requires explicit null — the struct uses
				// omitempty so an empty string would be silently dropped.
				if err := clearTeamParentID(ctx, client, d.Id()); err != nil {
					return diag.FromErr(err)
				}
				// parent_id is now cleared via raw call; no field to add to upd.
			}
		}

		// Only issue the main update call when at least one field is non-zero.
		// When all changes were handled by raw clear calls (description, aliases,
		// tags, parentId), upd is empty and sending {} to the API would be a
		// wasteful no-op.
		if updHasFields {
			if _, err := client.Scorecards.EntityManagementUpdateTeam(d.Id(), upd); err != nil {
				return diag.FromErr(err)
			}
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
