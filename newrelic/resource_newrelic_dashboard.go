package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	validIconValues = []string{
		"none",
		"archive",
		"bar-chart",
		"line-chart",
		"bullseye",
		"user",
		"usd",
		"money",
		"thumbs-up",
		"thumbs-down",
		"cloud",
		"bell",
		"bullhorn",
		"comments-o",
		"envelope",
		"globe",
		"shopping-cart",
		"sitemap",
		"clock-o",
		"crosshairs",
		"rocket",
		"users",
		"mobile",
		"tablet",
		"adjust",
		"dashboard",
		"flag",
		"flask",
		"road",
		"bolt",
		"cog",
		"leaf",
		"magic",
		"puzzle-piece",
		"bug",
		"fire",
		"legal",
		"trophy",
		"pie-chart",
		"sliders",
		"paper-plane",
		"life-ring",
		"heart",
	}

	validWidgetVisualizationValues = []string{
		"billboard",
		"gauge",
		"billboard_comparison",
		"facet_bar_chart",
		"faceted_line_chart",
		"facet_pie_chart",
		"facet_table",
		"faceted_area_chart",
		"heatmap",
		"attribute_sheet",
		"single_event",
		"histogram",
		"funnel",
		"raw_json",
		"event_feed",
		"event_table",
		"uniques_list",
		"line_chart",
		"comparison_line_chart",
		"markdown",
		"metric_line_chart",
		"application_breakdown",
	}
)

func resourceNewRelicDashboard() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Please use 'newrelic_one_dashboard' instead, for more information check out https://github.com/newrelic/terraform-provider-newrelic/issues/1297'",
		CreateContext:      resourceNewRelicDashboardCreate,
		ReadContext:        resourceNewRelicDashboardRead,
		UpdateContext:      resourceNewRelicDashboardUpdate,
		DeleteContext:      resourceNewRelicDashboardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the dashboard.",
			},
			"icon": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "bar-chart",
				ValidateFunc: validation.StringInSlice(validIconValues, false),
				Description:  "The icon for the dashboard.",
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: validation.StringInSlice([]string{"owner", "all"}, false),
				Description:  "Determines who can see the dashboard in an account. Valid values are all or owner. Defaults to all.",
			},
			"dashboard_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for viewing the dashboard.",
			},
			"editable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "editable_by_all",
				ValidateFunc: validation.StringInSlice([]string{"read_only", "editable_by_owner", "editable_by_all", "all"}, false),
				Description:  "Determines who can edit the dashboard in an account. Valid values are all, editable_by_all, editable_by_owner, or read_only. Defaults to editable_by_all.",
			},
			"grid_column_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				ValidateFunc: validation.IntInSlice([]int{3, 12}),
				Description:  "New Relic One supports a 3 column grid or a 12 column grid. New Relic Insights supports a 3 column grid.",
			},
			"filter": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A nested block that describes a dashboard filter. Exactly one nested filter block is allowed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_types": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
							Set:      schema.HashString,
						},
						"attributes": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Set:      schema.HashString,
						},
					},
				},
			},
			"widget": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    300,
				Description: "A nested block that describes a visualization. Up to 300 widget blocks are allowed in a dashboard definition.",
				Elem:        widgetSchemaElem(),
			},
		},
		SchemaVersion: 1,
	}
}

func widgetSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The target account ID to fetch data from, if not the current account.",
			},
			"widget_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ID of the widget.",
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A title for the widget.",
			},
			"visualization": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(validWidgetVisualizationValues, false),
				Description:  "How the widget visualizes data.",
			},
			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Width of the widget. Valid values are 1 to 3 inclusive. Defaults to 1.",
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "Height of the widget. Valid values are 1 to 3 inclusive. Defaults to 1.",
			},
			"row": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Row position of widget from top left, starting at 1.",
			},
			"column": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Column position of widget from top left, starting at 1.",
			},
			"notes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the widget.",
			},
			"nrql": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Valid NRQL query string.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The markdown source to be rendered in the widget.",
			},
			"threshold_red": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: float64AtLeast(0),
				Description:  "Threshold above which the displayed value will be styled with a red color.",
			},
			"threshold_yellow": {
				Type:         schema.TypeFloat,
				Optional:     true,
				ValidateFunc: float64AtLeast(0),
				Description:  "Threshold above which the displayed value will be styled with a yellow color.",
			},
			"drilldown_dashboard_id": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The ID of a dashboard to link to from the widget's facets.",
			},
			"duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"end_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"raw_metric_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"facet": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Set the order of result series.  Required when using `limit`.",
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The limit of distinct data series to display.  Requires `order_by` to be set.",
			},
			"entity_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "A collection of entity ids to display data for. These are typically application IDs.",
			},
			"metric": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A nested block that describes a metric.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The metric name to display.",
						},
						"units": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The metric units.",
						},
						"scope": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The metric scope.",
						},
						"values": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The metric values to display.",
						},
					},
				},
			},
			"compare_with": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A block describing a COMPARE WITH clause.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"offset_duration": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The offset duration for the COMPARE WITH clause.",
						},
						"presentation": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "The presentation settings for the rendered data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"color": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The color for the rendered data.",
									},
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name for the rendered data.",
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

func resourceNewRelicDashboardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("legacy dashboards have reached end of life, use `newrelic_one_dashboard` or `newrelic_one_dashboard_raw` instead")
}

func resourceNewRelicDashboardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("legacy dashboards have reached end of life, use `newrelic_one_dashboard` or `newrelic_one_dashboard_raw` instead")
}

func resourceNewRelicDashboardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("legacy dashboards have reached end of life, use `newrelic_one_dashboard` or `newrelic_one_dashboard_raw` instead")
}

func resourceNewRelicDashboardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return diag.Errorf("legacy dashboards have reached end of life, use `newrelic_one_dashboard` or `newrelic_one_dashboard_raw` instead")
}
