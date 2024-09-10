package newrelic

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func buildSyntheticsSimpleBrowserMonitor(d *schema.ResourceData) (synthetics.SyntheticsCreateSimpleBrowserMonitorInput, error) {
	inputBase := expandSyntheticsMonitorBase(d)
	simpleBrowserMonitorInput := &synthetics.SyntheticsCreateSimpleBrowserMonitorInput{
		Name:            inputBase.Name,
		Period:          inputBase.Period,
		Status:          inputBase.Status,
		Tags:            inputBase.Tags,
		AdvancedOptions: synthetics.SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{},
	}

	//Initializing an empty slice of CustomHeader if no custom headers block is provided in TF config.
	t := expandSyntheticsCustomHeaders(d.Get("custom_header").(*schema.Set).List())
	simpleBrowserMonitorInput.AdvancedOptions.CustomHeaders = &t

	if v, ok := d.GetOk("uri"); ok {
		simpleBrowserMonitorInput.Uri = v.(string)
	}

	if v := d.Get("enable_screenshot_on_failure_and_script"); v != nil {
		e := v.(bool)
		simpleBrowserMonitorInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}

	if v, ok := d.GetOk("locations_public"); ok {
		simpleBrowserMonitorInput.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_private"); ok {
		simpleBrowserMonitorInput.Locations.Private = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("validation_string"); ok {
		simpleBrowserMonitorInput.AdvancedOptions.ResponseValidationText = v.(string)
	}

	if v := d.Get("verify_ssl"); v != nil {
		vs := v.(bool)
		simpleBrowserMonitorInput.AdvancedOptions.UseTlsValidation = &vs
	}

	typedBrowsers := expandSyntheticsBrowsers(d.Get("browsers").(*schema.Set).List())
	simpleBrowserMonitorInput.Browsers = typedBrowsers

	typedDevices := expandSyntheticsDevices(d.Get("devices").(*schema.Set).List())
	simpleBrowserMonitorInput.Devices = typedDevices

	err := buildSyntheticsSimpleBrowserMonitorRuntimeAndDeviceEmulation(d, simpleBrowserMonitorInput)
	if err != nil {
		return *simpleBrowserMonitorInput, err
	}

	return *simpleBrowserMonitorInput, nil
}

func buildSyntheticsSimpleBrowserMonitorRuntimeAndDeviceEmulation(d *schema.ResourceData, simpleBrowserMonitorInput *synthetics.SyntheticsCreateSimpleBrowserMonitorInput) error {
	scriptLang, scriptLangOk := d.GetOk("script_language")
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if scriptLangOk || runtimeTypeOk || runtimeTypeVersionOk {
		simpleBrowserMonitorInput.Runtime = &synthetics.SyntheticsRuntimeInput{}

		if scriptLangOk {
			simpleBrowserMonitorInput.Runtime.ScriptLanguage = scriptLang.(string)
		}

		if runtimeTypeOk {
			simpleBrowserMonitorInput.Runtime.RuntimeType = runtimeType.(string)
		}

		if runtimeTypeVersionOk {
			simpleBrowserMonitorInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
		}
	}

	do, doOk := d.GetOk("device_orientation")
	dt, dtOk := d.GetOk("device_type")

	if !(runtimeTypeOk && runtimeTypeVersionOk) && (doOk && dtOk) {
		return errors.New("device emulation is not supported by legacy runtime")
	}

	if doOk && dtOk {
		simpleBrowserMonitorInput.AdvancedOptions.DeviceEmulation = &synthetics.SyntheticsDeviceEmulationInput{}
		simpleBrowserMonitorInput.AdvancedOptions.DeviceEmulation.DeviceOrientation = synthetics.SyntheticsDeviceOrientation(do.(string))
		simpleBrowserMonitorInput.AdvancedOptions.DeviceEmulation.DeviceType = synthetics.SyntheticsDeviceType(dt.(string))
	} else {
		if doOk || dtOk {
			return errors.New("both device_orientation and device_type should be specified to enable device emulation")
		}
	}

	return nil
}

func buildSyntheticsSimpleBrowserMonitorUpdateStruct(d *schema.ResourceData) (synthetics.SyntheticsUpdateSimpleBrowserMonitorInput, error) {
	inputBase := expandSyntheticsMonitorBase(d)

	simpleBrowserMonitorUpdateInput := &synthetics.SyntheticsUpdateSimpleBrowserMonitorInput{
		Name:            inputBase.Name,
		Period:          inputBase.Period,
		Status:          inputBase.Status,
		Tags:            inputBase.Tags,
		AdvancedOptions: synthetics.SyntheticsSimpleBrowserMonitorAdvancedOptionsInput{},
	}

	//Initializing an empty slice of CustomHeader if no custom headers block is provided in TF config.
	t := expandSyntheticsCustomHeaders(d.Get("custom_header").(*schema.Set).List())
	simpleBrowserMonitorUpdateInput.AdvancedOptions.CustomHeaders = &t

	if v, ok := d.GetOk("uri"); ok {
		simpleBrowserMonitorUpdateInput.Uri = v.(string)
	}

	if v, ok := d.GetOk("locations_public"); ok {
		simpleBrowserMonitorUpdateInput.Locations.Public = expandStringSlice(v.(*schema.Set).List())
	}

	if v, ok := d.GetOk("locations_private"); ok {
		simpleBrowserMonitorUpdateInput.Locations.Private = expandStringSlice(v.(*schema.Set).List())
	}

	if v := d.Get("enable_screenshot_on_failure_and_script"); v != nil {
		e := v.(bool)
		simpleBrowserMonitorUpdateInput.AdvancedOptions.EnableScreenshotOnFailureAndScript = &e
	}

	if v, ok := d.GetOk("validation_string"); ok {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.ResponseValidationText = v.(string)
	}

	if v := d.Get("verify_ssl"); v != nil {
		vs := v.(bool)
		simpleBrowserMonitorUpdateInput.AdvancedOptions.UseTlsValidation = &vs
	}

	typedBrowsers := expandSyntheticsBrowsers(d.Get("browsers").(*schema.Set).List())
	simpleBrowserMonitorUpdateInput.Browsers = typedBrowsers

	typedDevices := expandSyntheticsDevices(d.Get("devices").(*schema.Set).List())
	simpleBrowserMonitorUpdateInput.Devices = typedDevices

	err := buildSyntheticsSimpleBrowserMonitorRuntimeAndDeviceEmulationUpdateStruct(d, simpleBrowserMonitorUpdateInput)
	if err != nil {
		return *simpleBrowserMonitorUpdateInput, err
	}

	return *simpleBrowserMonitorUpdateInput, nil
}

func buildSyntheticsSimpleBrowserMonitorRuntimeAndDeviceEmulationUpdateStruct(d *schema.ResourceData, simpleBrowserMonitorUpdateInput *synthetics.SyntheticsUpdateSimpleBrowserMonitorInput) error {
	scriptLang, scriptLangOk := d.GetOk("script_language")
	runtimeType, runtimeTypeOk := d.GetOk("runtime_type")
	runtimeTypeVersion, runtimeTypeVersionOk := d.GetOk("runtime_type_version")

	if scriptLangOk || runtimeTypeOk || runtimeTypeVersionOk {
		simpleBrowserMonitorUpdateInput.Runtime = &synthetics.SyntheticsRuntimeInput{}

		if scriptLangOk {
			simpleBrowserMonitorUpdateInput.Runtime.ScriptLanguage = scriptLang.(string)
		}

		if runtimeTypeOk {
			simpleBrowserMonitorUpdateInput.Runtime.RuntimeType = runtimeType.(string)
		}

		if runtimeTypeVersionOk {
			simpleBrowserMonitorUpdateInput.Runtime.RuntimeTypeVersion = synthetics.SemVer(runtimeTypeVersion.(string))
		}
	}

	do, doOk := d.GetOk("device_orientation")
	dt, dtOk := d.GetOk("device_type")

	if !(runtimeTypeOk && runtimeTypeVersionOk) && (doOk && dtOk) {
		return errors.New("device emulation is not supported by legacy runtime")
	}

	if doOk && dtOk {
		simpleBrowserMonitorUpdateInput.AdvancedOptions.DeviceEmulation = &synthetics.SyntheticsDeviceEmulationInput{}
		simpleBrowserMonitorUpdateInput.AdvancedOptions.DeviceEmulation.DeviceOrientation = synthetics.SyntheticsDeviceOrientation(do.(string))
		simpleBrowserMonitorUpdateInput.AdvancedOptions.DeviceEmulation.DeviceType = synthetics.SyntheticsDeviceType(dt.(string))
	} else {
		if doOk || dtOk {
			return errors.New("both device_orientation and device_type should be specified to enable device emulation")
		}
	}

	return nil
}
