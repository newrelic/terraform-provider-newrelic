package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicOneDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicOneDashboardCreate,
		Read:   resourceNewRelicOneDashboardRead,
		Update: resourceNewRelicOneDashboardUpdate,
		Delete: resourceNewRelicOneDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:        schema.TypeList,
				Optional:    true,
				Description: "An area widget.",
				Elem:        dashboardWidgetGraphSchemaElem(),
			},
			"widget_bar": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A bar widget.",
				Elem:        dashboardWidgetGraphSchemaElem(),
			},
			"widget_billboard": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A billboard widget.",
				Elem:        dashboardWidgetBillboardSchemaElem(),
			},
			"widget_line": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A line widget.",
				Elem:        dashboardWidgetGraphSchemaElem(),
			},
			"widget_markdown": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A markdown widget.",
				Elem:        dashboardWidgetMarkdownSchemaElem(),
			},
			"widget_pie": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A pie widget.",
				Elem:        dashboardWidgetGraphSchemaElem(),
			},
			"widget_table": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A table widget.",
				Elem:        dashboardWidgetGraphSchemaElem(),
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
				Required:    true,
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

func dashboardWidgetBillboardSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the widget.",
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
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A title for the widget.",
			},
			"critical": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The critical threshold value.",
			},
			"warning": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The warning threshold value.",
			},
		},
	}
}

// dashboardWidgetGraphSchemaElem is a reusable schema element for graphs
func dashboardWidgetGraphSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the widget.",
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
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A title for the widget.",
			},
		},
	}
}
func dashboardWidgetMarkdownSchemaElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the widget.",
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
			"text": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A title for the widget.",
			},
		},
	}
}

func resourceNewRelicOneDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return fmt.Errorf("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	dashboard, err := expandDashboardInput(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating New Relic One dashboard: %s", dashboard.Name)

	created, err := client.Dashboards.DashboardCreate(accountID, *dashboard)
	if err != nil {
		return err
	}
	guid := created.EntityResult.GUID
	d.SetId(string(guid))

	return resourceNewRelicOneDashboardRead(d, meta)
}

// resourceNewRelicOneDashboardRead NerdGraph => Terraform reader
func resourceNewRelicOneDashboardRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return fmt.Errorf("err: NerdGraph support not present, but required for Read")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic One dashboard %s", d.Id())

	dashboard, err := client.Dashboards.GetDashboardEntity(entities.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenDashboardEntity(dashboard, d)
}

func resourceNewRelicOneDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return fmt.Errorf("err: NerdGraph support not present, but required for Update")
	}

	client := providerConfig.NewClient

	dashboard, err := expandDashboardInput(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Updating New Relic One dashboard '%s' (%s)", dashboard.Name, d.Id())

	result, err := client.Dashboards.DashboardUpdate(*dashboard, entities.EntityGUID(d.Id()))
	if err != nil {
		return err
	}

	// We have to use the Update Result, not a re-read of the entity as the changes take
	// some amount of time to be re-indexed
	return flattenDashboardUpdateResult(result, d)
}

func resourceNewRelicOneDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One dashboard %v", d.Id())

	if _, err := client.Dashboards.DashboardDelete(entities.EntityGUID(d.Id())); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return err
	}

	return nil
}
