package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/dashboards"
	"github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

func resourceNewRelicOneDashboard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicOneDashboardCreate,
		ReadContext:   resourceNewRelicOneDashboardRead,
		UpdateContext: resourceNewRelicOneDashboardUpdate,
		DeleteContext: resourceNewRelicOneDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dashboard's name.",
			},
			"page": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "",
				Elem:        dashboardPageSchemaElem(),
			},
			// Optional
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to create the dashboard.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The dashboard's description.",
			},
			"permissions": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "public_read_only",
				ValidateFunc: validation.StringInSlice([]string{"private", "public_read_only", "public_read_write"}, false),
				Description:  "Determines who can see or edit the dashboard. Valid values are private, public_read_only, public_read_write. Defaults to public_read_only.",
			},
			// Computed
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the dashboard in New Relic.",
			},
			"permalink": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the dashboard.",
			},
			"variable": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Dashboard-local variable definitions.",
				Elem:        dashboardVariableSchemaElem(),
			},
		},
		CustomizeDiff: validateDashboardArguments,
	}
}

func dashboardVariableSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"default_values": {
				Type:        schema.TypeList,
				Optional:    true,
				MinItems:    1,
				Description: "Default values for this variable.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_multi_selection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether this variable supports multiple selection or not. Only applies to variables of type NRQL or ENUM.",
			},
			"item": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of possible values for variables of type ENUM",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"title": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A human-friendly display string for this value.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A possible variable value",
						},
					},
				},
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The variable identifier.",
			},
			"options": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Options applied to the variable.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ignore_time_range": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Only applies to variables of type NRQL. With this turned on, the time range for the NRQL query will override the time picker on dashboards and other pages. Turn this off to use the time picker as normal.",
						},
						"excluded": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Only applies to variables of type NRQL. With this turned on, query condition defined with the variable will not be included in the query.",
						},
					},
				},
			},
			"nrql_query": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Configuration for variables of type NRQL.",
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "New Relic account ID(s) to issue the query against.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "NRQL formatted query.",
						},
					},
				},
			},
			"replacement_strategy": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Indicates the strategy to apply when replacing a variable in a NRQL query.",
				ValidateFunc: validation.StringInSlice([]string{"default", "identifier", "number", "string"}, false),
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Human-friendly display string for this variable.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Specifies the data type of the variable and where its possible values may come from.",
				ValidateFunc: validation.StringInSlice([]string{"enum", "nrql", "string"}, false),
			},
		},
	}
}

// dashboardPageElem returns the schema for a New Relic dashboard Page
func dashboardPageSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The dashboard page's description.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dashboard page's name.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the dashboard page in New Relic.",
			},

			// All the widget types below
			"widget_area": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An area widget.",
				Elem:        dashboardWidgetAreaSchemaElem(),
			},
			"widget_bar": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A bar widget.",
				Elem:        dashboardWidgetBarSchemaElem(),
			},
			"widget_billboard": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A billboard widget.",
				Elem:        dashboardWidgetBillboardSchemaElem(),
			},
			"widget_bullet": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A bullet widget.",
				Elem:        dashboardWidgetBulletSchemaElem(),
			},
			"widget_funnel": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A funnel widget.",
				Elem:        dashboardWidgetFunnelSchemaElem(),
			},
			"widget_heatmap": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A heatmap widget.",
				Elem:        dashboardWidgetHeatmapSchemaElem(),
			},
			"widget_histogram": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A histogram widget.",
				Elem:        dashboardWidgetHistogramSchemaElem(),
			},
			"widget_line": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A line widget.",
				Elem:        dashboardWidgetLineSchemaElem(),
			},
			"widget_markdown": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A markdown widget.",
				Elem:        dashboardWidgetMarkdownSchemaElem(),
			},
			"widget_pie": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A pie widget.",
				Elem:        dashboardWidgetPieSchemaElem(),
			},
			"widget_log_table": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A log table widget.",
				Elem:        dashboardWidgetLogTableSchemaElem(),
			},
			"widget_table": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A table widget.",
				Elem:        dashboardWidgetTableSchemaElem(),
			},
			"widget_json": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A JSON widget.",
				Elem:        dashboardWidgetJSONSchemaElem(),
			},
			"widget_stacked_bar": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A stacked bar widget.",
				Elem:        dashboardWidgetStackedBarSchemaElem(),
			},
		},
	}
}

func dashboardWidgetSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The ID of the widget.",
		},
		"title": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "A title for the widget.",
		},
		"column": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"height": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntAtLeast(1),
		},
		"row": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"width": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      4,
			ValidateFunc: validation.IntBetween(1, 12),
		},
		"nrql_query": {
			Type:     schema.TypeList,
			Required: true,
			Elem:     dashboardWidgetNRQLQuerySchemaElem(),
		},
		"refresh_rate": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"initial_sorting": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			MinItems: 1,
			Elem:     dashboardWidgetInitialSortingSchemaElem(),
		},
		"data_format": {
			Type:     schema.TypeList,
			Optional: true,
			MinItems: 1,
			Elem:     dashboardWidgetDataFormatSchemaElem(),
		},
		"ignore_time_range": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"facet_show_other_series": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"legend_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"y_axis_left_min": {
			Type:     schema.TypeFloat,
			Optional: true,
			Default:  0,
		},
		"y_axis_left_max": {
			Type:     schema.TypeFloat,
			Optional: true,
			Default:  0,
		},
		"null_values": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     dashboardWidgetNullValuesSchemaElem(),
		},
		"units": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     dashboardWidgetUnitsSchemaElem(),
		},
		"colors": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     dashboardWidgetColorSchemaElem(),
		},
	}
}

func dashboardWidgetColorSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"color": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"series_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"color": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Color code",
						},
						"series_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Series name",
						},
					},
				},
			},
		},
	}
}

func dashboardWidgetUnitsSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"unit": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"series_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Unit name",
						},
						"series_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Series name",
						},
					},
				},
			},
		},
	}
}

func dashboardWidgetNullValuesSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"null_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"series_overrides": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"null_value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Null value",
						},
						"series_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Series name",
						},
					},
				},
			},
		},
	}
}

func dashboardWidgetInitialSortingSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines the sort order. Either ascending or descending.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The column name to be sorted",
			},
		},
	}
}

func dashboardWidgetDataFormatSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The column name to be sorted",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Defines the type of the mentioned column",
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Defines the format of the mentioned type",
			},
			"precision": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The precision of the type",
			},
		},
	}
}

// dashboardWidgetNRQLQuerySchemaElem defines a NRQL query for use on a dashboard
//
// see: newrelic/newrelic-client-go/pkg/entities/DashboardWidgetQuery
func dashboardWidgetNRQLQuerySchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The account id used for the NRQL query.",
			},
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NRQL query.",
			},
		},
	}
}

func dashboardWidgetAreaSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetBarSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetBillboardSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["critical"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The critical threshold value.",
	}

	s["warning"] = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The warning threshold value.",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetBulletSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["limit"] = &schema.Schema{
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "The maximum value for the visualization",
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetFunnelSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetHeatmapSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetHistogramSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetLineSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["y_axis_left_zero"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specifies if the values on the graph to be rendered need to be fit to scale, or printed within the specified range.",
	}

	s["is_label_visible"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Specified if the label should be visible in the graph created when specified with thresholds.",
	}

	// adding this attribute in the schema of line widgets, since 'y_axis_right' is currently available only to line widgets
	s["y_axis_right"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"y_axis_right_zero": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "An attribute that helps specify the Y-Axis on the right of the line widget.",
				},
				"y_axis_right_series": {
					Type:        schema.TypeSet,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "A set of series that helps specify the Y-Axis on the right of the line widget.",
					Optional:    true,
				},
				"y_axis_right_min": {
					Type:        schema.TypeFloat,
					Description: "Minimum value of the range to be specified with the Y-Axis on the right of the line widget.",
					Optional:    true,
				},
				"y_axis_right_max": {
					Type:        schema.TypeFloat,
					Description: "Minimum value of the range to be specified with the Y-Axis on the right of the line widget.",
					Optional:    true,
				},
			},
		},
	}

	s["threshold"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The number from which the range starts in thresholds.",
				},
				"to": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The number at which the range ends in thresholds.",
				},
				"name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the threshold created.",
				},
				"severity": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.SUCCESS),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.WARNING),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.UNAVAILABLE),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.SEVERE),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.CRITICAL),
					}, false),
					Description: "Severity of the threshold, which would reflect in the widget, in the range of the threshold specified.",
				},
			},
		},
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetJSONSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetMarkdownSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	delete(s, "nrql_query") // No queries for Markdown

	s["text"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetStackedBarSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetPieSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetLogTableSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetTableSchemaElem() *schema.Resource {
	s := dashboardWidgetSchemaBase()

	s["linked_entity_guids"] = dashboardWidgetLinkedEntityGUIDsSchema()
	s["filter_current_dashboard"] = dashboardWidgetFilterCurrentDashboardSchema()

	s["threshold"] = &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"from": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The number from which the range starts in thresholds.",
				},
				"to": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The number at which the range ends in thresholds.",
				},
				"column_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Name of the column in the table, to which the threshold would be applied.",
				},
				"severity": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.SUCCESS),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.WARNING),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.UNAVAILABLE),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.SEVERE),
						string(dashboards.DashboardLineTableWidgetsAlertSeverityTypes.CRITICAL),
					}, false),
					Description: "Severity of the threshold, which would reflect in the widget, in the range of the threshold specified.",
				},
			},
		},
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardWidgetLinkedEntityGUIDsSchema() *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "Related entities. Currently only supports Dashboard entities, but may allow other cases in the future.",
	}
}

func dashboardWidgetFilterCurrentDashboardSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Use this item to filter the current dashboard",
	}
}

func resourceNewRelicOneDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}
	dashboard, err := expandDashboardInput(d, defaultInfo, "")
	if err != nil {
		return diag.FromErr(err)
	}

	filterWidgets, err := findDashboardWidgetFilterCurrentDashboard(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Creating New Relic One dashboard: %s", dashboard.Name))

	created, err := client.Dashboards.DashboardCreateWithContext(ctx, accountID, *dashboard)
	if err != nil {
		return diag.FromErr(err)
	}
	guid := created.EntityResult.GUID
	if guid == "" {
		var errMessages string
		for _, e := range created.Errors {
			errMessages += "[" + string(e.Type) + ": " + e.Description + "]"
		}

		return diag.Errorf("err: newrelic_one_dashboard Create failed: %s", errMessages)
	}

	tflog.Info(ctx, fmt.Sprintf("New Dashboard GUID: %s", guid))

	d.SetId(string(guid))

	res := resourceNewRelicOneDashboardRead(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Number of widgets with filter_current_dashboard: %d", len(filterWidgets)))
	if len(filterWidgets) > 0 {

		err = setDashboardWidgetFilterCurrentDashboardLinkedEntity(d, filterWidgets)
		if err != nil {
			return diag.FromErr(err)
		}

		dashboard, err := expandDashboardInput(d, defaultInfo, created.EntityResult.Name)
		if err != nil {
			return diag.FromErr(err)
		}

		result, err := client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, guid)
		if err != nil {
			return diag.FromErr(err)
		}

		return diag.FromErr(flattenDashboardUpdateResult(result, d))

	}

	return res
}

// resourceNewRelicOneDashboardRead NerdGraph => Terraform reader
func resourceNewRelicOneDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	tflog.Info(ctx, fmt.Sprintf("Reading New Relic One dashboard %s", d.Id()))

	dashboard, err := client.Dashboards.GetDashboardEntityWithContext(ctx, common.EntityGUID(d.Id()))

	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenDashboardEntity(dashboard, d))
}

func resourceNewRelicOneDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}

	filterWidgets, err := findDashboardWidgetFilterCurrentDashboard(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the linked guid entities, if the page guid value is empty it is set to nil
	err = setDashboardWidgetFilterCurrentDashboardLinkedEntity(d, filterWidgets)
	if err != nil {
		return diag.FromErr(err)
	}

	dashboard, err := expandDashboardInput(d, defaultInfo, "")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Updating New Relic One dashboard '%s' (%s)", dashboard.Name, d.Id()))

	updated, err := client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	guid := updated.EntityResult.GUID
	if guid == "" {
		var errMessages string
		for _, e := range updated.Errors {
			errMessages += "[" + string(e.Type) + ": " + e.Description + "]"
		}

		return diag.Errorf("err: newrelic_one_dashboard Update failed: %s", errMessages)
	}

	diagErr := resourceNewRelicOneDashboardRead(ctx, d, meta)
	if diagErr != nil {
		return diagErr
	}

	tflog.Info(ctx, fmt.Sprintf("Number of widgets with filter_current_dashboard: %d", len(filterWidgets)))
	// If there are widgets with filter_current_dashboard, we need to update the linked guid entities
	if len(filterWidgets) > 0 {
		err = setDashboardWidgetFilterCurrentDashboardLinkedEntity(d, filterWidgets)
		if err != nil {
			return diag.FromErr(err)
		}

		dashboard, err = expandDashboardInput(d, defaultInfo, updated.EntityResult.Name)
		if err != nil {
			return diag.FromErr(err)
		}

		updated, err = client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, guid)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// We have to use the Update Result, not a re-read of the entity as the changes take
	// some amount of time to be re-indexed
	return diag.FromErr(flattenDashboardUpdateResult(updated, d))
}

func resourceNewRelicOneDashboardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	tflog.Info(ctx, fmt.Sprintf("Deleting New Relic One dashboard %v", d.Id()))

	if _, err := client.Dashboards.DashboardDeleteWithContext(ctx, common.EntityGUID(d.Id())); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
