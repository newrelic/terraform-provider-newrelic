package newrelic

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicAlertPolicyChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAlertPolicyChannelCreate,
		ReadContext:   resourceNewRelicAlertPolicyChannelRead,
		// Update: Not currently supported in API
		DeleteContext: resourceNewRelicAlertPolicyChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy.",
			},
			"channel_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Description: "Array of channel IDs to apply to the specified policy. We recommended sorting channel IDs in ascending order to avoid drift your Terraform state.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceNewRelicAlertPolicyChannelV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateStateNewRelicAlertPolicyChannelV0toV1,
				Version: 0,
			},
		},
	}
}

func resourceNewRelicAlertPolicyChannelV0() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAlertPolicyChannelCreate,
		ReadContext:   resourceNewRelicAlertPolicyChannelRead,
		// Update: Not currently supported in API
		DeleteContext: resourceNewRelicAlertPolicyChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy.",
			},
			"channel_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MinItems:    1,
				Description: "Array of channel IDs to apply to the specified policy. We recommended sorting channel IDs in ascending order to avoid drift your Terraform state.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
		},
		SchemaVersion: 0,
	}
}
func resourceNewRelicAlertPolicyChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	policyChannels, err := expandAlertPolicyChannels(d)

	if err != nil {
		return diag.FromErr(err)
	}

	sortIntegerSlice(policyChannels.ChannelIDs)

	serializedID := serializeIDs(append(
		[]int{policyChannels.ID},
		policyChannels.ChannelIDs...,
	))

	log.Printf("[INFO] Creating New Relic alert policy channel %s", serializedID)

	_, err = client.Alerts.UpdatePolicyChannelsWithContext(
		ctx,
		policyChannels.ID,
		policyChannels.ChannelIDs,
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(serializedID)

	return resourceNewRelicAlertPolicyChannelRead(ctx, d, meta)
}

func resourceNewRelicAlertPolicyChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	parsedChannelIDs := ids[1:]

	sortIntegerSlice(parsedChannelIDs)

	log.Printf("[INFO] Reading New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelsExist(ctx, client, policyID, parsedChannelIDs)

	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		d.SetId("")
		return nil
	}

	return diag.FromErr(flattenAlertPolicyChannels(d, policyID, parsedChannelIDs))
}

func resourceNewRelicAlertPolicyChannelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := ids[0]
	channelIDs := ids[1:]

	log.Printf("[INFO] Deleting New Relic alert policy channel %s", d.Id())

	exists, err := policyChannelsExist(ctx, client, policyID, channelIDs)
	if err != nil {
		return diag.FromErr(err)
	}

	if exists {
		for _, id := range channelIDs {
			if _, err := client.Alerts.DeletePolicyChannelWithContext(ctx, policyID, id); err != nil {
				if _, ok := err.(*errors.NotFound); ok {
					return nil
				}
				return diag.FromErr(err)
			}
		}
	}

	return nil
}

func policyChannelsExist(
	ctx context.Context,
	client *newrelic.NewRelic,
	policyID int,
	channelIDs []int,
) (bool, error) {
	channels, err := client.Alerts.ListChannelsWithContext(ctx)

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
