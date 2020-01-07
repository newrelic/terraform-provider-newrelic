package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	synthetics "github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func resourceNewRelicSyntheticsMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicSyntheticsMonitorCreate,
		Read:   resourceNewRelicSyntheticsMonitorRead,
		Update: resourceNewRelicSyntheticsMonitorUpdate,
		Delete: resourceNewRelicSyntheticsMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SIMPLE",
					"BROWSER",
					"SCRIPT_API",
					"SCRIPT_BROWSER",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: intInSlice([]int{1, 5, 10, 15, 30, 60, 360, 720, 1440}),
			},
			"uri": {
				Type:     schema.TypeString,
				Optional: true,
				// TODO: ValidateFunc (required if SIMPLE or BROWSER)
			},
			"locations": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MinItems: 1,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ENABLED",
					"MUTED",
					"DISABLED",
				}, false),
			},
			"sla_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
				Default:  7,
			},
			// TODO: ValidationFunc (options only valid if SIMPLE or BROWSER)
			"validation_string": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"verify_ssl": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bypass_head_request": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"treat_redirect_as_failure": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func buildSyntheticsMonitorStruct(d *schema.ResourceData) synthetics.Monitor {
	monitor := synthetics.Monitor{
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
	return monitor
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

func readSyntheticsMonitorStruct(monitor *synthetics.Monitor, d *schema.ResourceData) error {
	d.Set("name", monitor.Name)
	d.Set("type", monitor.Type)
	d.Set("frequency", monitor.Frequency)
	d.Set("uri", monitor.URI)
	d.Set("locations", monitor.Locations)
	d.Set("status", monitor.Status)
	d.Set("sla_threshold", monitor.SLAThreshold)
	d.Set("verify_ssl", monitor.Options.VerifySSL)
	d.Set("validation_string", monitor.Options.ValidationString)
	d.Set("bypass_head_request", monitor.Options.BypassHEADRequest)
	d.Set("treat_redirect_as_failure", monitor.Options.TreatRedirectAsFailure)

	return nil
}

func resourceNewRelicSyntheticsMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics
	monitor := buildSyntheticsMonitorStruct(d)

	log.Printf("[INFO] Creating New Relic Synthetics monitor %s", monitor.Name)

	id, err := client.CreateMonitor(monitor)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceNewRelicSyntheticsMonitorRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	monitor, err := client.GetMonitor(d.Id())
	if err != nil {
		if _, ok := err.(*errors.ErrorNotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return readSyntheticsMonitorStruct(monitor, d)
}

func resourceNewRelicSyntheticsMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics
	monitor := buildSyntheticsUpdateMonitorArgs(d)

	log.Printf("[INFO] Updating New Relic Synthetics monitor %s", d.Id())

	err := client.UpdateMonitor(*monitor)
	if err != nil {
		return err
	}

	return resourceNewRelicSyntheticsMonitorRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	if err := client.DeleteMonitor(d.Id()); err != nil {
		return err
	}

	return nil
}
