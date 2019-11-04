package newrelic

import (
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func resourceNewRelicAlertPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicAlertPolicyCreate,
		Read:   resourceNewRelicAlertPolicyRead,
		Update: resourceNewRelicAlertPolicyUpdate,
		Delete: resourceNewRelicAlertPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"incident_preference": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PER_POLICY",
				ValidateFunc: validation.StringInSlice([]string{"PER_POLICY", "PER_CONDITION", "PER_CONDITION_AND_TARGET"}, false),
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

func buildAlertPolicyStruct(d *schema.ResourceData) *newrelic.AlertPolicy {
	policy := newrelic.AlertPolicy{
		Name: d.Get("name").(string),
	}

	if attr, ok := d.GetOk("incident_preference"); ok {
		policy.IncidentPreference = attr.(string)
	}

	return &policy
}

func resourceNewRelicAlertPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	policy := buildAlertPolicyStruct(d)

	log.Printf("[INFO] Creating New Relic alert policy %s", policy.Name)

	policy, err := client.CreateAlertPolicy(*policy)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(policy.ID))

	return nil
}

func unixMillis(msec int64) time.Time {
	sec := int64(msec / 1000)
	nsec := int64((msec - (sec * 1000)) * 1000000)
	// Note: this will default to local time
	created := time.Unix(sec, nsec)
	return created
}

func resourceNewRelicAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading New Relic alert policy %v", id)

	policy, err := client.GetAlertPolicy(int(id))
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	// New Relic provides created_at and updated_at as millisecond unix timestamps
	// https://www.terraform.io/docs/extend/schemas/schema-types.html#date-amp-time-data
	// "TypeString is also used for date/time data, the preferred format is RFC 3339."
	created := unixMillis(policy.CreatedAt).Format(time.RFC3339)
	updated := unixMillis(policy.UpdatedAt).Format(time.RFC3339)

	d.Set("name", policy.Name)
	d.Set("incident_preference", policy.IncidentPreference)
	d.Set("created_at", created)
	d.Set("updated_at", updated)

	return nil
}

func resourceNewRelicAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	policy := buildAlertPolicyStruct(d)

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	policy.ID = int(id)

	log.Printf("[INFO] Updating New Relic alert policy %d", id)
	respPolicy, err := client.UpdateAlertPolicy(*policy)
	if err != nil {
		return err
	}

	d.Set("created_at", respPolicy.CreatedAt)
	d.Set("updated_at", respPolicy.UpdatedAt)

	return nil
}

func resourceNewRelicAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic alert policy %v", id)

	if err := client.DeleteAlertPolicy(int(id)); err != nil {
		return err
	}

	return nil
}
