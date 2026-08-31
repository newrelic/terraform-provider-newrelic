package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kpi_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"account_id": {
							Type:     schema.TypeInt,
							Optional: true,
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
							Elem: &schema.Resource{
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
										Elem: &schema.Resource{
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
										},
									},
									"time_window": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"custom_range": {
													Type:     schema.TypeString,
													Optional: true,
												},
												"relative_range": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
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
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"stages": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
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
							Elem: &schema.Resource{
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
							},
						},
						"stage_kpis": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kpi_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"account_id": {
										Type:     schema.TypeInt,
										Optional: true,
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
							},
						},
						"levels": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"level_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"steps": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
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
												"entity_search_query": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
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
													},
												},
												"config": {
													Type:     schema.TypeList,
													Optional: true,
													MaxItems: 1,
													Elem: &schema.Resource{
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
													},
												},
												"signals": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
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
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"guid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
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
				Elem: &schema.Resource{
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
				},
			},
			"time_window": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"custom_range": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"relative_range": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
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
							},
						},
					},
				},
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
	_ = d.Set("account_id", accountID)

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

	if result == nil || result.GUID == "" {
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

var _ = fmt.Sprintf
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
	_ = d.Set("account_id", accountID)

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

	if result == nil || result.GUID == "" {
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