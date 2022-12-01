package newrelic

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	nrErrors "github.com/newrelic/newrelic-client-go/v2/pkg/errors"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrqldroprules"
)

func resourceNewRelicNRQLDropRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicNRQLDropRuleCreate,
		ReadContext:   resourceNewRelicNRQLDropRuleRead,
		DeleteContext: resourceNewRelicNRQLDropRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				ValidateFunc: validation.StringInSlice([]string{"drop_data", "drop_attributes", "drop_attributes_from_metric_aggregates"}, false),
				Description:  "The drop rule action (drop_data, drop_attributes, or drop_attributes_from_metric_aggregates).",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceNewRelicNRQLDropRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := []nrqldroprules.NRQLDropRulesCreateDropRuleInput{
		{
			Description: d.Get("description").(string),
			Action:      nrqldroprules.NRQLDropRulesAction(strings.ToUpper(d.Get("action").(string))),
			NRQL:        d.Get("nrql").(string),
		},
	}

	created, err := client.Nrqldroprules.NRQLDropRulesCreateWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}
	//retry needed to check failure
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		created, err := client.Nrqldroprules.NRQLDropRulesCreateWithContext(ctx, accountID, createInput)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(created.Successes) == 0 {
			//assuming failures are returned a bit late
			log.Printf("The value of the failure is : %v ", created.Failures)
			return resource.RetryableError(fmt.Errorf("err: drop rule create result wasn't returned. Validate the action value or NRQL query."))
		}
		return nil
	})
	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	//Setting the errors
	//var apiDiags diag.Diagnostics

	if created == nil || len(created.Successes) == 0 {

		//for _, err := range created.Failures {
		//	apiDiags = append(apiDiags, diag.Diagnostic{
		//		Severity: diag.Error,
		//		Summary:  err.Error.Description,
		//		Detail:   string(err.Error.Reason),
		//	})
		//}
		//return apiDiags
		return diag.Errorf("err: drop rule create result wasn't returned. Validate the action value or NRQL query.")
	}
	rule := created.Successes[0]

	id := fmt.Sprintf("%d:%s", rule.AccountID, rule.ID)

	d.SetId(id)

	//put retry if  required

	return nil
	//return resourceNewRelicNRQLDropRuleRead(ctx, d, meta)
}

func resourceNewRelicNRQLDropRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic NRQL Drop Rule for %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err, apiErr := getNRQLDropRuleByID(ctx, client, accountID, ruleID)

	if err != nil || apiErr != nil {
		if _, ok := err.(*nrErrors.NotFound); ok || apiErr.Reason == "RULE_NOT_FOUND" {
			d.SetId("")
			return nil
		}

		var apiDiags diag.Diagnostics
		apiDiags = append(apiDiags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  apiErr.Description,
			Detail:   string(apiErr.Reason),
		})
		if apiDiags.HasError() {
			return apiDiags
		}
		return diag.FromErr(err)
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("rule_id", ruleID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("action", strings.ToLower(string(rule.Action))); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", rule.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("nrql", rule.NRQL); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicNRQLDropRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic entity tags from entity guid %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInput := []string{ruleID}

	deleted, err := client.Nrqldroprules.NRQLDropRulesDeleteWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	//Setting the errors
	var apiDiags diag.Diagnostics
	if len(deleted.Failures) != 0 {
		for _, err := range deleted.Failures {
			apiDiags = append(apiDiags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Error.Description,
				Detail:   string(err.Error.Reason),
			})
		}
		return apiDiags
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

//getNRQLDropRuleByID() returns the rule with the given ID.
func getNRQLDropRuleByID(ctx context.Context, client *newrelic.NewRelic, accountID int, ruleID string) (*nrqldroprules.NRQLDropRulesDropRule, error, *nrqldroprules.NRQLDropRulesError) {
	rules, err := client.Nrqldroprules.GetListWithContext(ctx, accountID)
	if err != nil {
		return nil, err, nil
	}

	if &rules.Error != nil {
		//set the values
		return nil, err, &rules.Error
	}
	for _, v := range rules.Rules {
		if v.ID == ruleID {
			return &v, nil, nil
		}
	}
	return nil, errors.New("drop rule not found"), nil
}
