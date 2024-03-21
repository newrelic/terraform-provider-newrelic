package newrelic

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
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
			"account_id": {
				Type:        schema.TypeInt,
				Description: "ID of the newrelic account",
				Computed:    true,
				Optional:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The monitor type. Valid values are SIMPLE AND BROWSER.",
				ValidateFunc: validation.StringInSlice([]string{"SIMPLE", "BROWSER"}, false),
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
			"period_in_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The interval in minutes at which this monitor should run.",
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
				Description:  "The monitor status (ENABLED or DISABLED).",
				ValidateFunc: validateSyntheticMonitorStatus,
			},
			"validation_string": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The string to validate against in the response.",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
			"device_orientation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The device orientation the user would like to represent. Valid values are LANDSCAPE, PORTRAIT, or NONE.",
			},
			"device_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The device type that a user can select. Valid values are MOBILE, TABLET, or NONE.",
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
		simpleBrowserMonitorInput, browserMonitorErr := buildSyntheticsSimpleBrowserMonitor(d)
		if browserMonitorErr != nil {
			return diag.FromErr(browserMonitorErr)
		}
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

	setAttributesFromCreate(resp, d)

	return nil
}

func setAttributesFromCreate(res *synthetics.SyntheticsSimpleBrowserMonitorCreateMutationResult, d *schema.ResourceData) {
	_ = d.Set("validation_string", res.Monitor.AdvancedOptions.ResponseValidationText)
	_ = d.Set("verify_ssl", res.Monitor.AdvancedOptions.UseTlsValidation)
	_ = d.Set("name", res.Monitor.Name)
	_ = d.Set("status", string(res.Monitor.Status))
	_ = d.Set("period", string(res.Monitor.Period))
	_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[res.Monitor.Period])
	_ = d.Set("uri", res.Monitor.Uri)
	_ = d.Set("locations_public", res.Monitor.Locations.Public)
	_ = d.Set("locations_private", res.Monitor.Locations.Private)
	_ = d.Set("device_orientation", res.Monitor.AdvancedOptions.DeviceEmulation.DeviceOrientation)
	_ = d.Set("device_type", res.Monitor.AdvancedOptions.DeviceEmulation.DeviceType)

	if res.Monitor.Runtime.RuntimeType != "" {
		_ = d.Set("runtime_type", res.Monitor.Runtime.RuntimeType)
	}

	if res.Monitor.Runtime.RuntimeTypeVersion != "" {
		_ = d.Set("runtime_type_version", string(res.Monitor.Runtime.RuntimeTypeVersion))
	}

	if res.Monitor.Runtime.ScriptLanguage != "" {
		_ = d.Set("script_language", res.Monitor.Runtime.ScriptLanguage)
	}
}

func resourceNewRelicSyntheticsMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	// Detect old ID and convert to new format
	if len(d.Id()) == 36 {
		newGUID := fmt.Sprintf("%d|SYNTH|MONITOR|%s", accountID, d.Id())
		log.Printf("[INFO] Detected old ID %s converting to GUID %s", d.Id(), newGUID)
		d.SetId(base64.RawStdEncoding.EncodeToString([]byte(newGUID)))
	}

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(d.Id()))
	if err != nil {
		return diag.FromErr(err)
	}

	// This should probably be in go-client so we can use *errors.NotFound
	if *resp == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("account_id", accountID)
	setCommonSyntheticsMonitorAttributes(resp, d)

	return nil
}

// func to set output values in the read func.
func setCommonSyntheticsMonitorAttributes(v *entities.EntityInterface, d *schema.ResourceData) {
	switch e := (*v).(type) {
	case *entities.SyntheticMonitorEntity:
		err := setSyntheticsMonitorAttributes(d, map[string]string{
			"name":   e.Name,
			"type":   string(e.MonitorType),
			"uri":    e.MonitoredURL,
			"period": string(syntheticsMonitorPeriodValueMap[int(e.GetPeriod())]),
			"status": string(e.MonitorSummary.Status),
		})

		_ = d.Set("period_in_minutes", e.GetPeriod())

		if err != nil {
			diag.FromErr(err)
		}

		if e.MonitorType == entities.SyntheticMonitorTypeTypes.BROWSER {
			for _, t := range e.Tags {
				if k, ok := syntheticsMonitorTagKeyToSchemaAttrMap[t.Key]; ok {
					if len(t.Values) == 1 {
						_ = d.Set(k, t.Values[0])
					}
				}
			}
		}

		for _, t := range e.Tags {
			if t.Key == "responseValidationText" {
				if len(t.Values) == 1 {
					_ = d.Set("validation_string", t.Values[0])
				}
			}

			if t.Key == "useTlsValidation" {
				if len(t.Values) == 1 {
					v, _ := strconv.ParseBool(t.Values[0])
					_ = d.Set("verify_ssl", v)
				}
			}
		}
	}
}

func resourceNewRelicSyntheticsMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Updating New Relic Synthetics monitor %s", d.Id())

	var diags diag.Diagnostics

	monitorType, ok := d.GetOk("type")
	if !ok {
		log.Printf("monitor type must be specified")
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
		setSimpleMonitorAttributesFromUpdate(resp, d)

	case string(SyntheticsMonitorTypes.BROWSER):
		simpleBrowserMonitorUpdateInput, err := buildSyntheticsSimpleBrowserMonitorUpdateStruct(d)
		if err != nil {
			return diag.FromErr(err)
		}
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
		setSimpleBrowserAttributesFromUpdate(resp, d)
	}
	if len(diags) > 0 {
		return diags
	}
	return nil
}

func setSimpleMonitorAttributesFromUpdate(res *synthetics.SyntheticsSimpleMonitorUpdateMutationResult, d *schema.ResourceData) {
	_ = d.Set("name", res.Monitor.Name)
	_ = d.Set("period", string(res.Monitor.Period))
	_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[res.Monitor.Period])
	_ = d.Set("uri", res.Monitor.Uri)
	_ = d.Set("status", string(res.Monitor.Status))
	_ = d.Set("validation_string", res.Monitor.AdvancedOptions.ResponseValidationText)
	_ = d.Set("verify_ssl", res.Monitor.AdvancedOptions.UseTlsValidation)
	_ = d.Set("treat_redirect_as_failure", res.Monitor.AdvancedOptions.RedirectIsFailure)
	_ = d.Set("bypass_head_request", res.Monitor.AdvancedOptions.ShouldBypassHeadRequest)
	_ = d.Set("locations_public", res.Monitor.Locations.Public)
	_ = d.Set("locations_private", res.Monitor.Locations.Private)
}

func setSimpleBrowserAttributesFromUpdate(res *synthetics.SyntheticsSimpleBrowserMonitorUpdateMutationResult, d *schema.ResourceData) {
	_ = d.Set("name", res.Monitor.Name)
	_ = d.Set("period", string(res.Monitor.Period))
	_ = d.Set("period_in_minutes", syntheticsMonitorPeriodInMinutesValueMap[res.Monitor.Period])
	_ = d.Set("uri", res.Monitor.Uri)
	_ = d.Set("status", string(res.Monitor.Status))
	_ = d.Set("validation_string", res.Monitor.AdvancedOptions.ResponseValidationText)
	_ = d.Set("verify_ssl", res.Monitor.AdvancedOptions.UseTlsValidation)
	_ = d.Set("enable_screenshot_on_failure_and_script", res.Monitor.AdvancedOptions.EnableScreenshotOnFailureAndScript)
	_ = d.Set("locations_public", res.Monitor.Locations.Public)
	_ = d.Set("locations_private", res.Monitor.Locations.Private)
	_ = d.Set("device_orientation", res.Monitor.AdvancedOptions.DeviceEmulation.DeviceOrientation)
	_ = d.Set("device_type", res.Monitor.AdvancedOptions.DeviceEmulation.DeviceType)

	if res.Monitor.Runtime.RuntimeType != "" {
		_ = d.Set("runtime_type", res.Monitor.Runtime.RuntimeType)
	}

	if res.Monitor.Runtime.RuntimeTypeVersion != "" {
		_ = d.Set("runtime_type_version", string(res.Monitor.Runtime.RuntimeTypeVersion))
	}

	if res.Monitor.Runtime.ScriptLanguage != "" {
		_ = d.Set("script_language", res.Monitor.Runtime.ScriptLanguage)
	}
}

func resourceNewRelicSyntheticsMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	guid := synthetics.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	if _, err := client.Synthetics.SyntheticsDeleteMonitorWithContext(ctx, guid); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
