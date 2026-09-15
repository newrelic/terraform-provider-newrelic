package newrelic

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/scorecards"
)

// runIntervalAllowedMinutes is the complete list of values the NGEP API
// accepts for runInterval (minutes). Any other value is rejected.
var runIntervalAllowedMinutes = []int{60, 360, 720, 1440, 4320}

func resourceNewRelicScorecardRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicScorecardRuleCreate,
		ReadContext:   resourceNewRelicScorecardRuleRead,
		UpdateContext: resourceNewRelicScorecardRuleUpdate,
		DeleteContext: resourceNewRelicScorecardRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the scorecard rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A description of the rule.",
			},
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "The NRQL query. Must be FACET-ed and alias the result as 'score'.",
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
			"impact_weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Weight of this rule's impact on the overall scorecard score.",
				ValidateFunc: validation.IntAtLeast(0),
			},
			// progress_level references an ID defined in the parent scorecard's
			// progress_levels. Leave empty for no level association.
			"progress_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The progress level ID from the parent scorecard this rule maps to (e.g. 'red').",
			},
			// run_interval controls how often the rule is evaluated.
			// Allowed values: 60, 360, 720, 1440, 4320 (minutes).
			// Do NOT use the deprecated 'schedule' field.
			"run_interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "Evaluation frequency in minutes. Must be one of: 60, 360, 720, 1440, 4320.",
				ValidateFunc: validation.IntInSlice(runIntervalAllowedMinutes),
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Tags in 'key:value1,value2' format.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"organization_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The NGEP organization UUID. Auto-fetched if omitted.",
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

	orgID, err := getOrganizationID(ctx, providerConfig, d.Get("organization_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	input := scorecards.EntityManagementScorecardRuleEntityCreateInput{
		Name:    d.Get("name").(string),
		Enabled: d.Get("enabled").(bool),
		Scope: scorecards.EntityManagementScopedReferenceInput{
			ID:   orgID,
			Type: scorecards.EntityManagementEntityScopeTypes.ORGANIZATION,
		},
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("nrql_engine"); ok {
		input.NRQLEngine = expandNRQLEngineCreate(v.([]interface{}))
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
		input.Tags = expandNGEPTags(v.([]interface{}))
	}

	result, err := client.Scorecards.EntityManagementCreateScorecardRule(input)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.Entity.ID)
	_ = d.Set("organization_id", orgID)
	log.Printf("[INFO] Created NGEP scorecard rule %s", result.Entity.ID)

	return resourceNewRelicScorecardRuleRead(ctx, d, meta)
}

// ──────────────────────────────────────────────────────────────────────────────
// Read
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	rule, ok := (*entityIface).(*scorecards.EntityManagementScorecardRuleEntity)
	if !ok {
		return diag.Errorf("entity %s is not a ScorecardRuleEntity", d.Id())
	}

	_ = d.Set("name", rule.Name)
	_ = d.Set("description", rule.Description)
	_ = d.Set("enabled", rule.Enabled)
	_ = d.Set("tags", flattenNGEPTags(rule.Tags))
	_ = d.Set("organization_id", rule.Scope.ID)

	if rule.ImpactWeight != 0 {
		_ = d.Set("impact_weight", rule.ImpactWeight)
	}
	if rule.ProgressLevel != "" {
		_ = d.Set("progress_level", rule.ProgressLevel)
	}
	if rule.RunInterval != 0 {
		_ = d.Set("run_interval", rule.RunInterval)
	}

	_ = d.Set("nrql_engine", flattenNRQLEngine(rule.NRQLEngine))

	return nil
}

// ──────────────────────────────────────────────────────────────────────────────
// Update
// ──────────────────────────────────────────────────────────────────────────────

func resourceNewRelicScorecardRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	upd := scorecards.EntityManagementScorecardRuleEntityUpdateInput{}
	if d.HasChange("name") {
		upd.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		upd.Description = d.Get("description").(string)
	}
	if d.HasChange("enabled") {
		upd.Enabled = d.Get("enabled").(bool)
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
		upd.Tags = expandNGEPTags(d.Get("tags").([]interface{}))
	}

	if _, err := client.Scorecards.EntityManagementUpdateScorecardRule(d.Id(), upd); err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicScorecardRuleRead(ctx, d, meta)
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
