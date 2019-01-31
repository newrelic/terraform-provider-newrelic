package newrelic

import (
	"log"

	synthetics "github.com/dollarshaveclub/new-relic-synthetics-go"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SIMPLE",
				//TODO: Validate types, currently only SIMPLE is supported
			},
			"frequency": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: intInSlice([]int{1, 5, 10, 15, 30, 60, 360, 720, 1440}),
			},
			"uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"locations": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
				MinItems: 1,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "ENABLED",
				ValidateFunc: validation.StringInSlice([]string{"ENABLED", "MUTED", "DISABLED"}, false),
			},
			"sla_threshold": {
				Type:     schema.TypeFloat,
				Optional: true,
			},
			//TODO: Add advanced options
			//"options": {
			//  "validationString": string [only valid for SIMPLE and BROWSER types],
			//  "verifySSL": boolean (true, false) [only valid for SIMPLE and BROWSER types],
			//  "bypassHEADRequest": boolean (true, false) [only valid for SIMPLE types],
			//  "treatRedirectAsFailure": boolean (true, false) [only valid for SIMPLE types]
			//  }

		},
	}
}

func resourceNewRelicSyntheticsMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	locationSet := d.Get("locations").([]interface{})
	locations := make([]string, len(locationSet))
	for i, location := range locationSet {
		locations[i] = location.(string)
	}

	createMonitor := synthetics.CreateMonitorArgs{
		Name:         d.Get("name").(string),
		Type:         d.Get("type").(string),
		Frequency:    uint(d.Get("frequency").(int)),
		URI:          d.Get("uri").(string),
		Locations:    locations,
		Status:       d.Get("status").(string),
		SLAThreshold: d.Get("sla_threshold").(float64),
	}

	log.Printf("[INFO] Creating New Relic synthetics monitor %s", createMonitor.Name)

	monitor, err := client.CreateMonitor(&createMonitor)
	if err != nil {
		return err
	}

	d.SetId(monitor.ID)

	return nil
}

func resourceNewRelicSyntheticsMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Reading New Relic synthetics monitor %s", d.Id())

	monitor, err := client.GetMonitor(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", monitor.Name)
	d.Set("type", monitor.Type)
	d.Set("frequency", monitor.Frequency)
	d.Set("uri", monitor.URI)
	d.Set("locations", monitor.Locations)
	d.Set("status", monitor.Status)
	d.Set("sla_threshold", monitor.SLAThreshold)

	return nil
}

func resourceNewRelicSyntheticsMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	locationSet := d.Get("locations").([]interface{})
	locations := make([]string, len(locationSet))
	for i, location := range locationSet {
		locations[i] = location.(string)
	}

	updatedMonitor := synthetics.UpdateMonitorArgs{
		Name:         d.Get("name").(string),
		Frequency:    uint(d.Get("frequency").(int)),
		URI:          d.Get("uri").(string),
		Locations:    locations,
		Status:       d.Get("status").(string),
		SLAThreshold: d.Get("sla_threshold").(float64),
	}

	log.Printf("[INFO] Updating New Relic synthetics monitor %s", d.Id())

	monitor, err := client.UpdateMonitor(d.Id(), &updatedMonitor)
	if err != nil {
		return err
	}

	d.Set("name", monitor.Name)
	d.Set("frequency", monitor.Frequency)
	d.Set("uri", monitor.URI)
	d.Set("locations", monitor.Locations)
	d.Set("status", monitor.Status)
	d.Set("sla_threshold", monitor.SLAThreshold)

	return nil
}

func resourceNewRelicSyntheticsMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Synthetics

	log.Printf("[INFO] Deleting New Relic synthetics monitor %s", d.Id())

	if err := client.DeleteMonitor(d.Id()); err != nil {
		return err
	}

	d.SetId("")
	return nil
}
