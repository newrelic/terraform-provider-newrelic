package newrelic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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
		Create: resourceNewRelicDashboardCreate,
		Read:   resourceNewRelicDashboardRead,
		Update: resourceNewRelicDashboardUpdate,
		Delete: resourceNewRelicDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    300,
				Description: "A nested block that describes a visualization. Up to 300 widget blocks are allowed in a dashboard definition.",
				Elem:        widgetSchemaElem(),
			},
			"widgets": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    300,
				Description: "Array of widgets to display within the dashboard.",
				Elem:        widgetSchemaElem(),
			},
		},
	}
}

func resourceNewRelicDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	dashboard, err := expandDashboard(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating New Relic dashboard: %s", dashboard.Title)

	dashboard, err = client.Dashboards.CreateDashboard(*dashboard)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(dashboard.ID))

	return resourceNewRelicDashboardRead(d, meta)
	// return nil
}

func resourceNewRelicDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic dashboard %s", d.Id())

	dashboardID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard, err := client.Dashboards.GetDashboard(dashboardID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenDashboard(dashboard, d)
}

func resourceNewRelicDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	dashboard, err := expandDashboard(d)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard.ID = id
	log.Printf("[INFO] Updating New Relic dashboard %d", id)

	_, err = client.Dashboards.UpdateDashboard(*dashboard)
	if err != nil {
		return err
	}

	return resourceNewRelicDashboardRead(d, meta)
}

func resourceNewRelicDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic dashboard %v", id)

	if _, err := client.Dashboards.DeleteDashboard(id); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return err
	}

	return nil
}
