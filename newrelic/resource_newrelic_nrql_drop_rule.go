package newrelic

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/newrelic"
	nrErrors "github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/nrqldroprules"
)

func resourceNewRelicNRQLDropRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicNRQLDropRuleCreate,
		Read:   resourceNewRelicNRQLDropRuleRead,
		Delete: resourceNewRelicNRQLDropRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Account with the NRQL drop rule will be put.",
			},
			"action": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"drop_data", "drop_attributes"}, false),
				Description:  "The drop rule action (drop_data or drop_attributes).",
			},
			"nrql": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Explains which data to apply the drop rule to.",
			},
			"description": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Provides additional information about the rule.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id, uniquely identifying the rule.",
			},
		},
	}
}

func resourceNewRelicNRQLDropRuleCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Create")
	}

	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := []nrqldroprules.NRQLDropRulesCreateDropRuleInput{
		{
			Description: d.Get("description").(string),
			Action:      nrqldroprules.NRQLDropRulesAction(strings.ToUpper(d.Get("action").(string))),
			NRQL:        d.Get("nrql").(string),
		},
	}

	created, err := client.Nrqldroprules.NRQLDropRulesCreate(accountID, createInput)
	if err != nil {
		return err
	}

	if created == nil {
		return errors.New("err: drop rule create result wasn't returned")
	}
	rule := created.Successes[0]

	id := fmt.Sprintf("%d:%s", rule.AccountID, rule.ID)

	d.SetId(id)

	return resourceNewRelicNRQLDropRuleRead(d, meta)
}

func resourceNewRelicNRQLDropRuleRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Read")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic NRQL Drop Rule for %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return err
	}

	rule, err := getNRQLDropRuleByID(client, accountID, ruleID)

	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	if err := d.Set("account_id", accountID); err != nil {
		return err
	}

	if err := d.Set("rule_id", ruleID); err != nil {
		return err
	}

	if err := d.Set("action", strings.ToLower(string(rule.Action))); err != nil {
		return err
	}

	if err := d.Set("description", rule.Description); err != nil {
		return err
	}

	if err := d.Set("nrql", rule.NRQL); err != nil {
		return err
	}

	return nil
}

func resourceNewRelicNRQLDropRuleDelete(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)

	if !providerConfig.hasNerdGraphCredentials() {
		return errors.New("err: NerdGraph support not present, but required for Delete")
	}

	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic entity tags from entity guid %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return err
	}

	deleteInput := []string{ruleID}

	_, err = client.Nrqldroprules.NRQLDropRulesDelete(accountID, deleteInput)
	if err != nil {
		return err
	}

	return nil
}

func parseNRQLDropRuleIDs(id string) (int, string, error) {
	strIDs := strings.Split(id, ":")

	if len(strIDs) != 2 {
		return 0, "", errors.New("could not parse drop rule IDs")
	}

	accountID, err := strconv.Atoi(strIDs[0])
	if err != nil {
		return 0, "", err
	}

	return accountID, strIDs[1], nil
}

func getNRQLDropRuleByID(client *newrelic.NewRelic, accountID int, ruleID string) (*nrqldroprules.NRQLDropRulesDropRule, error) {
	rules, err := client.Nrqldroprules.GetList(accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range rules.Rules {
		if v.ID == ruleID {
			return &v, nil
		}
	}
	return nil, errors.New("drop rule not found")
}
