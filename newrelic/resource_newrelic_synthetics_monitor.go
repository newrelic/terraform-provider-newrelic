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
			"period": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The URI for the monitor to hit.",
				ValidateFunc: validation.StringInSlice(getSyntheticsMonitorPeriodTypesAsStrings(), false),
			},
			"uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URI for the monitor to hit.",
				// TODO: ValidateFunc (required if SIMPLE or BROWSER)
			},

			// TODO: Locations needs to include both private and public
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
				Optional:    true,
				Description: "",
			},
			"runtime_type_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			"script_language": {
				Type:        schema.TypeString,
				Optional:    true,
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
							Type:        schema.TypeList,
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
			"script": {
				Type:        schema.TypeString,
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
	inputBase := expandSyntheticsMonitorBase(d)

	scriptAPIMonitorInput := synthetics.SyntheticsCreateScriptAPIMonitorInput{}
	scriptAPIMonitorInput.Name = inputBase.Name
	scriptAPIMonitorInput.Period = inputBase.Period
	scriptAPIMonitorInput.Status = inputBase.Status
	scriptAPIMonitorInput.Tags = inputBase.Tags

	if v, ok := d.GetOk("script"); ok {
		scriptAPIMonitorInput.Script = v.(string)
	}

	if v, ok := d.GetOk("script_language"); ok {
		scriptAPIMonitorInput.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		scriptAPIMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	// Continue to work through....

	// scriptAPIMonitorInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])

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

func buildSyntheticsScriptBrowserMonitorStruct(d *schema.ResourceData) synthetics.SyntheticsCreateScriptBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	scriptBrowserMonitorInput := synthetics.SyntheticsCreateScriptBrowserMonitorInput{}
	scriptBrowserMonitorInput.Name = inputBase.Name
	scriptBrowserMonitorInput.Period = inputBase.Period
	scriptBrowserMonitorInput.Status = inputBase.Status
	scriptBrowserMonitorInput.Tags = inputBase.Tags

	if v, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		scriptBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = v.(bool)
	}

	if v, ok := d.GetOk("script"); ok {
		scriptBrowserMonitorInput.Script = v.(string)
	}

	if v, ok := d.GetOk("script_language"); ok {
		scriptBrowserMonitorInput.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		scriptBrowserMonitorInput.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		scriptBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	// Continue to work through...

	// scriptBrowserMonitorInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])

	return scriptBrowserMonitorInput
}

func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) synthetics.SyntheticsCreateSimpleBrowserMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorInput := synthetics.SyntheticsCreateSimpleBrowserMonitorInput{}
	simpleBrowserMonitorInput.Locations = inputBase.Locations
	simpleBrowserMonitorInput.Name = inputBase.Name
	simpleBrowserMonitorInput.Period = inputBase.Period
	simpleBrowserMonitorInput.Status = inputBase.Status
	simpleBrowserMonitorInput.Tags = inputBase.Tags
	simpleBrowserMonitorInput.Uri = inputBase.Uri

	if v, ok := d.GetOk("enable_screenshot_on_failure_and_script"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = v.(bool)
	}

	if v, ok := d.GetOk("script_language"); ok {
		simpleBrowserMonitorInput.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		simpleBrowserMonitorInput.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		simpleBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	return simpleBrowserMonitorInput
}

func expandCustomHeaders(headers interface{}) []synthetics.SyntheticsCustomHeaderInput {
	output := make([]synthetics.SyntheticsCustomHeaderInput, len(headers.([]interface{})))
	for i, v := range headers.([]interface{}) {
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
	inputBase := expandSyntheticsMonitorBase(d)

	simpleMonitorInput := synthetics.SyntheticsCreateSimpleMonitorInput{}
	simpleMonitorInput.Locations = inputBase.Locations
	simpleMonitorInput.Name = inputBase.Name
	simpleMonitorInput.Period = inputBase.Period
	simpleMonitorInput.Status = inputBase.Status
	simpleMonitorInput.Tags = inputBase.Tags
	simpleMonitorInput.Uri = inputBase.Uri

	if v, ok := d.GetOk("treat_redirect_as_failure"); ok {
		simpleMonitorInput.AdvancedOptions.RedirectIsFailure = v.(bool)
	}

	if v, ok := d.GetOk("validation_string"); ok {
		simpleMonitorInput.AdvancedOptions.ResponseValidationText = v.(string)
	}

	if v, ok := d.GetOk("bypass_head_request"); ok {
		simpleMonitorInput.AdvancedOptions.ShouldBypassHeadRequest = v.(bool)
	}

	if v, ok := d.GetOk("verify_ssl"); ok {
		simpleMonitorInput.AdvancedOptions.UseTlsValidation = v.(bool)
	}

	// Continue to work through...

	// simpleMonitorInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"])

	return simpleMonitorInput
}

// TODO: Reduce the cyclomatic complexity of this func
// nolint:gocyclo
func resourceNewRelicSyntheticsMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	var diags diag.Diagnostics

	monitorType := d.Get("type")

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

	resp, err := client.Entities.GetEntity(common.EntityGUID(d.Id()))
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
	case *entities.SyntheticMonitorEntityOutline:
		_ = d.Set("guid", e.GUID)
		_ = d.Set("name", e.Name)
		_ = d.Set("type", e.MonitorType)
		_ = d.Set("period", e.Period)
		_ = d.Set("tags", e.Tags)
		_ = d.Set("uri", e.MonitoredURL)
	}
}

func buildSyntheticsScriptAPIMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptAPIMonitorInput {
	inputBase := expandSyntheticsMonitorBase(d)

	scriptAPIMonitorUpdateInput := synthetics.SyntheticsUpdateScriptAPIMonitorInput{}
	scriptAPIMonitorUpdateInput.Name = inputBase.Name
	scriptAPIMonitorUpdateInput.Period = inputBase.Period
	scriptAPIMonitorUpdateInput.Status = inputBase.Status
	scriptAPIMonitorUpdateInput.Tags = inputBase.Tags

	if v, ok := d.GetOk("script"); ok {
		scriptAPIMonitorUpdateInput.Script = v.(string)
	}

	if v, ok := d.GetOk("script_language"); ok {
		scriptAPIMonitorUpdateInput.Runtime.ScriptLanguage = v.(string)
	}

	if v, ok := d.GetOk("runtime_type"); ok {
		scriptAPIMonitorUpdateInput.Runtime.RuntimeType = v.(string)
	}

	if v, ok := d.GetOk("runtime_type_version"); ok {
		scriptAPIMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(v.(string))
	}

	// Continue to work through ...

	// scriptAPIMonitorUpdateInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])\

	return scriptAPIMonitorUpdateInput
}

func buildSyntheticsScriptBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateScriptBrowserMonitorInput {
	scriptBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateScriptBrowserMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	scriptBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = input["enable_screenshot_on_failure_and_script"].(bool)
	scriptBrowserMonitorUpdateInput.Locations = expandSyntheticsScriptMonitorLocations(input["locations"])
	scriptBrowserMonitorUpdateInput.Name = input["name"].(string)
	scriptBrowserMonitorUpdateInput.Period = periodConvIntToString(input["frequency"])
	scriptBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	scriptBrowserMonitorUpdateInput.Runtime.ScriptLanguage = input["script_language"].(string)
	scriptBrowserMonitorUpdateInput.Runtime.RuntimeType = input["runtime_type"].(string)
	scriptBrowserMonitorUpdateInput.Script = input["script"].(string)
	scriptBrowserMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	scriptBrowserMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].([]interface{}))

	return scriptBrowserMonitorUpdateInput
}

func buildSyntheticsSimpleBrowserMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleBrowserMonitorInput {
	simpleBrowserMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleBrowserMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	simpleBrowserMonitorUpdateInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"])
	simpleBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = input["enable_screenshot_on_failure_and_script"].(bool)
	simpleBrowserMonitorUpdateInput.AdvancedOptions.ResponseValidationText = input["validation_string"].(string)
	simpleBrowserMonitorUpdateInput.AdvancedOptions.UseTlsValidation = input["verify_ssl"].(bool)
	simpleBrowserMonitorUpdateInput.Locations = expandSyntheticsLocations(input["locations"])
	simpleBrowserMonitorUpdateInput.Name = input["name"].(string)
	simpleBrowserMonitorUpdateInput.Period = periodConvIntToString(input["frequency"])
	simpleBrowserMonitorUpdateInput.Runtime.RuntimeType = input["runtime_type"].(string)
	simpleBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(input["runtime_type_version"].(string))
	simpleBrowserMonitorUpdateInput.Runtime.ScriptLanguage = input["script_language"].(string)
	simpleBrowserMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	simpleBrowserMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].([]interface{}))
	simpleBrowserMonitorUpdateInput.Uri = input["uri"].(string)

	return simpleBrowserMonitorUpdateInput
}

func buildSyntheticsSimpleMonitorUpdateStruct(d *schema.ResourceData) synthetics.SyntheticsUpdateSimpleMonitorInput {
	simpleMonitorUpdateInput := synthetics.SyntheticsUpdateSimpleMonitorInput{}

	input := buildSyntheticsMonitorBase(d)

	simpleMonitorUpdateInput.AdvancedOptions.CustomHeaders = expandCustomHeaders(input["custom_headers"])
	simpleMonitorUpdateInput.AdvancedOptions.RedirectIsFailure = input["treat_redirect_as_failure"].(bool)
	simpleMonitorUpdateInput.AdvancedOptions.ResponseValidationText = input["validation_string"].(string)
	simpleMonitorUpdateInput.AdvancedOptions.ShouldBypassHeadRequest = input["bypass_head_request"].(bool)
	simpleMonitorUpdateInput.AdvancedOptions.UseTlsValidation = input["verify_ssl"].(bool)
	simpleMonitorUpdateInput.Locations = expandSyntheticsLocations(input["locations"])
	simpleMonitorUpdateInput.Name = input["name"].(string)
	simpleMonitorUpdateInput.Period = periodConvIntToString(input["frequency"])
	simpleMonitorUpdateInput.Status = synthetics.SyntheticsMonitorStatus(input["status"].(string))
	simpleMonitorUpdateInput.Tags = expandSyntheticsTags(input["tags"].([]interface{}))
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
