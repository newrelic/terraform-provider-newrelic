package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func dataSourceNewRelicAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicAlertPolicyRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"incident_preference": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic Alert Policies")

	policies, err := client.ListAlertPolicies()
	if err != nil {
		return err
	}

	var policy *newrelic.AlertPolicy
	name := d.Get("name").(string)

	for _, c := range policies {
		if strings.Contains(strings.ToLower(c.Name), strings.ToLower(name)) {
			policy = &c
			break
		}
	}

	if policy == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic alert policy", name)
	}

	d.SetId(strconv.Itoa(policy.ID))
	d.Set("name", policy.Name)
	d.Set("incident_preference", policy.IncidentPreference)
	d.Set("created_at", policy.CreatedAt)
	d.Set("updated_at", policy.UpdatedAt)

	return nil
}
