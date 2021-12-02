package newrelic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsMonitorScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicSyntheticsMonitorScriptCreate,
		ReadContext:   resourceNewRelicSyntheticsMonitorScriptRead,
		UpdateContext: resourceNewRelicSyntheticsMonitorScriptUpdate,
		DeleteContext: resourceNewRelicSyntheticsMonitorScriptDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importSyntheticsMonitorScript,
		},
		Schema: map[string]*schema.Schema{
			"monitor_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the monitor to attach the script to.",
			},
			"text": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The plaintext representing the monitor script.",
			},
			"location": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of locations for a monitor script.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The monitor script location name",
						},
						"hmac": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The HMAC for the monitor script location. Use only one of `hmac` or `vse_password.`",
						},
						"vse_password": {
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							Description: "The password for the monitor script location used to calculate HMAC. Use only one of `vse_password` or `hmac.`",
						},
					},
				},
			},
		},
	}
}

func importSyntheticsMonitorScript(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	_ = d.Set("monitor_id", d.Id())
	return []*schema.ResourceData{d}, nil
}

func buildSyntheticsMonitorScriptStruct(d *schema.ResourceData) (*synthetics.MonitorScript, error) {
	locations, err := expandMonitorScriptLocations(d.Get("location").([]interface{}), d)
	if err != nil {
		return nil, err
	}

	script := synthetics.MonitorScript{
		Text:      d.Get("text").(string),
		Locations: locations,
	}

	return &script, err
}

func resourceNewRelicSyntheticsMonitorScriptCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	id := d.Get("monitor_id").(string)
	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", id)

	script, scriptErr := buildSyntheticsMonitorScriptStruct(d)
	if scriptErr != nil {
		return diag.FromErr(scriptErr)
	}

	_, err := client.Synthetics.UpdateMonitorScriptWithContext(ctx, id, *script)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceNewRelicSyntheticsMonitorScriptRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Synthetics script %s", d.Id())

	script, err := client.Synthetics.GetMonitorScriptWithContext(ctx, d.Id())
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	_ = d.Set("text", script.Text)

	return nil
}

func resourceNewRelicSyntheticsMonitorScriptUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Creating New Relic Synthetics monitor script %s", d.Id())

	script, scriptErr := buildSyntheticsMonitorScriptStruct(d)
	if scriptErr != nil {
		return diag.FromErr(scriptErr)
	}

	_, err := client.Synthetics.UpdateMonitorScriptWithContext(ctx, d.Id(), *script)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Id())
	return resourceNewRelicSyntheticsMonitorScriptRead(ctx, d, meta)
}

func resourceNewRelicSyntheticsMonitorScriptDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic Synthetics monitor script %s", d.Id())

	script := synthetics.MonitorScript{
		Text:      " ",
		Locations: make([]synthetics.MonitorScriptLocation, 0),
	}

	if _, err := client.Synthetics.UpdateMonitorScriptWithContext(ctx, d.Id(), script); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandMonitorScriptLocations(cfg []interface{}, d *schema.ResourceData) ([]synthetics.MonitorScriptLocation, error) {
	var locations []synthetics.MonitorScriptLocation

	if len(cfg) == 0 {
		return locations, nil
	}

	locations = make([]synthetics.MonitorScriptLocation, 0, len(cfg))

	for _, l := range cfg {
		cfgLocation := l.(map[string]interface{})

		location := synthetics.MonitorScriptLocation{}

		if n, ok := cfgLocation["name"]; ok {
			location.Name = n.(string)

			if v, ok := cfgLocation["vse_password"]; ok && v != "" {
				if h, ok := cfgLocation["hmac"]; ok && h != "" {
					return nil, fmt.Errorf("only set one of either `hmac` or `vse_password`")
				}

				key := []byte(v.(string))
				h := hmac.New(sha256.New, key)
				h.Write([]byte(d.Get("text").(string)))
				encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
				location.HMAC = encoded

				if h, ok := cfgLocation["hmac"]; ok && h != "" {
					if v, ok := cfgLocation["vse_password"]; ok && v != "" {
						return nil, fmt.Errorf("only set one of either `hmac` or `vse_password`")
					}
					location.HMAC = h.(string)
				}
			}
		}
		locations = append(locations, location)
	}
	return locations, nil
}