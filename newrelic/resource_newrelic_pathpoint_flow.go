package newrelic

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
)

func resourceNewRelicPathpointFlow() *schema.Resource {
	kpiQuerySelectSchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"aggregation_type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Aggregation function: AVERAGE, COUNT, HISTOGRAM, MAX, MIN, PERCENTILE, SUM, UNIQUE_COUNT.",
				ValidateFunc: validation.StringInSlice([]string{"AVERAGE", "COUNT", "HISTOGRAM", "MAX", "MIN", "PERCENTILE", "SUM", "UNIQUE_COUNT"}, false),
			},
			"alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional alias for the aggregated value.",
			},
			"attribute": {
				Type:        schema.TypeString,
				Optional:    true,
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
				ValidateFunc: validation.StringInSlice([]string{"SEVEN_DAYS", "SIXTY_MINUTES", "SIX_HOURS", "THIRTY_DAYS", "THIRTY_MINUTES", "THREE_HOURS", "TWENTY_FOUR_HOURS"}, false),
			},
			"compare_against": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The earlier window to compare against.",
				ValidateFunc: validation.StringInSlice([]string{"SEVEN_DAYS", "SIXTY_MINUTES", "SIX_HOURS", "THIRTY_DAYS", "THIRTY_MINUTES", "THREE_HOURS", "TWENTY_FOUR_HOURS"}, false),
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
				Description: "Account ID this KPI belongs to.",
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
				ValidateFunc: validation.StringInSlice([]string{"ENTITY", "ALERT"}, false),
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
				Description:  "How step health is rolled up: BEST_STATUS_WINS or WORST_STATUS_WINS.",
				ValidateFunc: validation.StringInSlice([]string{"BEST_STATUS_WINS", "WORST_STATUS_WINS"}, false),
			},
			"threshold_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Whether threshold is FIXED or PERCENTAGE.",
				ValidateFunc: validation.StringInSlice([]string{"FIXED", "PERCENTAGE"}, false),
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
				MaxItems:    1,
				Description: "Health evaluation configuration for this step.",
				Elem:        stepConfigSchema,
			},
			"signals": {
				Type:        schema.TypeList,
				Optional:    true,
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
				Description:  "Health rollup strategy: ALERT_CONDITIONS or AUTOMATIC_ROLL_UP.",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_CONDITIONS", "AUTOMATIC_ROLL_UP"}, false),
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
				Description: "KPIs tracked at the stage level.",
				Elem:        kpiSchema,
			},
			"levels": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    50,
				Description: "Ordered list of levels within this stage.",
				Elem:        levelSchema,
			},
		},
	}

	return &schema.Resource{
		CreateContext: resourceNewRelicPathpointFlowCreate,
		ReadContext:   resourceNewRelicPathpointFlowRead,
		UpdateContext: resourceNewRelicPathpointFlowUpdate,
		DeleteContext: resourceNewRelicPathpointFlowDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
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
				Description:  "Health rollup strategy: ALERT_CONDITIONS or AUTOMATIC_ROLL_UP.",
				ValidateFunc: validation.StringInSlice([]string{"ALERT_CONDITIONS", "AUTOMATIC_ROLL_UP"}, false),
			},
			"refresh_interval": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "How often health statuses refresh: ONE_MINUTE, FIVE_MINUTES, TEN_MINUTES, FIFTEEN_MINUTES, THIRTY_MINUTES.",
				ValidateFunc: validation.StringInSlice([]string{"ONE_MINUTE", "FIVE_MINUTES", "TEN_MINUTES", "FIFTEEN_MINUTES", "THIRTY_MINUTES"}, false),
			},
			"kpis": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "KPIs tracked at the flow level.",
				Elem:        kpiSchema,
			},
			"stages": {
				Type:        schema.TypeList,
				Optional:    true,
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
		},
	}
}

func resourceNewRelicPathpointFlowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	flowInput := expandPathpointFlowInput(d, accountID)
	scopeInput := pathpoint.PathPointScopeInput{
		ID:   accountID,
		Type: pathpoint.PathPointScopeTypeTypes.ACCOUNT,
	}

	result, err := client.PathPoint.PathPointCreate(flowInput, scopeInput)
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("error creating Pathpoint flow: empty response"))
	}

	d.SetId(string(result.GUID))
	log.Printf("[INFO] Created Pathpoint flow %s (GUID: %s)", result.Name, result.GUID)

	return flattenPathpointFlowResult(d, result)
}

func resourceNewRelicPathpointFlowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := pathpoint.EntityGUID(d.Id())

	// After import, account_id may not yet be in state — decode it from the GUID.
	accountID := selectAccountID(providerConfig, d)
	if accountID == 0 {
		accountID = accountIDFromGUID(d.Id())
	}

	result, err := client.PathPoint.GetFlow(accountID, guid)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("error reading Pathpoint flow (GUID: %s): %w", guid, err))
	}
	if result == nil || string(result.GUID) == "" {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", accountID)
	return flattenPathpointFlowResult(d, result)
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

func resourceNewRelicPathpointFlowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := pathpoint.EntityGUID(d.Id())

	versionStr := d.Get("version").(string)
	versionInt, _ := strconv.ParseInt(versionStr, 10, 64)
	version := nrtime.EpochMilliseconds(time.UnixMilli(versionInt))

	accountID := selectAccountID(providerConfig, d)
	updateInput := expandPathpointFlowUpdateInput(d, version, accountID)

	result, err := client.PathPoint.PathPointUpdate(guid, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}
	if result == nil {
		return diag.FromErr(fmt.Errorf("error updating Pathpoint flow: empty response"))
	}

	log.Printf("[INFO] Updated Pathpoint flow %s (GUID: %s)", result.Name, result.GUID)

	return flattenPathpointFlowResult(d, result)
}

func resourceNewRelicPathpointFlowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := pathpoint.EntityGUID(d.Id())

	_, err := client.PathPoint.PathPointDelete(guid)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleted Pathpoint flow (GUID: %s)", guid)
	return nil
}

func flattenPathpointFlowResult(d *schema.ResourceData, result *pathpoint.PathPointFlowResult) diag.Diagnostics {
	_ = d.Set("guid", string(result.GUID))
	_ = d.Set("name", result.Name)
	_ = d.Set("version", strconv.FormatInt(time.Time(result.Version).UnixMilli(), 10))
	_ = d.Set("description", result.Description)
	_ = d.Set("category", result.Category)

	if result.HealthRollup != "" {
		_ = d.Set("health_rollup", string(result.HealthRollup))
	}
	if result.RefreshInterval != "" {
		_ = d.Set("refresh_interval", string(result.RefreshInterval))
	}
	if result.Message != "" {
		log.Printf("[WARN] Pathpoint flow response message: %s", result.Message)
	}

	if err := d.Set("kpis", flattenPathpointKpis(result.Kpis)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting kpis: %w", err))
	}

	if err := d.Set("stages", flattenPathpointStages(result.Stages.Items)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting stages: %w", err))
	}

	return nil
}
