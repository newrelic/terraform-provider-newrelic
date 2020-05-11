package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
)

func resourceNewRelicApplicationSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicApplicationSettingsCreate,
		Read:   resourceNewRelicApplicationSettingsRead,
		Update: resourceNewRelicApplicationSettingsUpdate,
		Delete: resourceNewRelicApplicationSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"app_apdex_threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"end_user_apdex_threshold": {
				Type:     schema.TypeFloat,
				Required: true,
			},
			"enable_real_user_monitoring": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceNewRelicApplicationSettingsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	userApp := expandApplication(d)

	listParams := apm.ListApplicationsParams{
		Name: userApp.Name,
	}

	result, err := client.APM.ListApplications(&listParams)
	if err != nil {
		return err
	}

	if len(result) != 1 {
		return fmt.Errorf("more/less than one result from query for %s", userApp.Name)
	}

	app := *result[0]

	if app.Name != userApp.Name {
		return fmt.Errorf("the result name %s does not match requested name %s", app.Name, userApp.Name)
	}

	d.SetId(strconv.Itoa(app.ID))

	log.Printf("[INFO] Importing New Relic application %v", userApp.Name)
	return resourceNewRelicApplicationSettingsUpdate(d, meta)
}

func resourceNewRelicApplicationSettingsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	userApp := expandApplication(d)
	log.Printf("[INFO] Reading New Relic application %+v", userApp)

	app, err := client.APM.GetApplication(userApp.ID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Read found New Relic application %+v\n\n\n", app)

	return flattenApplication(app, d)
}

func resourceNewRelicApplicationSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	userApp := expandApplication(d)

	updateParams := apm.UpdateApplicationParams{
		Name:     userApp.Name,
		Settings: userApp.Settings,
	}

	log.Printf("[INFO] Updating New Relic application %+v with params: %+v", userApp, updateParams)

	_, err := client.APM.UpdateApplication(userApp.ID, updateParams)
	if err != nil {
		return err
	}

	return resourceNewRelicApplicationSettingsRead(d, meta)
}

func resourceNewRelicApplicationSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	// You can not delete application settings
	return nil
}
