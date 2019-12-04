package newrelic

import (
	"fmt"
	"log"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	util "github.com/dollarshaveclub/new-relic-synthetics-go/util"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
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
				ValidateFunc: validation.IntInSlice([]int{1, 5, 10, 15, 30, 60, 360, 720, 1440}),
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

func buildSyntheticsMonitorStruct(d *schema.ResourceData) *synthetics.CreateMonitorArgs {
	monitorType := d.Get("type").(string)

	monitor := synthetics.CreateMonitorArgs{
		Type:         monitorType,
		Name:         d.Get("name").(string),
		Frequency:    uint(d.Get("frequency").(int)),
		Status:       d.Get("status").(string),
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

	if monitorType == synthetics.TypeSimple || monitorType == synthetics.TypeBrowser {
		if validationString, ok := d.GetOk("validation_string"); ok {
			monitor.ValidationString = util.StrPtr(validationString.(string))
		}

		verifySSL := d.Get("verify_ssl")
		if verifySSL != nil {
			monitor.VerifySSL = util.BoolPtr(verifySSL.(bool))
		}
	}

	if monitorType == synthetics.TypeSimple {
		bypassHeadRequest := d.Get("bypass_head_request")
		treatRedirectAsFailure := d.Get("treat_redirect_as_failure")

		if bypassHeadRequest != nil {
			monitor.BypassHEADRequest = util.BoolPtr(bypassHeadRequest.(bool))
		}

		if treatRedirectAsFailure != nil {
			monitor.TreatRedirectAsFailure = util.BoolPtr(treatRedirectAsFailure.(bool))
		}
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

	if monitor.Type == synthetics.TypeSimple || monitor.Type == synthetics.TypeBrowser {
		if monitor.VerifySSL != nil {
			d.Set("verify_ssl", monitor.VerifySSL)
		}

		if monitor.ValidationString != nil {
			d.Set("validation_string", monitor.ValidationString)
		}
	}

	if monitor.Type == synthetics.TypeSimple {
		if monitor.BypassHEADRequest != nil {
			d.Set("bypass_head_request", monitor.BypassHEADRequest)
		}

		if monitor.TreatRedirectAsFailure != nil {
			d.Set("treat_redirect_as_failure", monitor.TreatRedirectAsFailure)
		}
	}

	return nil
}

func resourceNewRelicSyntheticsMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics
	monitor := buildSyntheticsMonitorStruct(d)

	log.Printf("[INFO] Creating New Relic Synthetics monitor %s", monitor.Name)

	condition, err := client.CreateMonitor(monitor)

	if err != nil {
		return err
	}

	d.SetId(condition.ID)
	return resourceNewRelicSyntheticsMonitorRead(d, meta)
}

func resourceNewRelicSyntheticsMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Reading New Relic Synthetics monitor %s", d.Id())

	monitor, err := client.GetMonitor(d.Id())

	if err != nil {
		if err == synthetics.ErrMonitorNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return readSyntheticsMonitorStruct(monitor, d)
}

func resourceNewRelicSyntheticsMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics
	monitor := buildSyntheticsMonitorStructUpdate(d)

	log.Printf("[INFO] Updating New Relic Synthetics monitor %s", d.Id())

	_, err := client.UpdateMonitor(d.Id(), monitor)

	if err != nil {
		return err
	}

	return resourceNewRelicSyntheticsMonitorRead(d, meta)
}

func buildSyntheticsMonitorStructUpdate(d *schema.ResourceData) *synthetics.UpdateMonitorArgs {
	monitorType := d.Get("type").(string)

	monitor := synthetics.UpdateMonitorArgs{
		Name:         d.Get("name").(string),
		Frequency:    uint(d.Get("frequency").(int)),
		Status:       d.Get("status").(string),
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
		monitor.ValidationString = util.StrPtr(validationString.(string))
	}

	verifySSL := d.Get("verify_ssl")

	if verifySSL != nil {
		monitor.VerifySSL = util.BoolPtr(verifySSL.(bool))
	}

	if monitorType == synthetics.TypeSimple {
		bypassHeadRequest := d.Get("bypass_head_request")
		treatRedirectAsFailure := d.Get("treat_redirect_as_failure")

		if bypassHeadRequest != nil {
			monitor.BypassHEADRequest = util.BoolPtr(bypassHeadRequest.(bool))
		}

		if treatRedirectAsFailure != nil {
			monitor.TreatRedirectAsFailure = util.BoolPtr(treatRedirectAsFailure.(bool))
		}
	}

	monitor.Locations = locations

	return &monitor
}

func resourceNewRelicSyntheticsMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Deleting New Relic Synthetics monitor %s", d.Id())

	if err := client.DeleteMonitor(d.Id()); err != nil {
		return err
	}

	return nil
}
