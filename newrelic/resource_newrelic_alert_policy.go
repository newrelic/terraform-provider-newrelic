package newrelic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The name of the policy.",
			},
			"incident_preference": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PER_POLICY",
				ValidateFunc: validation.StringInSlice([]string{"PER_POLICY", "PER_CONDITION", "PER_CONDITION_AND_TARGET"}, false),
				Description:  "The rollup strategy for the policy. Options include: PER_POLICY, PER_CONDITION, or PER_CONDITION_AND_TARGET. The default is PER_POLICY.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the policy was created.",
				Deprecated:  "Unavailable attribute in NerdGraph.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the policy was last updated.",
				Deprecated:  "Unavailable attribute in NerdGraph.",
			},
			"channel_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "An array of channel IDs (integers) to assign to the policy. Adding or removing channel IDs from this array will result in a new alert policy resource being created and the old one being destroyed. Also note that channel IDs cannot be imported via terraform import.",
			},
		},
	}
}

func resourceNewRelicAlertPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	p := expandAlertPolicy(d)

	log.Printf("[INFO] Creating New Relic alert policy %s", p.Name)

	policy, err := client.Alerts.CreatePolicy(*p)

	if err != nil {
		return err
	}

	channels := d.Get("channel_ids").([]interface{})

	if len(channels) > 0 {
		channelIDs := expandAlertChannelIDs(channels)
		matchedChannelIDs, err := findExistingChannelIDs(client, channelIDs)

		if err != nil {
			return err
		}

		log.Printf("[INFO] Adding channels %+v to policy %+v", matchedChannelIDs, policy.Name)

		_, err = client.Alerts.UpdatePolicyChannels(policy.ID, matchedChannelIDs)

		if err != nil {
			return err
		}
	}

	d.SetId(strconv.Itoa(policy.ID))

	return nil
}

func resourceNewRelicAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading New Relic alert policy %v", id)

	policy, err := client.Alerts.GetPolicy(int(id))

	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenAlertPolicy(policy, d)
}

func resourceNewRelicAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	policy := expandAlertPolicy(d)

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}
	policy.ID = int(id)

	log.Printf("[INFO] Updating New Relic alert policy %d", id)
	respPolicy, err := client.Alerts.UpdatePolicy(*policy)
	if err != nil {
		return err
	}

	d.Set("created_at", respPolicy.CreatedAt)
	d.Set("updated_at", respPolicy.UpdatedAt)

	return nil
}

func resourceNewRelicAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic alert policy %v", id)

	if _, err := client.Alerts.DeletePolicy(int(id)); err != nil {
		return err
	}

	return nil
}

func findExistingChannelIDs(client *newrelic.NewRelic, channelIDs []int) ([]int, error) {
	channels, err := client.Alerts.ListChannels()

	if err != nil {
		return nil, err
	}

	matched := make([]int, 0)

	for i := range channels {
		for n := range channelIDs {
			if channelIDs[n] == channels[i].ID {
				matched = append(matched, channelIDs[n])
			}
		}
	}

	return matched, nil
}
