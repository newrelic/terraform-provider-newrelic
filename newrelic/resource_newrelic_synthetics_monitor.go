package newrelic

import (
	"context"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The monitor type. Valid values are SIMPLE AND BROWSER.",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsScriptMonitorTypes(), false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of this monitor.",
			},
			"period": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorPeriods(), false),
			},
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URI for the monitor to hit.",
			},
			"locations_public": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"locations_public", "locations_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"locations_private": {
				Type:         schema.TypeSet,
				Elem:         &schema.Schema{Type: schema.TypeString},
				MinItems:     1,
				Optional:     true,
				AtLeastOneOf: []string{"locations_public", "locations_private"},
				Description:  "The locations in which this monitor should be run.",
			},
			"status": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
				ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorStatuses(), false),
			},
			"validation_string": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The string to validate against in the response.",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Verify SSL.",
			},
			"bypass_head_request": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Bypass HEAD request.",
			},
			"treat_redirect_as_failure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Fail the monitor check if redirected.",
			},
			"runtime_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The runtime type that the monitor will run",
			},
			"runtime_type_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specific version of the runtime type selected",
			},
			"script_language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The programing language that should execute the script",
			},
			"tag": {
				Type:        schema.TypeSet,
				Optional:    true,
				MinItems:    1,
				Description: "The tags that will be associated with the monitor",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the tag key",
						},
						"values": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Values associated with the tag key",
						},
					},
				},
			},
			"enable_screenshot_on_failure_and_script": {
				Type:        schema.TypeBool,
				Description: "Capture a screenshot during job execution",
				Optional:    true,
			},
			"custom_header": {
				Type:        schema.TypeSet,
				Description: "Custom headers to use in monitor job",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "Header name",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "Header value",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func resourceNewRelicSyntheticsMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	var diags diag.Diagnostics
	var resp *synthetics.SyntheticsSimpleBrowserMonitorCreateMutationResult
	var err error

	monitorType, ok := d.GetOk("type")
	if !ok {
		log.Printf("Not Monitor type specified")
	}

	switch monitorType.(string) {
	case string(SyntheticsMonitorTypes.SIMPLE):

		simpleMonitorInput := buildSyntheticsSimpleMonitor(d)

		resp, err = client.Synthetics.SyntheticsCreateSimpleMonitorWithContext(ctx, accountID, simpleMonitorInput)

	case string(SyntheticsMonitorTypes.BROWSER):

		simpleBrowserMonitorInput := buildSyntheticsSimpleBrowserMonitor(d)

		resp, err = client.Synthetics.SyntheticsCreateSimpleBrowserMonitorWithContext(ctx, accountID, simpleBrowserMonitorInput)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if len(resp.Errors) > 0 {
		for _, err := range resp.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  string(err.Type) + " " + err.Description,
			})
		}
	}

	d.SetId(string(resp.Monitor.GUID))

	if len(diags) > 0 {
		return diags
	}

	return nil
}

func resourceNewRelicSyntheticsMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))

	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setCommonSyntheticsMonitorAttributes(resp, d)

	return nil
}

//func to set output values in the read func.
func setCommonSyntheticsMonitorAttributes(v *entities.EntityInterface, d *schema.ResourceData) {

	switch e := (*v).(type) {

	case *entities.SyntheticMonitorEntity:
		_ = d.Set("name", e.Name)
		_ = d.Set("type", e.MonitorType)
		_ = d.Set("uri", e.MonitoredURL)
	}
}

func resourceNewRelicSyntheticsMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Updating New Relic Synthetics monitor %s", d.Id())

	var diags diag.Diagnostics

	monitorType, ok := d.GetOk("type")
	if !ok {
		log.Printf("Not Monitor type specified")
	}

	guid := synthetics.EntityGUID(d.Id())

	switch monitorType.(string) {
	case string(SyntheticsMonitorTypes.SIMPLE):
		simpleMonitorUpdateInput := buildSyntheticsSimpleMonitorUpdateStruct(d)

		resp, err := client.Synthetics.SyntheticsUpdateSimpleMonitorWithContext(ctx, guid, simpleMonitorUpdateInput)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(resp.Errors) > 0 {
			for _, err := range resp.Errors {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  string(err.Type) + " " + err.Description,
				})
			}
		}
		d.SetId(string(resp.Monitor.GUID))

	case string(SyntheticsMonitorTypes.BROWSER):
		simpleBrowserMonitorUpdateInput := buildSyntheticsSimpleBrowserMonitorUpdateStruct(d)

		resp, err := client.Synthetics.SyntheticsUpdateSimpleBrowserMonitorWithContext(ctx, guid, simpleBrowserMonitorUpdateInput)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(resp.Errors) > 0 {
			for _, err := range resp.Errors {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  string(err.Type) + " " + err.Description,
				})
			}
		}
		d.SetId(string(resp.Monitor.GUID))
	}
	if len(diags) > 0 {
		return diags
	}
	return nil
}

func resourceNewRelicSyntheticsMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(*ProviderConfig).NewClient

	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	_, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}
