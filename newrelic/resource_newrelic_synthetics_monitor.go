package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/newrelic/newrelic-client-go/pkg/common"

	"github.com/newrelic/newrelic-client-go/pkg/entities"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The monitor type. Valid values are SIMPLE, BROWSER, SCRIPT_BROWSER, and SCRIPT_API.",
				ValidateFunc: validation.StringInSlice([]string{
					"SIMPLE",
					"BROWSER",
					"SCRIPT_API",
					"SCRIPT_BROWSER",
				}, false),
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of this monitor.",
			},
			"frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: intInSlice([]int{1, 5, 10, 15, 30, 60, 360, 720, 1440}),
				Description:  "The interval (in minutes) at which this monitor should run. Valid values are 1, 5, 10, 15, 30, 60, 360, 720, or 1440.",
			},
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URI for the monitor to hit.",
				// TODO: ValidateFunc (required if SIMPLE or BROWSER)
			},
			"locations": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				MinItems:    1,
				Required:    true,
				Description: "The locations in which this monitor should be run.",
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"MUTED",
					"DISABLED",
				}, false),
			},
			"sla_threshold": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Default:     7,
				Description: "The base threshold (in seconds) to calculate the apdex score for use in the SLA report. (Default 7 seconds)",
			},
			// TODO: ValidationFunc (options only valid if SIMPLE or BROWSER)
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
				Required:    true,
				Description: "",
			},
			"runtime_type_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"script_language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "",
						},
					},
				},
			},
			"enable_screenshot_on_failure_and_script": {
				Type:        schema.TypeBool,
				Description: "",
				Optional:    true,
			},
			"custom_headers": {
				Type:        schema.TypeSet,
				Description: "",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "",
							Optional:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func buildSyntheticsMonitorBase(d *schema.ResourceData) map[string]interface{} {

	var monitorInputs = make(map[string]interface{})

	if screenShot, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		monitorInputs["enable_screenshot_on_failure_and_script"] = screenShot
	}
	if customHeaders, ok := d.GetOk("custom_headers"); ok {
		monitorInputs["custom_headers"] = customHeaders
	}
	if redirectFailure, ok := d.GetOk("treat_redirect_as_failure"); ok {
		monitorInputs["treat_redirect_as_failure"] = redirectFailure
	}
	if respValidateText, ok := d.GetOk("validation_string"); ok {
		monitorInputs["validation_string"] = respValidateText
	}
	if pass, ok := d.GetOk("bypass_head_request"); ok {
		monitorInputs["bypass_head_request"] = pass
	}
	if tlsValidations, ok := d.GetOk("verify_ssl"); ok {
		monitorInputs["verify_ssl"] = tlsValidations
	}
	if locations, ok := d.GetOk("locations"); ok {
		monitorInputs["locations"] = locations
	}
	if name, ok := d.GetOk("name"); ok {
		monitorInputs["name"] = name
	}
	if period, ok := d.GetOk("frequency"); ok {
		monitorInputs["frequency"] = period
	}
	if runtimeType, ok := d.GetOk("runtime_type"); ok {
		monitorInputs["runtime_type"] = runtimeType
	}
	if runtimeTypeVersion, ok := d.GetOk("runtime_type_version"); ok {
		monitorInputs["runtime_type_version"] = runtimeTypeVersion
	}
	if scriptLanguage, ok := d.GetOk("script_language"); ok {
		monitorInputs["script_language"] = scriptLanguage
	}
	if status, ok := d.GetOk("synthetics_monitor_status"); ok {
		monitorInputs["synthetics_monitor_status"] = status
	}
	if tags, ok := d.GetOk("tags"); ok {
		monitorInputs["tags"] = tags
	}
	if uri, ok := d.GetOk("uri"); ok {
		monitorInputs["uri"] = uri
	}
	if script, ok := d.GetOk("script"); ok {
		monitorInputs["script"] = script
	}
	if status, ok := d.GetOk("status"); ok {
		monitorInputs["status"] = status
	}
	return monitorInputs
}

//1, 5, 10, 15, 30, 60, 360, 720, 1440

func periodConvIntToString(v interface{}) synthetics.SyntheticsMonitorPeriod {
	var output synthetics.SyntheticsMonitorPeriod
	switch v.(int) {
	case 1:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_MINUTE
	case 5:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_5_MINUTES
	case 10:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_10_MINUTES
	case 15:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_15_MINUTES
	case 30:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_30_MINUTES
	case 60:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_HOUR
	case 360:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_6_HOURS
	case 720:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_12_HOURS
	case 1440:
		output = synthetics.SyntheticsMonitorPeriodTypes.EVERY_DAY
	}
	return output
}

func buildSyntheticsScriptAPIMonitorStruct(d *schema.ResourceData) synthetics.SyntheticsCreateScriptAPIMonitorInput {
	scriptAPIMonitorInput := synthetics.SyntheticsCreateScriptAPIMonitorInput{}
	input := buildSyntheticsMonitorBase(d)
	scriptAPIMonitorInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])
	scriptAPIMonitorInput.Name = input["name"].(string)
	scriptAPIMonitorInput.Period = periodConvIntToString(input["frequency"])
	scriptAPIMonitorInput.Runtime.RuntimeType = input["runtime_type"].(string)
	scriptAPIMonitorInput.Runtime.ScriptLanguage = input["script_language"].(string)
	scriptAPIMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	scriptAPIMonitorInput.Script = input["script"].(string)
	scriptAPIMonitorInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	scriptAPIMonitorInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())

	return scriptAPIMonitorInput
}

func expandSyntheticsScriptMonitorLocations(v interface{}) synthetics.SyntheticsScriptedMonitorLocationsInput {
	locationsRaw := v.(*schema.Set)
	locations := make([]string, locationsRaw.Len())
	for i, v := range locationsRaw.List() {
		locations[i] = fmt.Sprint(v)
	}
	inputLocations := synthetics.SyntheticsScriptedMonitorLocationsInput{
		Public: locations,
	}
	return inputLocations
}

func expandSyntheticsTags(tags []interface{}) []synthetics.SyntheticsTag {
	out := make([]synthetics.SyntheticsTag, len(tags))
	for i, v := range tags {
		tag := v.(map[string]interface{})
		expanded := synthetics.SyntheticsTag{
			Key:    tag["key"].(string),
			Values: expandSyntheticsTagValues(tag["values"].(*schema.Set).List()),
		}
		out[i] = expanded
	}
	return out
}

func expandSyntheticsTagValues(v []interface{}) []string {
	values := make([]string, len(v))
	for i, value := range v {
		values[i] = value.(string)
	}
	return values
}

func buildSyntheticsScriptBrowserMonitorStruct(d *schema.ResourceData) synthetics.SyntheticsCreateScriptBrowserMonitorInput {
	scriptBrowserMonitorInput := synthetics.SyntheticsCreateScriptBrowserMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	scriptBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = input["enable_screenshot_on_failure_and_script"].(bool)
	scriptBrowserMonitorInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])
	scriptBrowserMonitorInput.Name = input["name"].(string)
	scriptBrowserMonitorInput.Period = periodConvIntToString(input["frequency"])
	scriptBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	scriptBrowserMonitorInput.Runtime.ScriptLanguage = input["script_language"].(string)
	scriptBrowserMonitorInput.Runtime.RuntimeType = input["runtime_type"].(string)
	scriptBrowserMonitorInput.Script = input["script"].(string)
	scriptBrowserMonitorInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	scriptBrowserMonitorInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())

	return scriptBrowserMonitorInput
}

func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleBrowserMonitorInput {
	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	simpleBrowserMonitorInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"].(*schema.Set).List())
	simpleBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = input["enable_screenshot_on_failure_and_script"].(bool)
	simpleBrowserMonitorInput.AdvancedOptions.ResponseValidationText = input["validation_string"].(string)
	simpleBrowserMonitorInput.AdvancedOptions.UseTlsValidation = input["verify_ssl"].(bool)
	simpleBrowserMonitorInput.Locations = expandSyntheticsLocations(input["locations"])
	simpleBrowserMonitorInput.Name = input["name"].(string)
	simpleBrowserMonitorInput.Period = periodConvIntToString(input["frequency"])
	simpleBrowserMonitorInput.Runtime.RuntimeType = input["runtime_type"].(string)
	simpleBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	simpleBrowserMonitorInput.Runtime.ScriptLanguage = input["scriptLanguage"].(string)
	simpleBrowserMonitorInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	simpleBrowserMonitorInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())
	simpleBrowserMonitorInput.Uri = input["uri"].(string)

	return simpleBrowserMonitorInput
}

func expandSyntheticsLocations(v interface{}) synthetics.SyntheticsLocationsInput {
	locationsRaw := v.(*schema.Set)
	locations := make([]string, locationsRaw.Len())
	for i, v := range locationsRaw.List() {
		locations[i] = fmt.Sprint(v)
	}
	inputLocations := synthetics.SyntheticsLocationsInput{
		Public: locations,
	}
	return inputLocations
}

func expandCustomHeaders(headers []interface{}) []synthetics.SyntheticsCustomHeaderInput {
	output := make([]synthetics.SyntheticsCustomHeaderInput, len(headers))
	for i, v := range headers {
		header := v.(map[string]interface{})
		expanded := synthetics.SyntheticsCustomHeaderInput{
			Name:  header["name"].(string),
			Value: header["value"].(string),
		}
		output[i] = expanded
	}
	return output
}

func buildSyntheticsSimpleMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleMonitorInput {
	simpleMonitorInput := synthetics.SyntheticsCreateSimpleMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	simpleMonitorInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"].(*schema.Set).List())
	simpleMonitorInput.AdvancedOptions.RedirectIsFailure = input["treat_redirect_as_failure"].(bool)
	simpleMonitorInput.AdvancedOptions.ResponseValidationText = input["validation_string"].(string)
	simpleMonitorInput.AdvancedOptions.ShouldBypassHeadRequest = input["bypass_head_request"].(bool)
	simpleMonitorInput.AdvancedOptions.UseTlsValidation = input["verify_ssl"].(bool)
	simpleMonitorInput.Locations = expandSyntheticsLocations(input["locations"])
	simpleMonitorInput.Name = input["name"].(string)
	simpleMonitorInput.Period = periodConvIntToString(input["frequency"])
	simpleMonitorInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	simpleMonitorInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())
	simpleMonitorInput.Uri = input["uri"].(string)

	return simpleMonitorInput
}

// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func resourceNewRelicSyntheticsMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	var diags diag.Diagnostics

	monitorType, ok := d.GetOk("type")
	if !ok {
		log.Printf("Not Monitor type specified")
	}
	switch monitorType.(string) {

	case "SIMPLE":
		simpleMonitorInput := buildSyntheticsSimpleMonitor(d)
		resp, err := client.Synthetics.SyntheticsCreateSimpleMonitorWithContext(ctx, accountID, simpleMonitorInput)
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

	case "BROWSER":
		simpleBrowserMonitorInput := buildSyntheticsSimpleBrowserMonitor(d)
		resp, err := client.Synthetics.SyntheticsCreateSimpleBrowserMonitorWithContext(ctx, accountID, simpleBrowserMonitorInput)
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

	case "SCRIPT_API":
		scriptedAPIMonitorInput := buildSyntheticsScriptAPIMonitorStruct(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptAPIMonitorWithContext(ctx, accountID, scriptedAPIMonitorInput)
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

	case "SCRIPT_BROWSER":
		scriptedBrowserMonitorInput := buildSyntheticsScriptBrowserMonitorStruct(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptBrowserMonitorWithContext(ctx, accountID, scriptedBrowserMonitorInput)
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

	return nil
}

func resourceNewRelicSyntheticsMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	guid := common.EntityGUID(d.Id())
	resp, err := client.Entities.GetEntity(guid)
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

func setCommonSyntheticsMonitorAttributes(v *entities.EntityInterface, d *schema.ResourceData) {
	switch e := (*v).(type) {
	case *entities.SyntheticMonitorEntity:
		_ = d.Set("guid", e.GUID)
		_ = d.Set("name", e.Name)
		_ = d.Set("type", e.MonitorType)
		_ = d.Set("period", e.Period)
		_ = d.Set("tags", e.Tags)
		_ = d.Set("uri", e.MonitoredURL)
	}
}

func buildSyntheticsScriptAPIMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptAPIMonitorInput {
	scriptAPIMonitorUpdateInput := synthetics.SyntheticsUpdateScriptAPIMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	scriptAPIMonitorUpdateInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])
	scriptAPIMonitorUpdateInput.Name = input["name"].(string)
	scriptAPIMonitorUpdateInput.Period = synthetics.SyntheticsMonitorPeriod(input["frequency"].(string))
	scriptAPIMonitorUpdateInput.Runtime.RuntimeType = input["runtime_type"].(string)
	scriptAPIMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	scriptAPIMonitorUpdateInput.Runtime.ScriptLanguage = input["script_language"].(string)
	scriptAPIMonitorUpdateInput.Script = input["script"].(string)
	scriptAPIMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	scriptAPIMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())

	return scriptAPIMonitorUpdateInput
}

func buildSyntheticsScriptBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptBrowserMonitorInput {
	scriptBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateScriptBrowserMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	scriptBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = input["enable_screenshot_on_failure_and_script"].(bool)
	scriptBrowserMonitorUpdateInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])
	scriptBrowserMonitorUpdateInput.Name = input["name"].(string)
	scriptBrowserMonitorUpdateInput.Period = synthetics.SyntheticsMonitorPeriod(input["frequency"].(string))
	scriptBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	scriptBrowserMonitorUpdateInput.Runtime.ScriptLanguage = input["script_language"].(string)
	scriptBrowserMonitorUpdateInput.Runtime.RuntimeType = input["runtime_type"].(string)
	scriptBrowserMonitorUpdateInput.Script = input["script"].(string)
	scriptBrowserMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	scriptBrowserMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())

	return scriptBrowserMonitorUpdateInput
}

func buildSyntheticsSimpleBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleBrowserMonitorInput {
	simpleBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleBrowserMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	simpleBrowserMonitorUpdateInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"].(*schema.Set).List())
	simpleBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = input["enable_screenshot_on_failure_and_script"].(bool)
	simpleBrowserMonitorUpdateInput.AdvancedOptions.ResponseValidationText = input["validation_string"].(string)
	simpleBrowserMonitorUpdateInput.AdvancedOptions.UseTlsValidation = input["verify_ssl"].(bool)
	simpleBrowserMonitorUpdateInput.Locations = expandSyntheticsLocations(input["locations"])
	simpleBrowserMonitorUpdateInput.Name = input["name"].(string)
	simpleBrowserMonitorUpdateInput.Period = synthetics.SyntheticsMonitorPeriod(input["frequency"].(string))
	simpleBrowserMonitorUpdateInput.Runtime.RuntimeType = input["runtime_type"].(string)
	simpleBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	simpleBrowserMonitorUpdateInput.Runtime.ScriptLanguage = input["scriptLanguage"].(string)
	simpleBrowserMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	simpleBrowserMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())
	simpleBrowserMonitorUpdateInput.Uri = input["uri"].(string)

	return simpleBrowserMonitorUpdateInput
}

func buildSyntheticsSimpleMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleMonitorInput {
	simpleMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	simpleMonitorUpdateInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"].(*schema.Set).List())
	simpleMonitorUpdateInput.AdvancedOptions.RedirectIsFailure = input["treat_redirect_as_failure"].(bool)
	simpleMonitorUpdateInput.AdvancedOptions.ResponseValidationText = input["validation_string"].(string)
	simpleMonitorUpdateInput.AdvancedOptions.ShouldBypassHeadRequest = input["bypass_head_request"].(bool)
	simpleMonitorUpdateInput.AdvancedOptions.UseTlsValidation = input["verify_ssl"].(bool)
	simpleMonitorUpdateInput.Locations = expandSyntheticsLocations(input["locations"])
	simpleMonitorUpdateInput.Name = input["name"].(string)
	simpleMonitorUpdateInput.Period = synthetics.SyntheticsMonitorPeriod(input["frequency"].(string))
	simpleMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["SyntheticsMonitorStatus"].(string))
	simpleMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].(*schema.Set).List())
	simpleMonitorUpdateInput.Uri = input["uri"].(string)

	return simpleMonitorUpdateInput
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

	switch monitorType {

	case "SIMPLE":
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

	case "BROWSER":
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

	case "SCRIPT_API":
		scriptedAPIMonitorUpdateInput := buildSyntheticsScriptAPIMonitorUpdateStruct(d)
		resp, err := client.Synthetics.SyntheticsUpdateScriptAPIMonitorWithContext(ctx, guid, scriptedAPIMonitorUpdateInput)
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

	case "SCRIPT_BROWSER":
		scriptedBrowserMonitorUpdateInput := buildSyntheticsScriptBrowserMonitorUpdateStruct(d)

		resp, err := client.Synthetics.SyntheticsUpdateScriptBrowserMonitorWithContext(ctx, guid, scriptedBrowserMonitorUpdateInput)

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

	}

	return resourceNewRelicSyntheticsMonitorRead(ctx, d, meta)
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
