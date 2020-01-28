package newrelic

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

func expandAlertPolicyChannels(d *schema.ResourceData) (*alerts.PolicyChannels, error) {
	channelID := d.Get("channel_id").(int)
	channelIDs := d.Get("channel_ids").([]interface{})

	if channelID == 0 && len(channelIDs) == 0 {
		return nil, fmt.Errorf("must provide channel_id or channel_ids for resource newrelic_alert_policy_channel")
	}

	var ids []int

	if channelID != 0 {
		ids = []int{channelID}
	} else {
		ids = expandChannelIDs(channelIDs)
	}

	policyChannels := alerts.PolicyChannels{
		ID:         d.Get("policy_id").(int),
		ChannelIDs: ids,
	}

	return &policyChannels, nil
}

func expandChannelIDs(channelIDs []interface{}) []int {
	ids := make([]int, len(channelIDs))

	for i := range ids {
		ids[i] = channelIDs[i].(int)
	}

	return ids
}
