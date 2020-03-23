package newrelic

import (
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
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
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The New Relic account ID to operate on.",
				DefaultFunc: envAccountID,
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The New Relic policy ID of the resource.",
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
	cfg := meta.(*ProviderConfig)

	if !cfg.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Create")
	}

	client := cfg.NewClient
	accountID := d.Get("account_id").(int)

	policy := alerts.AlertsPolicyInput{}

	if attr, ok := d.GetOk("incident_preference"); ok {
		if attr.(string) != "" {
			policy.IncidentPreference = alerts.AlertsIncidentPreference(attr.(string))
		}
	}

	if attr, ok := d.GetOk("name"); ok {
		policy.Name = attr.(string)
	}

	createResult, err := client.Alerts.CreatePolicyMutation(accountID, policy)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{createResult.ID, accountID}))
	err = flattenAlertPolicy(createResult, d, accountID)
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

		_, err = client.Alerts.UpdatePolicyChannels(createResult.ID, matchedChannelIDs)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceNewRelicAlertPolicyRead(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*ProviderConfig)

	if !cfg.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Read")
	}

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	accountID := ids[1]

	log.Printf("[INFO] Reading New Relic alert policy %d from account %d", policyID, accountID)

	client := cfg.NewClient

	queryPolicy, queryErr := client.Alerts.QueryPolicy(accountID, policyID)
	if queryErr != nil {
		return queryErr
	}

	return flattenAlertPolicy(queryPolicy, d, accountID)
}

func resourceNewRelicAlertPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*ProviderConfig)

	if !cfg.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Update")
	}

	client := cfg.NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	accountID := ids[1]

	log.Printf("[INFO] Updating New Relic alert policy %d from account %d", policyID, accountID)

	updatePolicy := alerts.AlertsPolicyUpdateInput{}

	if attr, ok := d.GetOk("incident_preference"); ok {
		if attr.(string) != "" {
			updatePolicy.IncidentPreference = alerts.AlertsIncidentPreference(attr.(string))
		}
	}

	if attr, ok := d.GetOk("name"); ok {
		updatePolicy.Name = attr.(string)
	}

	updateResult, updateErr := client.Alerts.UpdatePolicyMutation(accountID, policyID, updatePolicy)
	if updateErr != nil {
		return updateErr
	}

	return flattenAlertPolicy(updateResult, d, accountID)
}

func resourceNewRelicAlertPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*ProviderConfig)

	if !cfg.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Delete")
	}

	client := cfg.NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	accountID := ids[1]

	log.Printf("[INFO] Deleting New Relic alert policy %d from account %d", policyID, accountID)

	_, err = client.Alerts.DeletePolicyMutation(accountID, policyID)
	if err != nil {
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
