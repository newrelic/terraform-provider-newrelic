package newrelic

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
)

// ── enum value slices (derived from Go client constants) ─────────────────────

func pathpointFlowHealthRollupValues() []string {
	return []string{
		string(pathpoint.PathPointFlowHealthRollupTypes.ALERT_CONDITIONS),
		string(pathpoint.PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP),
	}
}

func pathpointRefreshIntervalValues() []string {
	return []string{
		string(pathpoint.PathPointRefreshIntervalTypes.ONE_MINUTE),
		string(pathpoint.PathPointRefreshIntervalTypes.FIVE_MINUTES),
		string(pathpoint.PathPointRefreshIntervalTypes.TEN_MINUTES),
		string(pathpoint.PathPointRefreshIntervalTypes.FIFTEEN_MINUTES),
		string(pathpoint.PathPointRefreshIntervalTypes.THIRTY_MINUTES),
	}
}

func pathpointKpiAggregationValues() []string {
	return []string{
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.AVERAGE),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.COUNT),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.HISTOGRAM),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.MAX),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.MIN),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.PERCENTILE),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.SUM),
		string(pathpoint.PathPointKpiNRQLAggregationsTypes.UNIQUE_COUNT),
	}
}

func pathpointKpiTimeDurationValues() []string {
	return []string{
		string(pathpoint.PathPointKpiTimeDurationTypes.SEVEN_DAYS),
		string(pathpoint.PathPointKpiTimeDurationTypes.SIXTY_MINUTES),
		string(pathpoint.PathPointKpiTimeDurationTypes.SIX_HOURS),
		string(pathpoint.PathPointKpiTimeDurationTypes.THIRTY_DAYS),
		string(pathpoint.PathPointKpiTimeDurationTypes.THIRTY_MINUTES),
		string(pathpoint.PathPointKpiTimeDurationTypes.THREE_HOURS),
		string(pathpoint.PathPointKpiTimeDurationTypes.TWENTY_FOUR_HOURS),
	}
}

func pathpointSignalTypeValues() []string {
	return []string{
		string(pathpoint.PathPointSignalTypeTypes.ENTITY),
		string(pathpoint.PathPointSignalTypeTypes.ALERT),
	}
}

func pathpointStepHealthRollupValues() []string {
	return []string{
		string(pathpoint.PathPointStepHealthRollupTypes.BEST_STATUS_WINS),
		string(pathpoint.PathPointStepHealthRollupTypes.WORST_STATUS_WINS),
	}
}

func pathpointThresholdTypeValues() []string {
	return []string{
		string(pathpoint.PathPointThresholdTypeTypes.FIXED),
		string(pathpoint.PathPointThresholdTypeTypes.PERCENTAGE),
	}
}

func pathpointStageHealthRollupValues() []string {
	return []string{
		string(pathpoint.PathPointStageHealthRollupTypes.ALERT_CONDITIONS),
		string(pathpoint.PathPointStageHealthRollupTypes.AUTOMATIC_ROLL_UP),
	}
}

// ── schema definitions ────────────────────────────────────────────────────────

func pathpointFlowSchema() map[string]*schema.Schema {
	kpiQuerySelectSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"aggregation_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Aggregation function: AVERAGE, COUNT, HISTOGRAM, MAX, MIN, PERCENTILE, SUM, UNIQUE_COUNT.",
				ValidateFunc: validation.StringInSlice(pathpointKpiAggregationValues(), false),
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Optional alias for the aggregated value.",
			},
			"attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Attribute name to aggregate. Required for all functions except COUNT.",
			},
			"threshold": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Threshold used in the selected function.",
			},
		},
	}

	kpiRelativeRangeSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"since": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "How far back the KPI is evaluated.",
				ValidateFunc: validation.StringInSlice(pathpointKpiTimeDurationValues(), false),
			},
			"compare_against": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The earlier window to compare against.",
				ValidateFunc: validation.StringInSlice(pathpointKpiTimeDurationValues(), false),
			},
		},
	}

	kpiQuerySchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Data source to query from (e.g., Transaction, Metric, Log).",
			},
			"where": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional WHERE clause to filter data.",
			},
			"select": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "SELECT clause defining what to aggregate.",
				Elem:        kpiQuerySelectSchema,
			},
			"time_window": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Time window for KPI evaluation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_range": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Raw NRQL time fragment, e.g. 'SINCE 3 days ago COMPARE WITH 1 day ago'. Mutually exclusive with relative_range.",
						},
						"relative_range": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "Relative time window. Mutually exclusive with custom_range.",
							Elem:        kpiRelativeRangeSchema,
						},
					},
				},
			},
		},
	}

	kpiSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the KPI.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the KPI.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional description.",
			},
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional category to group KPIs.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Account ID this KPI belongs to. Defaults to the flow's account_id.",
			},
			"query": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "NRQL query definition for this KPI.",
				Elem:        kpiQuerySchema,
			},
			"metric_query": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "NRQL query using Metric, derived after processing event-to-metric rules. Read-only.",
			},
		},
	}

	signalSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Entity GUID of the signal.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Display name of the signal.",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Whether this GUID belongs to an entity or an alert condition: ENTITY or ALERT.",
				ValidateFunc: validation.StringInSlice(pathpointSignalTypeValues(), false),
			},
			"is_excluded": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this signal is excluded from step health calculation.",
			},
		},
	}

	stepConfigSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"health_rollup": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(pathpoint.PathPointStepHealthRollupTypes.WORST_STATUS_WINS),
				Description:  "How step health is rolled up: BEST_STATUS_WINS or WORST_STATUS_WINS.",
				ValidateFunc: validation.StringInSlice(pathpointStepHealthRollupValues(), false),
			},
			"threshold_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Whether threshold is FIXED or PERCENTAGE.",
				ValidateFunc: validation.StringInSlice(pathpointThresholdTypeValues(), false),
			},
			"threshold_value": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Numeric threshold value for step health evaluation.",
			},
		},
	}

	stepSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal step workload ID, used for updates.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the step.",
			},
			"is_excluded": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this step is excluded from level health calculation.",
			},
			"link": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional URL to an external resource.",
			},
			"scoped_accounts": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				Description: "Account IDs whose data is scoped to this step.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"entity_search_query": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Filter query used to fetch signals for this step.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter query for signals, e.g. domain='NR1' AND type='APPLICATION'.",
						},
						"is_excluded": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "When true, this query is excluded from health calculation.",
						},
					},
				},
			},
			"config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Health evaluation configuration for this step.",
				Elem:        stepConfigSchema,
			},
			"signals": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Description: "Entity signals associated with this step.",
				Elem:        signalSchema,
			},
		},
	}

	levelSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal level workload ID, used for updates.",
			},
			"steps": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				MaxItems:    50,
				Description: "Ordered list of steps within this level.",
				Elem:        stepSchema,
			},
		},
	}

	relatedSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this stage acts as a source to other stages.",
			},
			"target": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this stage acts as a target to other stages.",
			},
		},
	}

	stageSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Internal stage workload ID, used for updates.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the stage.",
			},
			"health_rollup": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      string(pathpoint.PathPointStageHealthRollupTypes.AUTOMATIC_ROLL_UP),
				Description:  "Health rollup strategy: ALERT_CONDITIONS or AUTOMATIC_ROLL_UP.",
				ValidateFunc: validation.StringInSlice(pathpointStageHealthRollupValues(), false),
			},
			"is_excluded": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this stage is excluded from flow health calculation.",
			},
			"link": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional URL to an external resource.",
			},
			"related": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Relationship role of this stage within the flow.",
				Elem:        relatedSchema,
			},
			"stage_kpis": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Description: "KPIs tracked at the stage level.",
				Elem:        kpiSchema,
			},
			"levels": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				MaxItems:    50,
				Description: "Ordered list of levels within this stage.",
				Elem:        levelSchema,
			},
		},
	}

	return map[string]*schema.Schema{
		"account_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The New Relic account ID that owns this Pathpoint flow.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Display name of the Pathpoint flow.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional description of the flow.",
		},
		"category": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional category used to group flows (e.g. Marketing, Checkout).",
		},
		"health_rollup": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(pathpoint.PathPointFlowHealthRollupTypes.AUTOMATIC_ROLL_UP),
			Description:  "Health rollup strategy: ALERT_CONDITIONS or AUTOMATIC_ROLL_UP.",
			ValidateFunc: validation.StringInSlice(pathpointFlowHealthRollupValues(), false),
		},
		"refresh_interval": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(pathpoint.PathPointRefreshIntervalTypes.THIRTY_MINUTES),
			Description:  "How often health statuses refresh: ONE_MINUTE, FIVE_MINUTES, TEN_MINUTES, FIFTEEN_MINUTES, THIRTY_MINUTES.",
			ValidateFunc: validation.StringInSlice(pathpointRefreshIntervalValues(), false),
		},
		"kpis": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			Description: "KPIs tracked at the flow level.",
			Elem:        kpiSchema,
		},
		"stages": {
			Type:        schema.TypeList,
			Optional:    true,
			MinItems:    1,
			MaxItems:    50,
			Description: "Ordered list of stages that make up this flow.",
			Elem:        stageSchema,
		},
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The entity GUID assigned to this Pathpoint flow.",
		},
		"version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Last updated timestamp, used for version control.",
		},
	}
}

// ── import helper ─────────────────────────────────────────────────────────────

// resourceNewRelicPathpointFlowImport supports two import ID formats:
//   - "<guid>"            — account_id is decoded from the GUID
//   - "<account_id>:<guid>" — explicit account_id provided by the user
func resourceNewRelicPathpointFlowImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) == 2 {
		accountID, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid account_id in import ID %q: %w", d.Id(), err)
		}
		_ = d.Set("account_id", accountID)
		d.SetId(parts[1])
	} else {
		accountID := accountIDFromGUID(d.Id())
		if accountID == 0 {
			return nil, fmt.Errorf("cannot determine account_id from GUID %q; use 'account_id:guid' import format", d.Id())
		}
		_ = d.Set("account_id", accountID)
	}

	diags := resourceNewRelicPathpointFlowRead(ctx, d, meta)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading Pathpoint flow during import: %s", diags[0].Summary)
	}
	return []*schema.ResourceData{d}, nil
}

// accountIDFromGUID decodes a New Relic entity GUID (base64 of "accountID|domain|type|id")
// and returns the account ID. Returns 0 if the GUID cannot be parsed.
func accountIDFromGUID(guid string) int {
	decoded, err := base64.StdEncoding.DecodeString(guid)
	if err != nil {
		return 0
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) < 1 {
		return 0
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}
	return id
}

// ── expand helpers (Terraform state → API input) ──────────────────────────────

func expandPathpointFlowInput(d *schema.ResourceData, accountID int) pathpoint.PathPointFlowInput {
	input := pathpoint.PathPointFlowInput{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("category"); ok {
		input.Category = v.(string)
	}
	if v := d.Get("health_rollup").(string); v != "" {
		input.HealthRollup = pathpoint.PathPointFlowHealthRollup(v)
	}
	if v := d.Get("refresh_interval").(string); v != "" {
		input.RefreshInterval = pathpoint.PathPointRefreshInterval(v)
	}
	if v, ok := d.GetOk("kpis"); ok {
		input.Kpis = expandPathpointKpiInputs(v.([]interface{}), accountID)
	}
	if v, ok := d.GetOk("stages"); ok {
		input.Stages = expandPathpointStageInputs(v.([]interface{}), accountID)
	}

	return input
}

func expandPathpointFlowUpdateInput(d *schema.ResourceData, version nrtime.EpochMilliseconds, accountID int) pathpoint.PathPointFlowUpdateInput {
	input := pathpoint.PathPointFlowUpdateInput{
		Name:    d.Get("name").(string),
		Version: version,
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("category"); ok {
		input.Category = v.(string)
	}
	if v := d.Get("health_rollup").(string); v != "" {
		input.HealthRollup = pathpoint.PathPointFlowHealthRollup(v)
	}
	if v := d.Get("refresh_interval").(string); v != "" {
		input.RefreshInterval = pathpoint.PathPointRefreshInterval(v)
	}
	if v, ok := d.GetOk("kpis"); ok {
		oldRaw, _ := d.GetChange("kpis")
		input.Kpis = expandPathpointKpiUpdateInputsResolved(v.([]interface{}), oldRaw.([]interface{}), accountID)
	}
	if v, ok := d.GetOk("stages"); ok {
		oldRaw, _ := d.GetChange("stages")
		input.Stages = expandPathpointStageUpdateInputsResolved(v.([]interface{}), oldRaw.([]interface{}), accountID)
	}

	return input
}

// expandPathpointKpiInputs builds KPI inputs for a CREATE call. No existing IDs are
// available yet, so each entry is populated purely from Terraform config.
// Contrast with expandPathpointKpiUpdateInputsResolved, which additionally carries over
// API-assigned IDs from prior state so the API can correlate old and new entries.
func expandPathpointKpiInputs(raw []interface{}, accountID int) []pathpoint.PathPointKpiInput {
	inputs := make([]pathpoint.PathPointKpiInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		kpi := pathpoint.PathPointKpiInput{
			Name: m["name"].(string),
		}
		if v, ok := m["description"].(string); ok && v != "" {
			kpi.Description = v
		}
		if v, ok := m["category"].(string); ok && v != "" {
			kpi.Category = v
		}
		if v, ok := m["account_id"].(int); ok && v != 0 {
			kpi.AccountID = v
		} else {
			kpi.AccountID = accountID
		}
		if q, ok := m["query"].([]interface{}); ok && len(q) > 0 {
			kpi.Query = expandPathpointKpiNRQLInput(q[0].(map[string]interface{}))
		}
		inputs = append(inputs, kpi)
	}
	return inputs
}

// expandPathpointKpiUpdateInputsResolved builds KPI update inputs for an UPDATE call.
// It resolves API-assigned IDs from oldRaw so the API can match existing entries.
// Matching strategy: name first (stable identity across renames), then position as a
// fallback (handles additions or deletions). The two-phase claim maps (claimedByName /
// claimedByPos) prevent the same old entry from being assigned to multiple new entries.
// Unlike expandPathpointKpiInputs (used on create), both old and new raw state are required.
func expandPathpointKpiUpdateInputsResolved(newRaw, oldRaw []interface{}, accountID int) []pathpoint.PathPointKpiUpdateInput {
	type oldInfo struct {
		id string
	}
	nameToOld := make(map[string]*oldInfo, len(oldRaw))
	oldOrdered := make([]*oldInfo, 0, len(oldRaw))
	claimedByName := make(map[string]bool, len(oldRaw))
	for _, r := range oldRaw {
		m := r.(map[string]interface{})
		info := &oldInfo{id: m["id"].(string)}
		nameToOld[m["name"].(string)] = info
		oldOrdered = append(oldOrdered, info)
	}
	claimedByPos := make([]bool, len(oldOrdered))

	inputs := make([]pathpoint.PathPointKpiUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		name := m["name"].(string)
		kpi := pathpoint.PathPointKpiUpdateInput{Name: name}

		// Match by name first, then fall back to position.
		if info, ok := nameToOld[name]; ok && !claimedByName[name] {
			if info.id != "" {
				kpi.ID = info.id
			}
			claimedByName[name] = true
		} else if i < len(oldOrdered) && !claimedByPos[i] {
			if oldOrdered[i].id != "" {
				kpi.ID = oldOrdered[i].id
			}
			claimedByPos[i] = true
		}

		if v, ok := m["description"].(string); ok && v != "" {
			kpi.Description = v
		}
		if v, ok := m["category"].(string); ok && v != "" {
			kpi.Category = v
		}
		if v, ok := m["account_id"].(int); ok && v != 0 {
			kpi.AccountID = v
		} else {
			kpi.AccountID = accountID
		}
		if q, ok := m["query"].([]interface{}); ok && len(q) > 0 {
			kpi.Query = expandPathpointKpiNRQLInput(q[0].(map[string]interface{}))
		}
		inputs = append(inputs, kpi)
	}
	return inputs
}

func expandPathpointKpiNRQLInput(m map[string]interface{}) pathpoint.PathPointKpiNRQLInput {
	q := pathpoint.PathPointKpiNRQLInput{
		From: pathpoint.NRQL(m["from"].(string)),
	}
	if v, ok := m["where"].(string); ok && v != "" {
		q.Where = v
	}
	if sel, ok := m["select"].([]interface{}); ok && len(sel) > 0 {
		q.Select = expandPathpointKpiNRQLSelectInput(sel[0].(map[string]interface{}))
	}
	if tw, ok := m["time_window"].([]interface{}); ok && len(tw) > 0 {
		twVal := expandPathpointKpiTimeWindowInput(tw[0].(map[string]interface{}))
		if twVal != nil {
			q.TimeWindow = twVal
		}
	}
	return q
}

func expandPathpointKpiNRQLSelectInput(m map[string]interface{}) pathpoint.PathPointKpiNRQLSelectInput {
	s := pathpoint.PathPointKpiNRQLSelectInput{
		AggregationType: pathpoint.PathPointKpiNRQLAggregations(m["aggregation_type"].(string)),
	}
	if v, ok := m["alias"].(string); ok && v != "" {
		s.Alias = v
	}
	if v, ok := m["attribute"].(string); ok && v != "" {
		s.Attribute = v
	}
	if v, ok := m["threshold"].(float64); ok && v != 0 {
		s.Threshold = v
	}
	return s
}

func expandPathpointKpiTimeWindowInput(m map[string]interface{}) *pathpoint.PathPointKpiTimeWindowInput {
	tw := &pathpoint.PathPointKpiTimeWindowInput{}
	hasContent := false
	if v, ok := m["custom_range"].(string); ok && v != "" {
		tw.CustomRange = pathpoint.NRQL(v)
		hasContent = true
	}
	if rr, ok := m["relative_range"].([]interface{}); ok && len(rr) > 0 {
		rrm := rr[0].(map[string]interface{})
		if since, ok := rrm["since"].(string); ok && since != "" {
			rel := &pathpoint.PathPointKpiTimeWindowRelativeRangeInput{
				Since: pathpoint.PathPointKpiTimeDuration(since),
			}
			if v, ok := rrm["compare_against"].(string); ok && v != "" {
				rel.CompareAgainst = pathpoint.PathPointKpiTimeDuration(v)
			}
			tw.RelativeRange = rel
			hasContent = true
		}
	}
	if !hasContent {
		return nil
	}
	return tw
}

// expandPathpointStageInputs builds stage inputs for a CREATE call.
// Nested stage_kpis and levels are also expanded without IDs.
// Contrast with expandPathpointStageUpdateInputsResolved, which resolves IDs at every
// nesting level (stage → stage KPIs → levels → steps).
func expandPathpointStageInputs(raw []interface{}, accountID int) []pathpoint.PathPointStageInput {
	stages := make([]pathpoint.PathPointStageInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		stage := pathpoint.PathPointStageInput{
			Name: m["name"].(string),
		}
		if v, ok := m["health_rollup"].(string); ok && v != "" {
			stage.HealthRollup = pathpoint.PathPointStageHealthRollup(v)
		}
		if v, ok := m["is_excluded"].(bool); ok {
			stage.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			stage.Link = v
		}
		if rel, ok := m["related"].([]interface{}); ok && len(rel) > 0 {
			stage.Related = expandPathpointRelatedInput(rel[0].(map[string]interface{}))
		}
		if kpis, ok := m["stage_kpis"].([]interface{}); ok {
			stage.StageKpis = expandPathpointKpiInputs(kpis, accountID)
		}
		if levels, ok := m["levels"].([]interface{}); ok {
			stage.Levels = expandPathpointLevelInputs(levels)
		}
		stages = append(stages, stage)
	}
	return stages
}

// expandPathpointStageUpdateInputsResolved builds stage update inputs for an UPDATE call.
// It resolves API-assigned IDs for the stage itself and propagates the matching old-state
// context (oldLevels, oldStageKpis) into the nested expanders so each level of the
// hierarchy can independently carry over its own IDs.
// Matching strategy: name first, then position — same as KPIs.
// Unlike expandPathpointStageInputs (used on create), both old and new raw state are required.
func expandPathpointStageUpdateInputsResolved(newRaw, oldRaw []interface{}, accountID int) []pathpoint.PathPointStageUpdateInput {
	type oldInfo struct {
		id        string
		levels    []interface{}
		stageKpis []interface{}
	}
	nameToOld := make(map[string]*oldInfo, len(oldRaw))
	oldOrdered := make([]*oldInfo, 0, len(oldRaw))
	claimedByName := make(map[string]bool, len(oldRaw))
	for _, r := range oldRaw {
		m := r.(map[string]interface{})
		levels, _ := m["levels"].([]interface{})
		stageKpis, _ := m["stage_kpis"].([]interface{})
		info := &oldInfo{id: m["id"].(string), levels: levels, stageKpis: stageKpis}
		nameToOld[m["name"].(string)] = info
		oldOrdered = append(oldOrdered, info)
	}
	claimedByPos := make([]bool, len(oldOrdered))

	stages := make([]pathpoint.PathPointStageUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		name := m["name"].(string)
		stage := pathpoint.PathPointStageUpdateInput{Name: name}

		// Match by name first, then fall back to position.
		var matched *oldInfo
		if info, ok := nameToOld[name]; ok && !claimedByName[name] {
			matched = info
			claimedByName[name] = true
		} else if i < len(oldOrdered) && !claimedByPos[i] {
			matched = oldOrdered[i]
			claimedByPos[i] = true
		}

		var oldLevels, oldStageKpis []interface{}
		if matched != nil {
			if matched.id != "" {
				stage.ID = matched.id
			}
			oldLevels = matched.levels
			oldStageKpis = matched.stageKpis
		}

		if v, ok := m["health_rollup"].(string); ok && v != "" {
			stage.HealthRollup = pathpoint.PathPointStageHealthRollup(v)
		}
		if v, ok := m["is_excluded"].(bool); ok {
			stage.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			stage.Link = v
		}
		if rel, ok := m["related"].([]interface{}); ok && len(rel) > 0 {
			stage.Related = expandPathpointRelatedInput(rel[0].(map[string]interface{}))
		}
		if kpis, ok := m["stage_kpis"].([]interface{}); ok {
			stage.StageKpis = expandPathpointKpiUpdateInputsResolved(kpis, oldStageKpis, accountID)
		}
		if levels, ok := m["levels"].([]interface{}); ok {
			stage.Levels = expandPathpointLevelUpdateInputsResolved(levels, oldLevels)
		}
		stages = append(stages, stage)
	}
	return stages
}

// expandPathpointLevelUpdateInputsResolved builds level update inputs for an UPDATE call.
// Levels have no display name, so name-based matching is not possible. Instead, each new
// level is matched to the old level whose step names overlap the most (best-overlap wins).
// Falls back to positional matching when no overlap is found (e.g. all steps were renamed).
// Contrast with expandPathpointLevelInputs (used on create), which simply iterates by
// position and never needs to resolve IDs.
func expandPathpointLevelUpdateInputsResolved(newRaw, oldRaw []interface{}) []pathpoint.PathPointLevelUpdateInput {
	type oldInfo struct {
		id        string
		steps     []interface{}
		stepNames map[string]bool
	}
	oldLevels := make([]oldInfo, 0, len(oldRaw))
	for _, r := range oldRaw {
		m := r.(map[string]interface{})
		steps, _ := m["steps"].([]interface{})
		names := make(map[string]bool, len(steps))
		for _, s := range steps {
			names[s.(map[string]interface{})["name"].(string)] = true
		}
		oldLevels = append(oldLevels, oldInfo{id: m["id"].(string), steps: steps, stepNames: names})
	}
	claimed := make([]bool, len(oldLevels))

	levels := make([]pathpoint.PathPointLevelUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		level := pathpoint.PathPointLevelUpdateInput{}
		newSteps, _ := m["steps"].([]interface{})

		// Match by step-name overlap first; fall back to position if no overlap found.
		bestIdx, bestOverlap := -1, 0
		for j := range oldLevels {
			if claimed[j] {
				continue
			}
			overlap := 0
			for _, s := range newSteps {
				if oldLevels[j].stepNames[s.(map[string]interface{})["name"].(string)] {
					overlap++
				}
			}
			if overlap > bestOverlap {
				bestOverlap = overlap
				bestIdx = j
			}
		}
		if bestIdx < 0 && i < len(oldLevels) && !claimed[i] {
			bestIdx = i
		}

		var oldSteps []interface{}
		if bestIdx >= 0 {
			claimed[bestIdx] = true
			if oldLevels[bestIdx].id != "" {
				level.ID = oldLevels[bestIdx].id
			}
			oldSteps = oldLevels[bestIdx].steps
		}

		level.Steps = expandPathpointStepUpdateInputsResolved(newSteps, oldSteps)
		levels = append(levels, level)
	}
	return levels
}

// expandPathpointStepUpdateInputsResolved builds step update inputs for an UPDATE call.
// It resolves API-assigned step IDs from oldRaw so the API can match existing entries.
// Matching strategy: step name first, then position — same as KPIs and stages.
// Unlike expandPathpointStepInputs (used on create), old state is required to resolve IDs.
func expandPathpointStepUpdateInputsResolved(newRaw, oldRaw []interface{}) []pathpoint.PathPointStepUpdateInput {
	type oldInfo struct {
		id  string
		pos int
	}
	nameToOld := make(map[string]*oldInfo, len(oldRaw))
	oldOrdered := make([]*oldInfo, 0, len(oldRaw))
	claimedByName := make(map[string]bool, len(oldRaw))
	for i, r := range oldRaw {
		m := r.(map[string]interface{})
		info := &oldInfo{id: m["id"].(string), pos: i}
		nameToOld[m["name"].(string)] = info
		oldOrdered = append(oldOrdered, info)
	}
	claimedByPos := make([]bool, len(oldOrdered))

	steps := make([]pathpoint.PathPointStepUpdateInput, 0, len(newRaw))
	for i, r := range newRaw {
		m := r.(map[string]interface{})
		name := m["name"].(string)
		step := pathpoint.PathPointStepUpdateInput{Name: name}

		// Match by name first, then fall back to position.
		if info, ok := nameToOld[name]; ok && !claimedByName[name] {
			if info.id != "" {
				step.ID = info.id
			}
			claimedByName[name] = true
		} else if i < len(oldOrdered) && !claimedByPos[i] {
			if oldOrdered[i].id != "" {
				step.ID = oldOrdered[i].id
			}
			claimedByPos[i] = true
		}
		if v, ok := m["is_excluded"].(bool); ok {
			step.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			step.Link = v
		}
		if v, ok := m["scoped_accounts"].([]interface{}); ok {
			for _, id := range v {
				step.ScopedAccounts = append(step.ScopedAccounts, id.(int))
			}
		}
		if esq, ok := m["entity_search_query"].([]interface{}); ok && len(esq) > 0 {
			step.EntitySearchQuery = expandPathpointSignalQueryInput(esq[0].(map[string]interface{}))
		}
		if cfg, ok := m["config"].([]interface{}); ok && len(cfg) > 0 {
			step.Config = expandPathpointStepStatusThresholdInput(cfg[0].(map[string]interface{}))
		}
		if sigs, ok := m["signals"].([]interface{}); ok {
			step.Signals = expandPathpointSignalInputs(sigs)
		}
		steps = append(steps, step)
	}
	return steps
}

// expandPathpointLevelInputs builds level inputs for a CREATE call.
// Levels have no display name; they are ordered by position and expanded without IDs.
// Contrast with expandPathpointLevelUpdateInputsResolved, which uses step-name overlap
// to match new levels to old ones and carry over API-assigned IDs.
func expandPathpointLevelInputs(raw []interface{}) []pathpoint.PathPointLevelInput {
	levels := make([]pathpoint.PathPointLevelInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		level := pathpoint.PathPointLevelInput{}
		if steps, ok := m["steps"].([]interface{}); ok {
			level.Steps = expandPathpointStepInputs(steps)
		}
		levels = append(levels, level)
	}
	return levels
}

// expandPathpointStepInputs builds step inputs for a CREATE call.
// All fields are populated from Terraform config; no existing IDs are available.
// Contrast with expandPathpointStepUpdateInputsResolved, which carries over API-assigned IDs.
func expandPathpointStepInputs(raw []interface{}) []pathpoint.PathPointStepInput {
	steps := make([]pathpoint.PathPointStepInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		step := pathpoint.PathPointStepInput{
			Name: m["name"].(string),
		}
		if v, ok := m["is_excluded"].(bool); ok {
			step.IsExcluded = v
		}
		if v, ok := m["link"].(string); ok && v != "" {
			step.Link = v
		}
		if v, ok := m["scoped_accounts"].([]interface{}); ok {
			for _, id := range v {
				step.ScopedAccounts = append(step.ScopedAccounts, id.(int))
			}
		}
		if esq, ok := m["entity_search_query"].([]interface{}); ok && len(esq) > 0 {
			step.EntitySearchQuery = expandPathpointSignalQueryInput(esq[0].(map[string]interface{}))
		}
		if cfg, ok := m["config"].([]interface{}); ok && len(cfg) > 0 {
			step.Config = expandPathpointStepStatusThresholdInput(cfg[0].(map[string]interface{}))
		}
		if sigs, ok := m["signals"].([]interface{}); ok {
			step.Signals = expandPathpointSignalInputs(sigs)
		}
		steps = append(steps, step)
	}
	return steps
}

func expandPathpointSignalQueryInput(m map[string]interface{}) *pathpoint.PathPointSignalQueryInput {
	q := &pathpoint.PathPointSignalQueryInput{
		Query: m["query"].(string),
	}
	if v, ok := m["is_excluded"].(bool); ok {
		q.IsExcluded = v
	}
	return q
}

func expandPathpointStepStatusThresholdInput(m map[string]interface{}) pathpoint.PathPointStepStatusThresholdInput {
	cfg := pathpoint.PathPointStepStatusThresholdInput{}
	if v, ok := m["health_rollup"].(string); ok && v != "" {
		cfg.HealthRollup = pathpoint.PathPointStepHealthRollup(v)
	}
	if v, ok := m["threshold_type"].(string); ok && v != "" {
		cfg.ThresholdType = pathpoint.PathPointThresholdType(v)
	}
	if v, ok := m["threshold_value"].(int); ok {
		cfg.ThresholdValue = v
	}
	return cfg
}

func expandPathpointSignalInputs(raw []interface{}) []pathpoint.PathPointSignalInput {
	signals := make([]pathpoint.PathPointSignalInput, 0, len(raw))
	for _, r := range raw {
		m := r.(map[string]interface{})
		sig := pathpoint.PathPointSignalInput{
			GUID: pathpoint.EntityGUID(m["guid"].(string)),
		}
		if v, ok := m["name"].(string); ok && v != "" {
			sig.Name = v
		}
		if v, ok := m["type"].(string); ok && v != "" {
			sig.Type = pathpoint.PathPointSignalType(v)
		}
		if v, ok := m["is_excluded"].(bool); ok {
			sig.IsExcluded = v
		}
		signals = append(signals, sig)
	}
	return signals
}

func expandPathpointRelatedInput(m map[string]interface{}) pathpoint.PathPointRelatedInput {
	return pathpoint.PathPointRelatedInput{
		Source: m["source"].(bool),
		Target: m["target"].(bool),
	}
}

// ── flatten helpers (API response → Terraform state) ─────────────────────────

func flattenPathpointKpis(kpis []pathpoint.PathPointKpi) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(kpis))
	for _, k := range kpis {
		m := map[string]interface{}{
			"id":           k.ID,
			"name":         k.Name,
			"description":  k.Description,
			"category":     k.Category,
			"query":        flattenPathpointKpiNRQL(k.Query),
			"metric_query": string(k.MetricQuery),
		}
		if k.AccountID != 0 {
			m["account_id"] = k.AccountID
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointKpiNRQL(q pathpoint.PathPointKpiNRQL) []map[string]interface{} {
	m := map[string]interface{}{
		"from":  q.From,
		"where": q.Where,
		"select": []map[string]interface{}{
			{
				"aggregation_type": string(q.Select.AggregationType),
				"alias":            q.Select.Alias,
				"attribute":        q.Select.Attribute,
				"threshold":        q.Select.Threshold,
			},
		},
	}
	if tw := flattenPathpointKpiTimeWindow(q.TimeWindow); len(tw) > 0 {
		m["time_window"] = tw
	}
	return []map[string]interface{}{m}
}

func flattenPathpointKpiTimeWindow(tw pathpoint.PathPointKpiTimeWindow) []map[string]interface{} {
	if tw.CustomRange == "" && tw.RelativeRange.Since == "" {
		return nil
	}
	m := map[string]interface{}{
		"custom_range": string(tw.CustomRange),
	}
	if tw.RelativeRange.Since != "" {
		rr := map[string]interface{}{
			"since":           string(tw.RelativeRange.Since),
			"compare_against": string(tw.RelativeRange.CompareAgainst),
		}
		m["relative_range"] = []map[string]interface{}{rr}
	}
	return []map[string]interface{}{m}
}

func flattenPathpointStages(stages []pathpoint.PathPointStage) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(stages))
	for _, s := range stages {
		m := map[string]interface{}{
			"id":            s.ID,
			"name":          s.Name,
			"health_rollup": string(s.HealthRollup),
			"is_excluded":   s.IsExcluded,
			"link":          s.Link,
			"stage_kpis":    flattenPathpointKpis(s.StageKpis),
			"levels":        flattenPathpointLevels(s.Levels.Items),
		}
		if s.Related.Source || s.Related.Target {
			m["related"] = []map[string]interface{}{
				{
					"source": s.Related.Source,
					"target": s.Related.Target,
				},
			}
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointLevels(levels []pathpoint.PathPointLevel) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(levels))
	for _, l := range levels {
		m := map[string]interface{}{
			"id":    l.ID,
			"steps": flattenPathpointSteps(l.Steps.Items),
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointSteps(steps []pathpoint.PathPointStep) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(steps))
	for _, s := range steps {
		scopedAccounts := make([]interface{}, 0, len(s.ScopedAccounts))
		for _, id := range s.ScopedAccounts {
			scopedAccounts = append(scopedAccounts, id)
		}
		m := map[string]interface{}{
			"id":              s.ID,
			"name":            s.Name,
			"is_excluded":     s.IsExcluded,
			"link":            s.Link,
			"scoped_accounts": scopedAccounts,
			"signals":         flattenPathpointSignals(s.Signals),
		}
		if s.EntitySearchQuery.Query != "" {
			m["entity_search_query"] = []map[string]interface{}{
				{
					"query":       s.EntitySearchQuery.Query,
					"is_excluded": s.EntitySearchQuery.IsExcluded,
				},
			}
		}
		if s.Config.HealthRollup != "" || s.Config.ThresholdType != "" || s.Config.ThresholdValue != 0 {
			m["config"] = []map[string]interface{}{
				{
					"health_rollup":   string(s.Config.HealthRollup),
					"threshold_type":  string(s.Config.ThresholdType),
					"threshold_value": s.Config.ThresholdValue,
				},
			}
		}
		result = append(result, m)
	}
	return result
}

func flattenPathpointSignals(signals []pathpoint.PathPointSignal) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(signals))
	for _, s := range signals {
		result = append(result, map[string]interface{}{
			"guid":        string(s.GUID),
			"name":        s.Name,
			"type":        string(s.Type),
			"is_excluded": s.IsExcluded,
		})
	}
	return result
}
