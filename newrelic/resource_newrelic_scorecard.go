package newrelic

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// customizeScorecardDiff is intentionally left empty — progress_levels is now
// a TypeSet (hashed on "id"), which makes ordering irrelevant at the schema
// level. Terraform natively handles the ForceNew via the schema's ForceNew: true
// flag when any set element actually changes. No manual diff override is needed.

func resourceNewRelicScorecard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicScorecardCreate,
		ReadContext:   resourceNewRelicScorecardRead,
		UpdateContext: resourceNewRelicScorecardUpdate,
		DeleteContext: resourceNewRelicScorecardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the scorecard.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the scorecard.",
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
			// progress_levels are set at create time only. The NGEP API currently
			// rejects updates to progress levels due to a backend validation bug.
			// Define them once when the scorecard is created.
			//
			// TypeSet (hashed on "id") makes the ordering in config irrelevant —
			// Terraform compares by content so reordering blocks does not trigger
			// a diff. The ForceNew flag ensures any actual content change
			// (add, remove, or modify a level) causes resource recreation.
			"progress_levels": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Description: "Progress levels for this scorecard (e.g. Red/Amber/Green). " +
					"Defined at create time — changes require resource recreation. " +
					"If omitted, the API applies org defaults.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A unique identifier for this level (e.g. 'red', 'green').",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Display name of the level.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of what this level means.",
						},
						"hex_color_code": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "Hex colour for this level, e.g. '#FF0000'.",
							ValidateFunc: validation.StringLenBetween(4, 9),
						},
					},
				},
				// Hash on "id" so reordering blocks in config never produces a diff.
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["id"].(string))
				},
			},
			// rule_ids lists the rule entity GUIDs attached to this scorecard via
			// its auto-created rules collection. Rules are independent entities
			// (newrelic_scorecard_rule) that can be shared across scorecards.
			"rule_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Description: "Set of newrelic_scorecard_rule GUIDs to attach to this scorecard. " +
					"Rules must exist before being attached. Removing a GUID detaches the rule " +
					"without deleting it — rules are standalone entities.",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			// organization_id is resolved automatically from the provider configuration.
			// Customers should never supply it — it is fetched and stored as Computed.
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NGEP organization UUID. Resolved automatically from the provider account.",
			},
			// Computed: the auto-created rules collection ID, exposed so callers
			// can reference it if needed.
			"rules_collection_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "GUID of the auto-created rules collection. Read-only.",
			},
		},
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Create
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// organization_id is always resolved automatically — customers never supply it.
	orgID, err := getOrganizationID(ctx, providerConfig, "")
	if err != nil {
		return diag.FromErr(err)
	}

	input := scorecards.EntityManagementScorecardEntityCreateInput{
		Name: d.Get("name").(string),
		Scope: scorecards.EntityManagementScopedReferenceInput{
			ID:   orgID,
			Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION,
		},
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("tags"); ok {
		input.Tags = expandNGEPTags(v.(*schema.Set).List())
	}
	if v, ok := d.GetOk("progress_levels"); ok {
		input.ProgressLevels = expandProgressLevels(v.(*schema.Set).List())
	}

	result, err := client.Scorecards.EntityManagementCreateScorecard(input)
	if err != nil {
		return diag.FromErr(err)
	}

	scorecardID := result.Entity.ID
	rulesColID := result.Entity.Rules.ID

	log.Printf("[INFO] Created NGEP scorecard %s (rules_collection=%s)", scorecardID, rulesColID)

	d.SetId(scorecardID)
	_ = d.Set("organization_id", orgID)
	_ = d.Set("rules_collection_id", rulesColID)

	// Attach any declared rules immediately after create.
	if wantedRules := expandRuleIDsFromSet(d.Get("rule_ids").(*schema.Set)); len(wantedRules) > 0 {
		if _, err := client.Scorecards.EntityManagementAddCollectionMembers(rulesColID, wantedRules); err != nil {
			if strings.Contains(err.Error(), "already belongs to collection") || strings.Contains(err.Error(), "already exists in one collection") {
				return diag.Errorf(
					"cannot attach one or more rules to scorecard %s:\n\n"+
						"  • Each scorecard rule can only belong to ONE scorecard at a time.\n"+
						"  • One or more of the rule_ids you specified is already attached to a different scorecard.\n"+
						"  • To reuse a rule, first remove it from its current scorecard (remove its GUID from that scorecard's rule_ids).\n\n"+
						"API error: %v", scorecardID, err)
			}
			return diag.Errorf("attaching rules to scorecard %s: %v", scorecardID, err)
		}
	}

	// Indexing gate — wait until the scorecard is visible in entitySearch.
	if err := waitForNGEPEntityIndexed(ctx, &client.Entities, scorecardID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicScorecardRead(ctx, d, meta)
}

// ──────────────────────────────────────────────────────────────────────────────
// Read
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

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

	sc, ok := (*entityIface).(*scorecards.EntityManagementScorecardEntity)
	if !ok {
		return diag.Errorf("entity %s is not a ScorecardEntity", d.Id())
	}

	_ = d.Set("name", sc.Name)
	_ = d.Set("description", sc.Description)
	_ = d.Set("organization_id", sc.Scope.ID)
	_ = d.Set("tags", flattenNGEPTags(sc.Tags))
	_ = d.Set("rules_collection_id", sc.Rules.ID)

	// progress_levels: always read from the entity. ForceNew on the schema
	// ensures Terraform will never call Update with a progress_levels change
	// (it would trigger recreation instead), so we never hit the NGEP backend
	// bug that blocks updates to progress levels via entityManagementUpdateScorecard.
	_ = d.Set("progress_levels", flattenProgressLevels(sc.ProgressLevels))

	// Read attached rules from the rules collection.
	ruleGUIDs, err := readScorecardRuleGUIDs(ctx, &client.Scorecards, sc.Rules.ID)
	if err != nil {
		log.Printf("[WARN] Could not read rules collection for scorecard %s: %v", d.Id(), err)
	} else {
		_ = d.Set("rule_ids", ruleGUIDs)
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	if d.HasChangesExcept("rule_ids") {
		upd := scorecards.EntityManagementScorecardEntityUpdateInput{}

		// Always include description — it has no omitempty in the API input type,
		// so guarding behind HasChange would send "" and clear it on any tag/name update.
		upd.Description = d.Get("description").(string)

		if d.HasChange("name") {
			upd.Name = d.Get("name").(string)
		}
		if d.HasChange("tags") {
			userTags := expandNGEPTags(d.Get("tags").(*schema.Set).List())
			sysTags := fetchEntitySystemTags(ctx, &client.Scorecards, d.Id())
			upd.Tags = mergeWithSystemTags(userTags, sysTags)
		}

		if _, err := client.Scorecards.EntityManagementUpdateScorecard(d.Id(), upd); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("rule_ids") {
		rulesColID := d.Get("rules_collection_id").(string)
		oldRaw, newRaw := d.GetChange("rule_ids")
		toAdd, toRemove := stringSetDelta(
			expandRuleIDsFromSet(oldRaw.(*schema.Set)),
			expandRuleIDsFromSet(newRaw.(*schema.Set)),
		)
		if len(toAdd) > 0 {
			if _, err := client.Scorecards.EntityManagementAddCollectionMembers(rulesColID, toAdd); err != nil {
				if strings.Contains(err.Error(), "already belongs to collection") {
					return diag.Errorf(
						"attaching rules to scorecard %s: one or more rules are already attached to another scorecard — "+
							"each rule can only belong to one scorecard at a time: %v", d.Id(), err)
				}
				return diag.Errorf("attaching rules to scorecard %s: %v", d.Id(), err)
			}
		}
		if len(toRemove) > 0 {
			if _, err := client.Scorecards.EntityManagementRemoveCollectionMembers(rulesColID, toRemove); err != nil {
				return diag.Errorf("detaching rules from scorecard %s: %v", d.Id(), err)
			}
		}
	}

	return resourceNewRelicScorecardRead(ctx, d, meta)
}

// ──────────────────────────────────────────────────────────────────────────────
// Delete
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Deleting NGEP scorecard %s", d.Id())
	if _, err := client.Scorecards.EntityManagementDelete(d.Id()); err != nil {
		var notFound *nrErrors.NotFound
		if errors.As(err, &notFound) {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}
