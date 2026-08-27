package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicPathpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicPathpointCreate,
		ReadContext:   resourceNewRelicPathpointRead,
		UpdateContext: resourceNewRelicPathpointUpdate,
		DeleteContext: resourceNewRelicPathpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"health_rollup": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALERT_CONDITIONS",
					"AUTOMATIC_ROLL_UP",
				}, false),
			},
			"refresh_interval": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FIFTEEN_MINUTES",
					"FIVE_MINUTES",
					"ONE_MINUTE",
					"TEN_MINUTES",
					"THIRTY_MINUTES",
				}, false),
			},
			"kpis": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNewRelicPathpointKpiSchema(),
			},
			"stages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNewRelicPathpointStageSchema(),
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flow_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"excluded_kpis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNewRelicPathpointKpiSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kpi_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointKpiNRQLSchema(),
			},
		},
	}
}

func resourceNewRelicPathpointKpiNRQLSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"from": {
				Type:     schema.TypeString,
				Required: true,
			},
			"where": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"select": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointKpiNRQLSelectSchema(),
			},
			"time_window": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointKpiTimeWindowSchema(),
			},
		},
	}
}

func resourceNewRelicPathpointKpiNRQLSelectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"aggregation_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AVERAGE",
					"COUNT",
					"HISTOGRAM",
					"MAX",
					"MIN",
					"PERCENTILE",
					"SUM",
					"UNIQUE_COUNT",
				}, false),
			},
			"alias": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
		},
	}
}

func resourceNewRelicPathpointKpiTimeWindowSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"relative_range": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointKpiTimeWindowRelativeRangeSchema(),
			},
		},
	}
}

func resourceNewRelicPathpointKpiTimeWindowRelativeRangeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"since": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SEVEN_DAYS",
					"SIXTY_MINUTES",
					"SIX_HOURS",
					"THIRTY_DAYS",
					"THIRTY_MINUTES",
					"THREE_HOURS",
					"TWENTY_FOUR_HOURS",
				}, false),
			},
			"compare_against": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SEVEN_DAYS",
					"SIXTY_MINUTES",
					"SIX_HOURS",
					"THIRTY_DAYS",
					"THIRTY_MINUTES",
					"THREE_HOURS",
					"TWENTY_FOUR_HOURS",
				}, false),
			},
		},
	}
}

func resourceNewRelicPathpointStageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"stage_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"health_rollup": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALERT_CONDITIONS",
					"AUTOMATIC_ROLL_UP",
				}, false),
			},
			"is_excluded": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"link": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"related": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointRelatedSchema(),
			},
			"stage_kpis": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNewRelicPathpointKpiSchema(),
			},
			"levels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNewRelicPathpointLevelSchema(),
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNewRelicPathpointRelatedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"target": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNewRelicPathpointLevelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"level_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"steps": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNewRelicPathpointStepSchema(),
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNewRelicPathpointStepSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"step_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_excluded": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"link": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scoped_accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointStepConfigSchema(),
			},
			"entity_search_query": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     resourceNewRelicPathpointSignalQuerySchema(),
			},
			"signals": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     resourceNewRelicPathpointSignalSchema(),
			},
			"health_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNewRelicPathpointStepConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"health_rollup": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"BEST_STATUS_WINS",
					"WORST_STATUS_WINS",
				}, false),
			},
			"threshold_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"FIXED",
					"PERCENTAGE",
				}, false),
			},
			"threshold_value": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceNewRelicPathpointSignalQuerySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_excluded": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNewRelicPathpointSignalSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"guid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ALERT",
					"ENTITY",
				}, false),
			},
			"is_excluded": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}
func resourceNewRelicPathpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	input, err := expandPathpointFlowInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	scope := pathpoint.PathPointScopeInput{
		ID:   accountID,
		Type: pathpoint.PathPointScopeTypeTypes.ACCOUNT,
	}

	result, err := client.PathPoint.PathPointCreateWithContext(ctx, *input, scope)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(result.GUID))
	return resourceNewRelicPathpointRead(ctx, d, meta)
}

func resourceNewRelicPathpointRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	result, err := client.PathPoint.GetFlowWithContext(ctx, accountID, pathpoint.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	return diag.FromErr(flattenPathpointFlowResult(result, d))
}

func resourceNewRelicPathpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	input, err := expandPathpointFlowUpdateInput(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.PathPoint.PathPointUpdateWithContext(ctx, pathpoint.EntityGUID(d.Id()), *input)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = result
	return resourceNewRelicPathpointRead(ctx, d, meta)
}

func resourceNewRelicPathpointDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	_, err := client.PathPoint.PathPointDeleteWithContext(ctx, pathpoint.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}