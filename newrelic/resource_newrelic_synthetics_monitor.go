package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"log"

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
				Elem: schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "",
						},
						"values": {
							Type:        schema.TypeSet,
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
				Elem: schema.Resource{
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
	if status, ok := d.GetOk("SyntheticsMonitorStatus"); ok {
		monitorInputs["SyntheticsMonitorStatus"] = status
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
	return monitorInputs
}

func buildSyntheticsScriptAPIMonitorStruct(d *schema.ResourceData) synthetics.SyntheticsCreateScriptAPIMonitorInput {
	scriptAPIMonitorInput := synthetics.SyntheticsCreateScriptAPIMonitorInput{}
	input := buildSyntheticsMonitorBase(d)
	scriptAPIMonitorInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])
	scriptAPIMonitorInput.Name = input["name"].(string)
	scriptAPIMonitorInput.Period = synthetics.SyntheticsMonitorPeriod(input["frequency"].(string))
	scriptAPIMonitorInput.Runtime.RuntimeType = input["runtime_type"].(string)
	scriptAPIMonitorInput.Runtime.ScriptLanguage = input["script_language"].(string)
	scriptAPIMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	scriptAPIMonitorInput.Script = input["script"].(string)
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
	scriptBrowserMonitorInput.Script = input["script"].(string)
	scriptBrowserMonitorInput.Runtime.ScriptLanguage = input["script_language"].(string)
	return scriptBrowserMonitorInput
}

func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleBrowserMonitorInput {
	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{}
	if customHeaders, ok := d.GetOk("custom_headers"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(customHeaders.(*schema.Set).List())
	}
	if screenshot, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = screenshot.(bool)
	}
	if respValidateText, ok := d.GetOk("validation_string"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.ResponseValidationText = respValidateText.(string)
	}
	if tlsValidations, ok := d.GetOk("verify_ssl"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.UseTlsValidation = tlsValidations.(bool)
	}
	if locations, ok := d.GetOk("locations"); ok {
		simpleBrowserMonitorInput.Locations = expandSyntheticsLocations(locations)
	}
	if name, ok := d.GetOk("name"); ok {
		simpleBrowserMonitorInput.Name = name.(string)
	}
	if period, ok := d.GetOk("frequency"); ok {
		simpleBrowserMonitorInput.Period = synthetics.SyntheticsMonitorPeriod(period.(string))
	}
	if runtimeType, ok := d.GetOk("runtime_type"); ok {
		simpleBrowserMonitorInput.Runtime.RuntimeType = runtimeType.(string)
	}
	if runtimeTypeVersion, ok := d.GetOk("runtime_type_version"); ok {
		simpleBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
	}
	if scriptLanguage, ok := d.GetOk("scriptLanguage"); ok {
		simpleBrowserMonitorInput.Runtime.ScriptLanguage = scriptLanguage.(string)
	}
	if status, ok := d.GetOk("SyntheticsMonitorStatus"); ok {
		simpleBrowserMonitorInput.Status = synthetics.SyntheticsMonitorStatus(status.(string))
	}
	if tags, ok := d.GetOk("tags"); ok {
		simpleBrowserMonitorInput.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}
	if uri, ok := d.GetOk("uri"); ok {
		simpleBrowserMonitorInput.Uri = uri.(string)
	}
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
	if customHeaders, ok := d.GetOk("custom_headers"); ok {
		simpleMonitorInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(customHeaders.(*schema.Set).List())
	}
	if redirectFailure, ok := d.GetOk("treat_redirect_as_failure"); ok {
		simpleMonitorInput.AdvancedOptions.RedirectIsFailure = redirectFailure.(bool)
	}
	if respValidateText, ok := d.GetOk("validation_string"); ok {
		simpleMonitorInput.AdvancedOptions.ResponseValidationText = respValidateText.(string)
	}
	if pass, ok := d.GetOk("bypass_head_request"); ok {
		simpleMonitorInput.AdvancedOptions.ShouldBypassHeadRequest = pass.(bool)
	}
	if tlsValidations, ok := d.GetOk("verify_ssl"); ok {
		simpleMonitorInput.AdvancedOptions.UseTlsValidation = tlsValidations.(bool)
	}
	if locations, ok := d.GetOk("locations"); ok {
		simpleMonitorInput.Locations = expandSyntheticsLocations(locations)
	}
	if name, ok := d.GetOk("name"); ok {
		simpleMonitorInput.Name = name.(string)
	}
	if period, ok := d.GetOk("frequency"); ok {
		simpleMonitorInput.Period = synthetics.SyntheticsMonitorPeriod(period.(string))
	}
	if status, ok := d.GetOk("SyntheticsMonitorStatus"); ok {
		simpleMonitorInput.Status = synthetics.SyntheticsMonitorStatus(status.(string))
	}
	if tags, ok := d.GetOk("tags"); ok {
		simpleMonitorInput.Tags = expandSyntheticsTags(tags.(*schema.Set).List())
	}
	if uri, ok := d.GetOk("uri"); ok {
		simpleMonitorInput.Uri = uri.(string)
	}
	return simpleMonitorInput
}

func buildSyntheticsUpdateMonitorArgs(d *schema.ResourceData) *synthetics.Monitor {
	monitor := synthetics.Monitor{
		ID:           d.Id(),
		Name:         d.Get("name").(string),
		Type:         synthetics.MonitorType(d.Get("type").(string)),
		Frequency:    uint(d.Get("frequency").(int)),
		Status:       synthetics.MonitorStatusType(d.Get("status").(string)),
		SLAThreshold: d.Get("sla_threshold").(float64),
	}

	if uri, ok := d.GetOk("uri"); ok {
		monitor.URI = uri.(string)
	}

	locationsRaw := d.Get("locations").(*schema.Set)
	locations := make([]string, locationsRaw.Len())
	for i, v := range locationsRaw.List() {
		locations[i] = fmt.Sprint(v)
	}

	if validationString, ok := d.GetOk("validation_string"); ok {
		monitor.Options.ValidationString = validationString.(string)
	}

	if verifySSL, ok := d.GetOkExists("verify_ssl"); ok {
		monitor.Options.VerifySSL = verifySSL.(bool)
	}

	if bypassHeadRequest, ok := d.GetOkExists("bypass_head_request"); ok {
		monitor.Options.BypassHEADRequest = bypassHeadRequest.(bool)
	}

	if treatRedirectAsFailure, ok := d.GetOkExists("treat_redirect_as_failure"); ok {
		monitor.Options.TreatRedirectAsFailure = treatRedirectAsFailure.(bool)
	}

	monitor.Locations = locations
	return &monitor
}

func readSyntheticsMonitorStruct(monitor *synthetics.Monitor, d *schema.ResourceData) {
	_ = d.Set("name", monitor.Name)
	_ = d.Set("type", monitor.Type)
	_ = d.Set("frequency", monitor.Frequency)
	_ = d.Set("uri", monitor.URI)
	_ = d.Set("locations", monitor.Locations)
	_ = d.Set("status", monitor.Status)
	_ = d.Set("sla_threshold", monitor.SLAThreshold)
	_ = d.Set("verify_ssl", monitor.Options.VerifySSL)
	_ = d.Set("validation_string", monitor.Options.ValidationString)
	_ = d.Set("bypass_head_request", monitor.Options.BypassHEADRequest)
	_ = d.Set("treat_redirect_as_failure", monitor.Options.TreatRedirectAsFailure)
}

func resourceNewRelicSyntheticsMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	/////////////////////////////
	providerConfig := meta.(*ProviderConfig)

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
		gUID := resp.Monitor.GUID
		err = d.Set("guid", gUID)
		if err != nil {
			return diag.FromErr(err)
		}

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
		gUID := resp.Monitor.GUID
		err = d.Set("guid", gUID)
		if err != nil {
			return diag.FromErr(err)
		}

	case "SCRIPT_API":
		scriptedApiMonitorInput := buildSyntheticsScriptAPIMonitorStruct(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptAPIMonitorWithContext(ctx, accountID, scriptedApiMonitorInput)
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
		gUID := resp.Monitor.GUID
		err = d.Set("guid", gUID)
		if err != nil {
			return diag.FromErr(err)
		}

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
		gUID := resp.Monitor.GUID

		err = d.Set("guid", gUID)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	///////////////////////////

	//monitorStruct := buildSyntheticsMonitorStruct(d)
	//
	//log.Printf("[INFO] Creating New Relic Synthetics monitor %s", monitorStruct.Name)
	//
	//monitor, err := client.Synthetics.CreateMonitorWithContext(ctx, monitorStruct)
	//if err != nil {
	//	return diag.FromErr(err)
	//}
	//
	//d.SetId(monitor.ID)
	return resourceNewRelicSyntheticsMonitorRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	resp, err := client.Entities.GetEntity(common.EntityGUID(d.Id()))
	setCommonSyntheticsMonitorAttributes(resp, d)
	monitor, err := client.Synthetics.GetMonitorWithContext(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	readSyntheticsMonitorStruct(monitor, d)

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

func resourceNewRelicSyntheticsMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	log.Printf("[INFO] Updating New Relic Synthetics monitor %s", d.Id())

	//providerConfig := meta.(*ProviderConfig)

	//accountID := selectAccountID(providerConfig, d)
	//
	//resp, err := client.Synthetics.SyntheticsUpdateSimpleMonitorWithContext()
	_, err := client.Synthetics.UpdateMonitorWithContext(ctx, *buildSyntheticsUpdateMonitorArgs(d))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicSyntheticsMonitorRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	if err := client.Synthetics.DeleteMonitorWithContext(ctx, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
