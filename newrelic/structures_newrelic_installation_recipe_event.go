package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/installevents"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

func expandInstallationRecipeEvent(d *schema.ResourceData) (*installevents.InstallationRecipeStatus, error) {
	input := installevents.InstallationRecipeStatus{
		CliVersion:                     d.Get("cli_version").(string),
		Complete:                       d.Get("complete").(bool),
		DisplayName:                    d.Get("display_name").(string),
		EntityGUID:                     common.EntityGUID(d.Get("entity_guid").(string)),
		HostName:                       d.Get("host_name").(string),
		KernelArch:                     d.Get("kernel_arch").(string),
		KernelVersion:                  d.Get("kernel_version").(string),
		LogFilePath:                    d.Get("log_file_path").(string),
		Name:                           d.Get("name").(string),
		Os:                             d.Get("os").(string),
		Platform:                       d.Get("platform").(string),
		PlatformFamily:                 d.Get("platform_family").(string),
		PlatformVersion:                d.Get("platform_version").(string),
		Status:                         installevents.InstallationRecipeStatusType(d.Get("status").(string)),
		TargetedInstall:                d.Get("targeted_install").(bool),
		ValidationDurationMilliseconds: int64(d.Get("validation_duration_milliseconds").(int)),
	}

	if v, ok := d.GetOk("install_id"); ok {
		input.InstallId = v.(string)
	}

	if v, ok := d.GetOk("install_library_version"); ok {
		input.InstallLibraryVersion = v.(string)
	}

	if v, ok := d.GetOk("redirect_url"); ok {
		input.RedirectURL = v.(string)
	}

	if v, ok := d.GetOk("task_path"); ok {
		input.TaskPath = v.(string)
	}

	if v, ok := d.GetOk("timestamp"); ok {
		input.Timestamp = nrtime.EpochSeconds(float64(v.(int)))
	}

	if v, ok := d.GetOk("error"); ok {
		items := v.([]interface{})
		if len(items) > 0 && items[0] != nil {
			input.Error = expandInstallationStatusErrorInput(items[0].(map[string]interface{}))
		}
	}

	return &input, nil
}

func expandInstallationStatusErrorInput(cfg map[string]interface{}) installevents.InstallationStatusErrorInput {
	out := installevents.InstallationStatusErrorInput{}

	if v, ok := cfg["details"].(string); ok {
		out.Details = v
	}

	if v, ok := cfg["message"].(string); ok {
		out.Message = v
	}

	if v, ok := cfg["optimized_message"].(string); ok {
		out.OptimizedMessage = v
	}

	return out
}

func expandInstallationRecipeEventUpdate(d *schema.ResourceData) (*installevents.InstallationRecipeStatus, error) {
	return expandInstallationRecipeEvent(d)
}

var _ = (*schema.ResourceData)(nil)
func flattenInstallationRecipeEvent(result *installevents.InstallationRecipeEvent, d *schema.ResourceData) error {
	_ = d.Set("cli_version", result.CliVersion)
	_ = d.Set("complete", result.Complete)
	_ = d.Set("display_name", result.DisplayName)
	_ = d.Set("entity_guid", string(result.EntityGUID))
	_ = d.Set("host_name", result.HostName)
	_ = d.Set("install_id", result.InstallId)
	_ = d.Set("install_library_version", result.InstallLibraryVersion)
	_ = d.Set("kernel_arch", result.KernelArch)
	_ = d.Set("kernel_version", result.KernelVersion)
	_ = d.Set("log_file_path", result.LogFilePath)
	_ = d.Set("name", result.Name)
	_ = d.Set("os", result.Os)
	_ = d.Set("platform", result.Platform)
	_ = d.Set("platform_family", result.PlatformFamily)
	_ = d.Set("platform_version", result.PlatformVersion)
	_ = d.Set("redirect_url", result.RedirectURL)
	_ = d.Set("status", string(result.Status))
	_ = d.Set("targeted_install", result.TargetedInstall)
	_ = d.Set("task_path", result.TaskPath)
	_ = d.Set("timestamp", int(result.Timestamp))
	_ = d.Set("validation_duration_milliseconds", int(result.ValidationDurationMilliseconds))

	if err := d.Set("error", flattenInstallationRecipeEventError(result.Error)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting `error`: %v", err)
	}

	return nil
}

func flattenInstallationRecipeEventError(e installevents.InstallationStatusError) []interface{} {
	if e.Details == "" && e.Message == "" && e.OptimizedMessage == "" {
		return []interface{}{}
	}
	m := map[string]interface{}{
		"details":           e.Details,
		"message":           e.Message,
		"optimized_message": e.OptimizedMessage,
	}
	return []interface{}{m}
}