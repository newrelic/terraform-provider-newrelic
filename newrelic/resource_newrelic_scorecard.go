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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the scorecard.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the scorecard.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags in 'key:value1,value2' format.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			// progress_levels are set at create time only. The NGEP API currently
			// rejects updates to progress levels due to a backend validation bug.
			// Define them once when the scorecard is created.
			"progress_levels": {
				Type:     schema.TypeList,
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
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The NGEP organization UUID. Auto-fetched if omitted.",
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

	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
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
		input.Tags = expandNGEPTags(v.([]interface{}))
	}
	if v, ok := d.GetOk("progress_levels"); ok {
		input.ProgressLevels = expandProgressLevels(v.([]interface{}))
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
			return diag.Errorf("attaching rules to scorecard %s: %v", scorecardID, err)
		}
	}

	// Indexing gate — same as team resource.
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		res, err := client.Entities.GetEntitySearchByQueryWithContext(
			ctx,
			entities.EntitySearchOptions{},
			fmt.Sprintf("id IN ('%s')", scorecardID),
			[]entities.EntitySearchSortCriteria{},
		)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("entitySearch for scorecard %s: %w", scorecardID, err))
		}
		if res == nil || len(res.Results.Entities) == 0 {
			return resource.RetryableError(fmt.Errorf("scorecard %s not yet indexed", scorecardID))
		}
		return nil
	})
	if retryErr != nil {
		return diag.FromErr(retryErr)
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

	sc, ok := (*entityIface).(*scorecards.EntityManagementScorecardEntity)
	if !ok {
		return diag.Errorf("entity %s is not a ScorecardEntity", d.Id())
	}

	_ = d.Set("name", sc.Name)
	_ = d.Set("description", sc.Description)
	_ = d.Set("organization_id", sc.Scope.ID)
	_ = d.Set("tags", flattenNGEPTags(sc.Tags))
	_ = d.Set("rules_collection_id", sc.Rules.ID)

	// progress_levels: ForceNew so only set from create result; on Read we
	// preserve what's in state to avoid spurious diffs from the update bug.
	if _, alreadySet := d.GetOk("progress_levels"); !alreadySet {
		_ = d.Set("progress_levels", flattenProgressLevels(sc.ProgressLevels))
	}

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
		if d.HasChange("name") {
			upd.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			upd.Description = d.Get("description").(string)
		}
		if d.HasChange("tags") {
			upd.Tags = expandNGEPTags(d.Get("tags").([]interface{}))
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

// ── Helpers ───────────────────────────────────────────────────────────────────

// readScorecardRuleGUIDs pages through the scorecard's rules collection and
// returns the GUID of every ScorecardRule entity in it.
func readScorecardRuleGUIDs(ctx context.Context, client *scorecards.Scorecards, rulesColID string) ([]string, error) {
	if rulesColID == "" {
		return nil, nil
	}
	var guids []string
	err := pageCollectionItems(ctx, client, rulesColID, func(item scorecards.EntityManagementEntityInterface) {
		if r, ok := item.(*scorecards.EntityManagementScorecardRuleEntity); ok {
			guids = append(guids, r.ID)
		}
	})
	return guids, err
}

