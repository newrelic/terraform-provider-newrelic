package newrelic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
)

func expandApplication(d *schema.ResourceData) *apm.Application {
	a := apm.Application{
		Name: d.Get("name").(string),
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		log.Printf("[ERROR] expanding application, id %s", err)
	} else {
		a.ID = id
	}

	a.Settings.AppApdexThreshold = d.Get("app_apdex_threshold").(float64)
	a.Settings.EndUserApdexThreshold = d.Get("end_user_apdex_threshold").(float64)
	a.Settings.EnableRealUserMonitoring = d.Get("enable_real_user_monitoring").(bool)

	return &a
}

func flattenApplication(a *apm.Application, d *schema.ResourceData) {
	d.SetId(strconv.Itoa(a.ID))
	d.Set("name", a.Name)
	d.Set("app_apdex_threshold", a.Settings.AppApdexThreshold)
	d.Set("end_user_apdex_threshold", a.Settings.EndUserApdexThreshold)
	d.Set("enable_real_user_monitoring", a.Settings.EnableRealUserMonitoring)
}
