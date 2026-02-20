package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/alerts"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
)

func resourceNewRelicAlertPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicAlertPolicyCreate,
		ReadContext:   resourceNewRelicAlertPolicyRead,
		UpdateContext: resourceNewRelicAlertPolicyUpdate,
		DeleteContext: resourceNewRelicAlertPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateWithMetadata(1, "account_id"),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The name of the policy.",
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The New Relic account ID to operate on.",
			},
			"incident_preference": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(alerts.AlertsIncidentPreferenceTypes.PER_POLICY),
				ValidateFunc: validation.StringInSlice([]string{
					string(alerts.AlertsIncidentPreferenceTypes.PER_POLICY),
					string(alerts.AlertsIncidentPreferenceTypes.PER_CONDITION),
					string(alerts.AlertsIncidentPreferenceTypes.PER_CONDITION_AND_TARGET),
				},
					false,
				),
				Description: "The rollup strategy for the policy. Options include: PER_POLICY, PER_CONDITION, or PER_CONDITION_AND_TARGET. The default is PER_POLICY.",
			},
			"channel_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "An array of channel IDs (integers) to assign to the policy. Adding or removing channel IDs from this array will result in a new alert policy resource being created and the old one being destroyed. Also note that channel IDs cannot be imported via terraform import.",
				Deprecated:  "The `channel_ids` attribute is deprecated and will be removed in the next major release of the provider.",
			},
			"entity_guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the alert policy in New Relic.",
			},
		},
	}
}

func resourceNewRelicAlertPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	// We're creating this updatedContext for the findExistingChannelIDs function which is still REST based
	// It does not have an effect on the GraphQL calls
	updatedContext := updateContextWithAccountID(ctx, accountID)

	policy := alerts.AlertsPolicyInput{}

	if attr, ok := d.GetOk("incident_preference"); ok {
		if attr.(string) != "" {
			policy.IncidentPreference = alerts.AlertsIncidentPreference(attr.(string))
		}
	}

	if attr, ok := d.GetOk("name"); ok {
		policy.Name = attr.(string)
	}

	createResult, err := client.Alerts.CreatePolicyMutationWithContext(updatedContext, accountID, policy)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createResult.ID)
	if err = flattenAlertPolicyWithEntityGUID(ctx, client, createResult, d, accountID); err != nil {
		return diag.FromErr(err)
	}

	channels := d.Get("channel_ids").([]interface{})

	if len(channels) > 0 {
		channelIDs := expandAlertChannelIDs(channels)
		matchedChannelIDs, err := findExistingChannelIDs(updatedContext, client, channelIDs)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] Adding channels %+v to policy %+v", matchedChannelIDs, policy.Name)

		createResultID, err := strconv.Atoi(createResult.ID)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.Alerts.UpdatePolicyChannelsWithContext(updatedContext, createResultID, matchedChannelIDs)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceNewRelicAlertPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var accountID int
	var policyID int

	if len(ids) == 1 {
		policyID = ids[0]
		accountID = selectAccountID(providerConfig, d)
	} else if len(ids) == 2 {
		policyID = ids[0]
		accountID = ids[1]
	} else {
		err := fmt.Errorf("unhandled id format %s", d.Id())
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading New Relic alert policy %d from account %d", policyID, accountID)

	id := strconv.Itoa(policyID)

	queryPolicy, queryErr := client.Alerts.QueryPolicyWithContext(ctx, accountID, id)

	if queryErr != nil {
		if _, ok := queryErr.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(queryErr)
	}

	return diag.FromErr(flattenAlertPolicyWithEntityGUID(ctx, client, queryPolicy, d, accountID))
}

func resourceNewRelicAlertPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Updating New Relic alert policy %s from account %d", d.Id(), accountID)

	updatePolicy := alerts.AlertsPolicyUpdateInput{}

	if attr, ok := d.GetOk("incident_preference"); ok {
		if attr.(string) != "" {
			updatePolicy.IncidentPreference = alerts.AlertsIncidentPreference(attr.(string))
		}
	}

	if attr, ok := d.GetOk("name"); ok {
		updatePolicy.Name = attr.(string)
	}

	updateResult, updateErr := client.Alerts.UpdatePolicyMutationWithContext(ctx, accountID, d.Id(), updatePolicy)
	if updateErr != nil {
		return diag.FromErr(updateErr)
	}

	return diag.FromErr(flattenAlertPolicyWithEntityGUID(ctx, client, updateResult, d, accountID))
}

func resourceNewRelicAlertPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Deleting New Relic alert policy %s from account %d", d.Id(), accountID)

	_, err := client.Alerts.DeletePolicyMutationWithContext(ctx, accountID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func findExistingChannelIDs(ctx context.Context, client *newrelic.NewRelic, channelIDs []int) ([]int, error) {
	channels, err := client.Alerts.ListChannelsWithContext(ctx)

	if err != nil {
		return nil, err
	}

	matched := make([]int, 0)

	for n := range channelIDs {
		found := false
		for i := range channels {
			if channelIDs[n] == channels[i].ID {
				found = true
				matched = append(matched, channelIDs[n])
			}
		}
		if !found {
			return nil, fmt.Errorf("[ERROR] Channel ID %d not found", channelIDs[n])
		}
	}

	return matched, nil
}
