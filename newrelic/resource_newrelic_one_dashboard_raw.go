package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicOneDashboardRaw() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicOneDashboardRawCreate,
		ReadContext:   resourceNewRelicOneDashboardRawRead,
		UpdateContext: resourceNewRelicOneDashboardRawUpdate,
		DeleteContext: resourceNewRelicOneDashboardRawDelete,
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
				Elem:        dashboardRawPageSchemaElem(),
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

// dashboardPageSchemaElem returns the schema for a New Relic dashboard Page
func dashboardRawPageSchemaElem() *schema.Resource {
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
			"widget": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A generic widget configuration. Visualization id is required.",
				Elem:        dashboardRawWidgetSchemaElem(),
			},
		},
	}
}

func dashboardRawWidgetSchemaElem() *schema.Resource {
	s := dashboardRawWidgetSchemaBase()

	delete(s, "nrql_query") // No queries for Raw

	// Possibly call it VisualizationId
	s["visualization_id"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The visualization ID of the widget.",
	}

	// TODO: raw_configuration
	s["configuration"] = &schema.Schema{
		Type:             schema.TypeString,
		Required:         true,
		Description:      "The configuration of the widget.",
		DiffSuppressFunc: structure.SuppressJsonDiff,
	}

	// Expose linked_entity_guids
	s["linked_entity_guids"] = &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:         true,
		Description:      "(Optional) Related entity GUIDs. Currently only supports Dashboard entity GUIDs.",
		DiffSuppressFunc: structure.SuppressJsonDiff,
	}

	return &schema.Resource{
		Schema: s,
	}
}

func dashboardRawWidgetSchemaBase() map[string]*schema.Schema {
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
	}
}

func resourceNewRelicOneDashboardRawCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}
	dashboard, err := expandDashboardRawInput(d, defaultInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic One dashboard: %s", dashboard.Name)

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

	d.SetId(string(guid))

	return resourceNewRelicOneDashboardRawRead(ctx, d, meta)
}

// resourceNewRelicOneDashboardRawRead NerdGraph => Terraform reader
func resourceNewRelicOneDashboardRawRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Read")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic One dashboard %s", d.Id())

	dashboard, err := client.Dashboards.GetDashboardEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenDashboardRawEntity(dashboard, d))
}

func resourceNewRelicOneDashboardRawUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return diag.Errorf("err: NerdGraph support not present, but required for Update")
	}

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	defaultInfo := map[string]interface{}{
		"account_id": accountID,
	}
	dashboard, err := expandDashboardRawInput(d, defaultInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating New Relic One dashboard '%s' (%s)", dashboard.Name, d.Id())

	result, err := client.Dashboards.DashboardUpdateWithContext(ctx, *dashboard, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// We have to use the Update Result, not a re-read of the entity as the changes take
	// some amount of time to be re-indexed
	return diag.FromErr(flattenDashboardRawUpdateResult(result, d))
}

func resourceNewRelicOneDashboardRawDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One dashboard %v", d.Id())

	if _, err := client.Dashboards.DashboardDeleteWithContext(ctx, common.EntityGUID(d.Id())); err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return nil
		}
		return diag.FromErr(err)
	}

	return nil
}
