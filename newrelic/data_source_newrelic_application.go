package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
)

func dataSourceNewRelicApplication() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "Use the `newrelic_entity` data source instead.",
		ReadContext:        dataSourceNewRelicApplicationRead,
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

func dataSourceNewRelicApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic applications")

	name := d.Get("name").(string)
	params := apm.ListApplicationsParams{
		Name: name,
	}

	applications, err := client.APM.ListApplications(&params)
	if err != nil {
		return diag.FromErr(err)
	}

	var application *apm.Application

	for _, a := range applications {
		if a.Name == name {
			application = a
			break
		}
	}

	if application == nil {
		return diag.FromErr(fmt.Errorf("the name '%s' does not match any New Relic applications", name))
	}

	return diag.FromErr(flattenApplicationData(application, d))
}

func flattenApplicationData(a *apm.Application, d *schema.ResourceData) error {
	d.SetId(strconv.Itoa(a.ID))
	var err error

	err = d.Set("name", a.Name)
	if err != nil {
		return err
	}

	err = d.Set("instance_ids", a.Links.InstanceIDs)
	if err != nil {
		return err
	}

	err = d.Set("host_ids", a.Links.HostIDs)
	if err != nil {
		return err
	}

	return nil
}
