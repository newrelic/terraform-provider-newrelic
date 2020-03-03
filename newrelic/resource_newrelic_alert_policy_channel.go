package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
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
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"channel_ids"},
				Deprecated:    "use `channel_ids` argument instead",
			},
			"channel_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MinItems:      1,
				ConflictsWith: []string{"channel_id"},
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
	}
}

func resourceNewRelicAlertPolicyChannelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	policyChannels, err := expandAlertPolicyChannels(d)

	if err != nil {
		return err
	}

	sortIntegerSlice(policyChannels.ChannelIDs)

	serializedID := serializeIDs(append(
		[]int{policyChannels.ID},
		policyChannels.ChannelIDs...,
	))

	log.Printf("[INFO] Creating New Relic alert policy channel %s", serializedID)

	_, err = client.Alerts.UpdatePolicyChannels(
		policyChannels.ID,
		policyChannels.ChannelIDs,
	)

	if err != nil {
		return err
	}

	d.SetId(serializedID)

	return resourceNewRelicAlertPolicyChannelRead(d, meta)
}

func resourceNewRelicAlertPolicyChannelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	policyID := ids[0]
	parsedChannelIDs := ids[1:]

	sortIntegerSlice(parsedChannelIDs)

	log.Printf("[INFO] Reading New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelsExist(client, policyID, parsedChannelIDs)

	if err != nil {
		return err
	}

	if !exists {
		d.SetId("")
		return nil
	}

	return flattenAlertPolicyChannels(d, policyID, parsedChannelIDs)
}

func resourceNewRelicAlertPolicyChannelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	policyID := ids[0]
	channelIDs := ids[1:]

	log.Printf("[INFO] Deleting New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelsExist(client, policyID, channelIDs)
	if err != nil {
		return err
	}

	if exists {
		for _, id := range channelIDs {
			if _, err := client.Alerts.DeletePolicyChannel(policyID, id); err != nil {
				if _, ok := err.(*errors.NotFound); ok {
					return nil
				}
				return err
			}
		}
	}

	return nil
}

func policyChannelsExist(
	client *newrelic.NewRelic,
	policyID int,
	channelIDs []int,
) (bool, error) {
	channels, err := client.Alerts.ListChannels()

	if err != nil {
		return false, err
	}

	var foundChannels []*alerts.Channel
	for _, id := range channelIDs {
		ch := findChannel(channels, id)

		if ch == nil {
			return false, nil
		}

		foundChannels = append(foundChannels, ch)
	}

	var policyChannelsCount int
	for _, channel := range foundChannels {
		for _, id := range channel.Links.PolicyIDs {
			if id == policyID {
				policyChannelsCount++
			}
		}
	}

	// Ensure all channels exist on the policy
	if policyChannelsCount == len(channelIDs) {
		return true, nil
	}

	return false, nil
}

func findChannel(channels []*alerts.Channel, id int) *alerts.Channel {
	for _, channel := range channels {
		if channel.ID == id {
			return channel
		}
	}

	return nil
}
