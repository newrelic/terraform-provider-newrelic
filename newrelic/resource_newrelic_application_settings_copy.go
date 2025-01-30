package newrelic

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/apm"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func resourceNewRelicApplicationSettingsCopy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicApplicationSettingsCreateCopy,
		ReadContext:   resourceNewRelicApplicationSettingsReadCopy,
		UpdateContext: resourceNewRelicApplicationSettingsUpdateCopy,
		DeleteContext: resourceNewRelicApplicationSettingsDeleteCopy,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: mergeSchemas(
			applicationSettingCommonSchema(),
			apmApplicationSettingsSchema(),
			mobileApplicationSettingsSchema(),
			browserApplicationSettingsSchema(),
		),
	}
}

func applicationSettingCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"guid": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func apmApplicationSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alias": {
			Type:     schema.TypeString,
			Optional: true,

			Description: "A name for this application in new relic",
		},
		"apm_config": {
			Type:     schema.TypeList,
			Optional: true,

			Description: "A specification of when the Monitor Downtime should end its repeat cycle, by number of occurrences or date.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"apdex_target": {
						Type:        schema.TypeFloat,
						Optional:    true,
						Description: "Response time threshold value for apdex",
					},
					"enable_server_side_config": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable server side configuration",
					},
				},
			},
		},
		"transaction_tracing": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "A specification of when the Monitor Downtime should end its repeat cycle, by number of occurrences or date.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "A Boolean to enable transaction tracing",
					},
					"transaction_threshold_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"APDEX_F", "VALUE"}, true),
						Description:  "Response time threshold value for apdex",
					},
					"transaction_threshold_value": {
						Type:     schema.TypeFloat,
						Optional: true,
						//RequiredWith: []string{"transaction_threshold_type"},
						Description: "Response time threshold value for apdex",
					},

					"log_sql": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "A Boolean to enable SQL tracing",
					},
					"record_sql": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"OBFUSCATED", "OFF", "RAW"}, true),
						Description:  "A Boolean to enable SQL tracing",
					},
					"stack_trace_threshold_value": {
						Type:        schema.TypeFloat,
						Optional:    true,
						Description: "Response time threshold value for stack trace of SQL",
					},
					"explain_enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Enable or disable the explain plan feature for slow SQL queries.",
					},

					"explain_threshold_type": {
						Type:         schema.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringInSlice([]string{"APDEX_F", "VALUE"}, true),
						Description:  "TBD",
					},
					"explain_threshold_value": {
						Type:     schema.TypeFloat,
						Optional: true,
						//RequiredWith: []string{"explain_threshold_type"},
						Description: "TBD",
					},
				},
			},
		},
		"error_collector": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "TBD",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,

						Description: "TBD",
					},
					"expected_error_classes": {
						Type:     schema.TypeSet,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,

						Description: "TBD",
					},
					"expected_error_codes": {
						Type:     schema.TypeSet,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,

						Description: "TBD",
					},
					"ignored_error_classes": {
						Type:     schema.TypeSet,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,

						Description: "TBD",
					},
					"ignored_error_codes": {
						Type:     schema.TypeSet,
						Elem:     &schema.Schema{Type: schema.TypeString},
						Optional: true,

						Description: "TBD",
					},
				},
			},
		},
		"tracer_type": {
			Type:     schema.TypeString,
			Optional: true,
			//Computed: true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return strings.EqualFold(old, new) // Case fold this attribute when diffing
			},
			ValidateFunc: validation.StringInSlice([]string{"CROSS_APPLICATION_TRACER", "DISTRIBUTED_TRACING", "NONE", "OPT_OUT"}, true),
			Description:  "TBD",
		},
		"thread_profiler_enabled": {
			Type:     schema.TypeBool,
			Optional: true,

			Description: "TBD",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				fmt.Println("new balue and old values", new, old)
				if new == "" {
					return true // Ignore changes when not set in config
				}
				return old == new
			},
		},
	}
}

func mobileApplicationSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"alias": {
			Type:     schema.TypeString,
			Optional: true,

			Description: "A name for this application in new relic",
		},
		"log_reporting": {
			Type:        schema.TypeList,
			MinItems:    1,
			Optional:    true,
			Description: "TBD",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "TBD",
					},
					"level": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "TBD",
					},
					"sampling_rate": {
						Type:         schema.TypeFloat,
						ValidateFunc: validation.FloatBetween(0, 100),
						Optional:     true,
						Description:  "TBD",
					},
				},
			},
		},
		"use_crash_reports": {
			Type:         schema.TypeBool,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{"CROSS_APPLICATION_TRACER", "DISTRIBUTED_TRACING", "NONE", "OPT_OUT"}, true),
			Description:  "TBD",
		},
		"application_exit_info": {
			Type:        schema.TypeList,
			MinItems:    1,
			Optional:    true,
			Description: "TBD",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"enabled": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "TBD",
					},
				},
			},
		},
	}
}

func browserApplicationSettingsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"session_replay": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"auto_start": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"error_sampling_rate": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"sampling_rate": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"mask_input_options": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"color": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"datetime_local": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"date": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"email": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"month": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"number": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"range": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"search": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"select": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"tel": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"text": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"text_area": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"time": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"url": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"week": {
									Type:     schema.TypeBool,
									Optional: true,
								},
							},
						},
					},
					"mask_all_inputs": {
						Type:     schema.TypeBool,
						Optional: true,
					},
					"block_selector": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"mask_text_selector": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},
		"session_trace": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"error_sampling_rate": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"mode": {
						Type:     schema.TypeString,
						Optional: true,
					},
					"sampling_rate": {
						Type:     schema.TypeInt,
						Optional: true,
					},
					"enabled": {
						Type:     schema.TypeBool,
						Optional: true,
					},
				},
			},
		},
		"browser_monitoring": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"distributed_tracing": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"cors_enabled": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"exclude_newrelic_header": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"cors_use_newrelic_header": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"cors_use_trace_context_headers": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"enabled": {
									Type:     schema.TypeBool,
									Optional: true,
								},
								"allowed_origins": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					"ajax": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"deny_list": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
		"browser_config": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"apdex_target": {
						Type:     schema.TypeInt,
						Optional: true,
					},
				},
			},
		},
	}
}

func resourceNewRelicApplicationSettingsCreateCopy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return resourceNewRelicApplicationSettingsUpdateCopy(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsReadCopy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Get("guid").(string)

	log.Printf("[INFO] Reading New Relic application %+v", guid)

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(guid))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("No New Relic application found with given guid %s", guid))
	}

	var dig diag.Diagnostics
	switch (*resp).(type) {
	case *entities.ApmApplicationEntity:
		entity := (*resp).(*entities.ApmApplicationEntity)
		d.SetId(string(entity.GUID))
		_ = d.Set("guid", string(entity.GUID))
		dig = diag.FromErr(setAPMApplicationValues(d, entity.ApmSettings))
	case *entities.MobileApplicationEntity:
		entity := (*resp).(*entities.MobileApplicationEntity)
		d.SetId(string(entity.GUID))
		_ = d.Set("guid", string(entity.GUID))
		dig = diag.FromErr(setMobileApplicationValues(d, entity.MobileSettings))
	case *entities.BrowserApplicationEntity:
		entity := (*resp).(*entities.BrowserApplicationEntity)
		d.SetId(string(entity.GUID))
		_ = d.Set("guid", string(entity.GUID))
		dig = diag.FromErr(setBrowserApplicationValues(d, entity.BrowserSettings))
	default:
		dig = diag.FromErr(fmt.Errorf("problem in retrieving application with GUID %s", guid))
	}
	return dig
}

func resourceNewRelicApplicationSettingsUpdateCopy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	updateApplicationParams := expandApplicationCopy(d)

	guid := d.Get("guid").(string)

	log.Printf("[INFO] Updating New Relic application %+v with params: %+v", guid, updateApplicationParams)

	agentApplicationSettingResult, err := client.APM.AgentApplicationSettingsUpdate(apm.EntityGUID(guid), *updateApplicationParams)

	if err != nil {
		return diag.FromErr(err)
	}
	if agentApplicationSettingResult == nil {
		return diag.FromErr(fmt.Errorf("something went wrong while Updating New Relic application"))
	}

	time.Sleep(2 * time.Second)

	d.SetId(string(agentApplicationSettingResult.GUID))

	return resourceNewRelicApplicationSettingsReadCopy(ctx, d, meta)
}

func resourceNewRelicApplicationSettingsDeleteCopy(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// You can not delete application settings
	return nil
}
