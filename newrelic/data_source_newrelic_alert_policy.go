package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		if strings.EqualFold(c.Name, name) {
			policy = &c
			break
		}
	}

	if policy == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic alert policy", name)
	}

	// New Relic provides created_at and updated_at as millisecond unix timestamps
	// https://www.terraform.io/docs/extend/schemas/schema-types.html#date-amp-time-data
	// "TypeString is also used for date/time data, the preferred format is RFC 3339."
	created := unixMillis(policy.CreatedAt).Format(time.RFC3339)
	updated := unixMillis(policy.UpdatedAt).Format(time.RFC3339)

	d.SetId(strconv.Itoa(policy.ID))
	d.Set("name", policy.Name)
	d.Set("incident_preference", policy.IncidentPreference)
	d.Set("created_at", created)
	d.Set("updated_at", updated)

	return nil
}
