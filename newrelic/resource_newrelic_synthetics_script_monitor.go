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

type SyntheticsMonitorType string

// nolint:revive
var SyntheticsMonitorTypes = struct {
	SIMPLE         SyntheticsMonitorType
	BROWSER        SyntheticsMonitorType
	SCRIPT_API     SyntheticsMonitorType
	SCRIPT_BROWSER SyntheticsMonitorType
}{
	SIMPLE:         "SIMPLE",
	BROWSER:        "BROWSER",
	SCRIPT_API:     "SCRIPT_API",
	SCRIPT_BROWSER: "SCRIPT_BROWSER",
}

func resourceNewRelicSyntheticsScriptMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsScriptMonitorCreate,
		ReadContext:   resourceNewRelicSyntheticsScriptMonitorRead,
		UpdateContext: resourceNewRelicSyntheticsScriptMonitorUpdate,
		DeleteContext: resourceNewRelicSyntheticsScriptMonitorDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: importSyntheticsMonitorScript,
		// },
		Schema: mergeSchemas(
			syntheticsMonitorCommonSchema(),
			syntheticsScriptMonitorCommonSchema(),
			syntheticsScriptMonitorLocationsSchema(),
			syntheticsScriptBrowserMonitorAdvancedOptionsSchema(),
		),
	}
}

// TODO: Find a good home for this since it can be shared by all monitor schemas.
//
// Returns the common schema attributes shared by all Synthetics monitor types.
func syntheticsMonitorCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The title of this monitor.",
		},
		"status": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The monitor status (i.e. ENABLED, MUTED, DISABLED).",
			ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorStatuses(), false),
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
		"period": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "The interval at which this monitor should run. Valid values are EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, or EVERY_DAY.",
			ValidateFunc: validation.StringInSlice(listValidSyntheticsMonitorPeriods(), false),
		},
	}
}

// Scripted API and scripted browser monitors can have public and/or private locations set.
func syntheticsScriptMonitorLocationsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"locations_private": {
			Type:        schema.TypeSet,
			Description: "",
			Optional:    true, // Note: Optional
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"guid": {
						Type:        schema.TypeString,
						Description: "The unique identifier for the Synthetics private location in New Relic.",
						Required:    true,
					},
					"vse_password": {
						Type:        schema.TypeString,
						Description: "The location's Verified Script Execution password (Only necessary if Verified Script Execution is enabled for the location).",
						Optional:    true,
						Sensitive:   true,
					},
				},
			},
		},
		"locations_public": {
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			MinItems:    1,
			Optional:    true,
			Description: "The public location(s) that the monitor will run jobs from.",
		},
	}
}

// Scripted browser monitors have advanced options, but scripted API monitors do not.
func syntheticsScriptBrowserMonitorAdvancedOptionsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enable_screenshot_on_failure_and_script": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Capture a screenshot during job execution.",
		},
	}
}

// Returns common schema attributes shared by both scripted browser and scripted API monitors.
func syntheticsScriptMonitorCommonSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"guid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The unique entity identifier of the monitor in New Relic.",
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			Description:  "The monitor type. Valid values are SCRIPT_BROWSER, and SCRIPT_API.",
			ValidateFunc: validation.StringInSlice(listValidSyntheticsScriptMonitorTypes(), false),
		},
		"script": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The script that the monitor runs.",
		},
		"script_language": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The programing language that should execute the script.",
		},
		"runtime_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The runtime type that the monitor will run.",
		},
		"runtime_type_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The specific semver version of the runtime type.",
		},
	}
}

func flattenSyntheticsScriptAPIMonitor(monitor synthetics.SyntheticsScriptAPIMonitor, d *schema.ResourceData) diag.Diagnostics {
	_ = d.Set("guid", string(monitor.GUID))
	_ = d.Set("locations_public", monitor.Locations.Public)
	_ = d.Set("name", monitor.Name)
	_ = d.Set("period", string(monitor.Period))
	_ = d.Set("status", monitor.Status)
	_ = d.Set("script_language", monitor.Runtime.ScriptLanguage)
	_ = d.Set("runtime_type", monitor.Runtime.RuntimeType)
	_ = d.Set("runtime_type_version", string(monitor.Runtime.RuntimeTypeVersion))

	return nil
}

func resourceNewRelicSyntheticsScriptMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	var diags diag.Diagnostics
	monitorType := d.Get("type").(string)

	var resp *synthetics.SyntheticsScriptAPIMonitorCreateMutationResult
	var err error
	switch monitorType {
	case string(SyntheticsMonitorTypes.SCRIPT_API):
		monitorInput := buildSyntheticsScriptAPIMonitorInput(d)
		resp, err = client.Synthetics.SyntheticsCreateScriptAPIMonitorWithContext(ctx, accountID, monitorInput)
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

		flattenSyntheticsScriptAPIMonitor(resp.Monitor, d)

	case string(SyntheticsMonitorTypes.SCRIPT_BROWSER):
		monitorInput := buildSyntheticsScriptBrowserMonitorInput(d)
		resp, err := client.Synthetics.SyntheticsCreateScriptBrowserMonitorWithContext(ctx, accountID, monitorInput)
		if err != nil {
			diag.FromErr(err)
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

// WIP - This won't work. The entity endpoints don't return all the data we need to set all the attributes.
func resourceNewRelicSyntheticsScriptMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	setCommonSyntheticsScriptMonitorAttributes(resp, d)
	return nil
}

//func to set output values in the read func.
func setCommonSyntheticsScriptMonitorAttributes(v *entities.EntityInterface, d *schema.ResourceData) {
	switch e := (*v).(type) {

	case *entities.SyntheticMonitorEntity:
		_ = d.Set("name", e.Name)
		_ = d.Set("type", e.MonitorType)
		_ = d.Set("uri", e.MonitoredURL)

	}
}

func resourceNewRelicSyntheticsScriptMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var diags diag.Diagnostics
	monitorType := d.Get("type").(string)

	guid := synthetics.EntityGUID(d.Id())

	switch monitorType {
	case string(SyntheticsMonitorTypes.SCRIPT_API):
		monitorInput := buildSyntheticsScriptAPIMonitorUpdateInput(d)
		resp, err := client.Synthetics.SyntheticsUpdateScriptAPIMonitorWithContext(ctx, guid, monitorInput)
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
	case string(SyntheticsMonitorTypes.SCRIPT_BROWSER):
		monitorInput := buildSyntheticsScriptBrowserUpdateInput(d)
		resp, err := client.Synthetics.SyntheticsUpdateScriptBrowserMonitorWithContext(ctx, guid, monitorInput)
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
	return nil
}

func resourceNewRelicSyntheticsScriptMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

// TODO: Move to proper helper file
func mergeSchemas(schemas ...map[string]*schema.Schema) map[string]*schema.Schema {
	schema := map[string]*schema.Schema{}
	for _, s := range schemas {
		for k, v := range s {
			schema[k] = v
		}
	}
	return schema
}
