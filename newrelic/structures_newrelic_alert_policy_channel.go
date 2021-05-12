package newrelic

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
)

// migrateStateNewRelicAlertPolicyChannelV0toV1 currently facilitates migrating
// the `channel_ids` attribute from TypeList to TypeSet. Since the underlying
// data structure is []int for both, we don't need to do anything other than
// return the state and Terraform will handle the rest.
func migrateStateNewRelicAlertPolicyChannelV0toV1(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	return rawState, nil
}

func expandAlertPolicyChannels(d *schema.ResourceData) (*alerts.PolicyChannels, error) {
	channelIDs := d.Get("channel_ids").(*schema.Set)

	if channelIDs.Len() == 0 {
		return nil, fmt.Errorf("must provide channel_ids for resource newrelic_alert_policy_channel")
	}

	ids := expandChannelIDs(channelIDs.List())

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

func flattenAlertPolicyChannels(d *schema.ResourceData, policyID int, channelIDs []int) error {
	d.Set("policy_id", policyID)

	_, channelIDsOk := d.GetOk("channel_ids")

	if channelIDsOk && len(channelIDs) > 0 {
		d.Set("channel_ids", channelIDs)
	}

	// Handle import (set `channel_ids` since this resource doesn't exist in state yet)
	if !channelIDsOk {
		d.Set("channel_ids", channelIDs)
	}

	return nil
}
