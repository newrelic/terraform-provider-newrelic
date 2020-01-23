package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicAlertPolicyChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicAlertPolicyChannelCreate,
		Read:   resourceNewRelicAlertPolicyChannelRead,
		// Update: Not currently supported in API
		Delete: resourceNewRelicAlertPolicyChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"channel_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNewRelicAlertPolicyChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	policyID := d.Get("policy_id").(int)
	channelID := d.Get("channel_id").(int)

	serializedID := serializeIDs([]int{policyID, channelID})

	log.Printf("[INFO] Creating New Relic alert policy channel %s", serializedID)

	exists, err := policyChannelExists(client, policyID, channelID)
	if err != nil {
		return err
	}

	if !exists {
		_, err = client.Alerts.UpdatePolicyChannels(policyID, []int{channelID})
		if err != nil {
			return err
		}
	}

	d.SetId(serializedID)

	return nil
}

func resourceNewRelicAlertPolicyChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	channelID := ids[1]

	log.Printf("[INFO] Reading New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelExists(client, policyID, channelID)
	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
		return nil
	}

	d.Set("policy_id", policyID)
	d.Set("channel_id", channelID)

	return nil
}

func resourceNewRelicAlertPolicyChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	channelID := ids[1]

	log.Printf("[INFO] Deleting New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelExists(client, policyID, channelID)
	if err != nil {
		return err
	}

	if exists {
		if _, err := client.Alerts.DeletePolicyChannel(policyID, channelID); err != nil {
			if _, ok := err.(*errors.NotFound); ok {
				return nil
			}
			return err
		}
	}

	return nil
}

func policyChannelExists(client *newrelic.NewRelic, policyID int, channelID int) (bool, error) {
	channel, err := client.Alerts.GetChannel(channelID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			return false, nil
		}

		return false, err
	}

	for _, id := range channel.Links.PolicyIDs {
		if id == policyID {
			return true, nil
		}
	}

	return false, nil
}
