package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
)

func dataSourceNewRelicApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the application in New Relic.",
			},
			"instance_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: "A list of instance IDs associated with the application.",
			},
			"host_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: "A list of host IDs associated with the application.",
			},
		},
	}
}

func dataSourceNewRelicApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic applications")

	name := d.Get("name").(string)
	params := apm.ListApplicationsParams{
		Name: name,
	}

	applications, err := client.APM.ListApplications(&params)

	if err != nil {
		return err
	}

	var application *apm.Application

	for _, a := range applications {
		if a.Name == name {
			application = a
			break
		}
	}

	if application == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic applications", name)
	}

	flattenApplication(application, d)

	return nil
}

func flattenApplication(a *apm.Application, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(a.ID))
	d.Set("name", a.Name)
	d.Set("instance_ids", a.Links.InstanceIDs)
	d.Set("host_ids", a.Links.HostIDs)
}
