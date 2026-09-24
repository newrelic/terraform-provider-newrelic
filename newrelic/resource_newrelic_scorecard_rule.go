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

// customizeScorecardRuleDiff validates the nrql_engine accounts/join_accounts
// lists for semantic correctness at plan time:
//  1. No account ID ≤ 0 in either list.
//  2. No duplicates within accounts.
//  3. No duplicates within join_accounts.
//  4. No account ID appears in both accounts AND join_accounts.
func customizeScorecardRuleDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	rawEngine := d.Get("nrql_engine").([]interface{})
	if len(rawEngine) == 0 || rawEngine[0] == nil {
		return nil
	}
	m := rawEngine[0].(map[string]interface{})

	toIntSlice := func(raw []interface{}) []int {
		out := make([]int, 0, len(raw))
		for _, v := range raw {
			out = append(out, v.(int))
		}
		return out
	}

	accounts := toIntSlice(m["accounts"].([]interface{}))
	joinAccounts := toIntSlice(m["join_accounts"].([]interface{}))

	// 1. No account ID ≤ 0
	for _, id := range accounts {
		if id <= 0 {
			return fmt.Errorf("nrql_engine.accounts: account ID %d is invalid (must be > 0)", id)
		}
	}
	for _, id := range joinAccounts {
		if id <= 0 {
			return fmt.Errorf("nrql_engine.join_accounts: account ID %d is invalid (must be > 0)", id)
		}
	}

	// 2. No duplicates within accounts
	seen := make(map[int]bool, len(accounts))
	for _, id := range accounts {
		if seen[id] {
			return fmt.Errorf("nrql_engine.accounts: duplicate account ID %d", id)
		}
		seen[id] = true
	}

	// 3. No duplicates within join_accounts
	seenJoin := make(map[int]bool, len(joinAccounts))
	for _, id := range joinAccounts {
		if seenJoin[id] {
			return fmt.Errorf("nrql_engine.join_accounts: duplicate account ID %d", id)
		}
		seenJoin[id] = true
	}

	// 4. No intersection between accounts and join_accounts
	for _, id := range joinAccounts {
		if seen[id] {
			return fmt.Errorf("nrql_engine: account ID %d appears in both accounts and join_accounts — these lists must be disjoint", id)
		}
	}

	return nil
}

func resourceNewRelicScorecardRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicScorecardRuleCreate,
		ReadContext:   resourceNewRelicScorecardRuleRead,
		UpdateContext: resourceNewRelicScorecardRuleUpdate,
		DeleteContext: resourceNewRelicScorecardRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customizeScorecardRuleDiff,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of the scorecard rule.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the rule.",
			},
			// enabled uses schema.TypeBool which the Terraform SDK validates strictly:
			// only true/false/1/0/yes/no are accepted; any other value produces a
			// schema validation error before Create/Update is invoked.
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether this rule is active. Disabled rules are not evaluated.",
			},
			// nrql_engine is required — it is the rule's evaluation query.
			// Use runInterval + enabled (not schedule which is deprecated).
			"nrql_engine": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Description: "NRQL query configuration for this rule. The query must be FACET-ed " +
					"and use if(latest(...), 1, 0) AS 'score' in the SELECT clause.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The NRQL query. Must be FACET-ed and alias the result as 'score'.",
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"accounts": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Account IDs where this rule runs.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"join_accounts": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Additional account IDs to join with the query accounts.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			// TODO: Confirm the valid upper bound for impact_weight with the Scorecards
			// team — currently accepts any non-negative integer. Update the ValidateFunc
			// and the resource documentation once the valid range is confirmed.
			"impact_weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Weight of this rule's impact on the overall scorecard score.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			// progress_level references an ID defined in the parent scorecard's
			// progress_levels. The value must match a progress_levels.id in the
			// scorecard this rule is attached to. Leave empty for no level association.
			"progress_level": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The progress level ID from the parent scorecard this rule maps to (e.g. 'red'). Must match a progress_levels.id defined on the scorecard.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			// run_interval controls how often the rule is evaluated (minutes).
			// The API accepts only 60, 360, 720, or 1440 minutes.
			// Do NOT use the deprecated 'schedule' field.
			"run_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "Evaluation frequency in minutes. " +
					"Must be one of: 60 (hourly), 360 (6h), 720 (12h), 1440 (daily). " +
					"Omit to use the API default.",
				ValidateFunc: validation.IntInSlice([]int{60, 360, 720, 1440}),
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
			// organization_id is resolved automatically from the provider configuration.
			// Customers should never supply it — it is fetched and stored as Computed.
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NGEP organization UUID. Resolved automatically from the provider account.",
			},
		},
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Create
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// organization_id is always resolved automatically — customers never supply it.
	orgID, err := getOrganizationID(ctx, providerConfig, "")
	if err != nil {
		return diag.FromErr(err)
	}

	// nrql_engine is Required — use d.Get, not d.GetOk.
	input := scorecards.EntityManagementScorecardRuleEntityCreateInput{
		Name:       d.Get("name").(string),
		Enabled:    d.Get("enabled").(bool),
		NRQLEngine: expandNRQLEngineCreate(d.Get("nrql_engine").([]interface{})),
		Scope: scorecards.EntityManagementScopedReferenceInput{
			ID:   orgID,
			Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION,
		},
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("impact_weight"); ok {
		input.ImpactWeight = v.(int)
	}
	if v, ok := d.GetOk("progress_level"); ok {
		input.ProgressLevel = v.(string)
	}
	if v, ok := d.GetOk("run_interval"); ok {
		input.RunInterval = v.(int)
	}
	if v, ok := d.GetOk("tags"); ok {
		input.Tags = expandNGEPTags(v.(*schema.Set).List())
	}

	result, err := client.Scorecards.EntityManagementCreateScorecardRule(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Entity.ID)
	_ = d.Set("organization_id", orgID)
	log.Printf("[INFO] Created NGEP scorecard rule %s", result.Entity.ID)

	// Set all state from what we sent — the API is a passthrough for these fields
	// so there is no need to issue a Read call. The indexing gate below ensures the
	// entity is queryable before Terraform considers Create complete.
	_ = d.Set("name", input.Name)
	_ = d.Set("enabled", input.Enabled)
	_ = d.Set("description", input.Description)
	_ = d.Set("impact_weight", input.ImpactWeight)
	_ = d.Set("progress_level", input.ProgressLevel)
	_ = d.Set("run_interval", input.RunInterval)
	if input.NRQLEngine != nil {
		_ = d.Set("nrql_engine", flattenNRQLEngine(scorecards.EntityManagementNRQLRuleEngine{
			Query:        input.NRQLEngine.Query,
			Accounts:     input.NRQLEngine.Accounts,
			JoinAccounts: input.NRQLEngine.JoinAccounts,
		}))
	}
	// Tags: convert TagInput → Tag (identical fields) for the flatten helper.
	if len(input.Tags) > 0 {
		tagSlice := make([]scorecards.EntityManagementTag, len(input.Tags))
		for i, t := range input.Tags {
			tagSlice[i] = scorecards.EntityManagementTag(t)
		}
		_ = d.Set("tags", flattenNGEPTags(tagSlice))
	}

	// Indexing gate — wait until the rule is visible in entitySearch before
	// returning. Subsequent Read calls will then find the entity immediately.
	if err := waitForNGEPEntityIndexed(ctx, &client.Entities, result.Entity.ID, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Read
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	rule, ok := (*entityIface).(*scorecards.EntityManagementScorecardRuleEntity)
	if !ok {
		return diag.Errorf("entity %s is not a ScorecardRuleEntity", d.Id())
	}

	_ = d.Set("name", rule.Name)
	_ = d.Set("description", rule.Description)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("tags", flattenNGEPTags(rule.Tags))
	_ = d.Set("organization_id", rule.Scope.ID)

	// Always set these optional int/string fields so that clearing them (setting
	// to zero/empty) is reflected in state. Guarding with != 0 / != "" would
	// leave stale non-zero values in state when the server returns 0/"".
	_ = d.Set("impact_weight", rule.ImpactWeight)
	_ = d.Set("progress_level", rule.ProgressLevel)
	_ = d.Set("run_interval", rule.RunInterval)

	_ = d.Set("nrql_engine", flattenNRQLEngine(rule.NRQLEngine))

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	upd := scorecards.EntityManagementScorecardRuleEntityUpdateInput{}

	// Always include enabled and description — they have no omitempty in the API
	// input type, so guarding behind HasChange would send the zero value (false/"")
	// for any update that touches other fields, accidentally disabling the rule or
	// clearing its description.
	upd.Enabled = d.Get("enabled").(bool)
	upd.Description = d.Get("description").(string)

	if d.HasChange("name") {
		upd.Name = d.Get("name").(string)
	}
	if d.HasChange("nrql_engine") {
		upd.NRQLEngine = expandNRQLEngineUpdate(d.Get("nrql_engine").([]interface{}))
	}
	if d.HasChange("impact_weight") {
		upd.ImpactWeight = d.Get("impact_weight").(int)
	}
	if d.HasChange("progress_level") {
		upd.ProgressLevel = d.Get("progress_level").(string)
	}
	if d.HasChange("run_interval") {
		upd.RunInterval = d.Get("run_interval").(int)
	}
	if d.HasChange("tags") {
		userTags := expandNGEPTags(d.Get("tags").(*schema.Set).List())
		sysTags := fetchEntitySystemTags(ctx, &client.Scorecards, d.Id())
		upd.Tags = mergeWithSystemTags(userTags, sysTags)
	}

	if _, err := client.Scorecards.EntityManagementUpdateScorecardRule(d.Id(), upd); err != nil {
		return diag.FromErr(err)
	}
	// Terraform's ResourceData already tracks all changed values; no Read
	// round-trip is needed. The next plan cycle will call Read to reconcile
	// any server-side divergence.
	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Delete
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Deleting NGEP scorecard rule %s", d.Id())
	if _, err := client.Scorecards.EntityManagementDelete(d.Id()); err != nil {
		var notFound *nrErrors.NotFound
		if errors.As(err, &notFound) {
			return nil
		}
		return diag.FromErr(err)
	}
	return nil
}
