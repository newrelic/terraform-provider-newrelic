package newrelic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(60 * time.Second),
		},
		DeprecationMessage: `The "newrelic_nrql_drop_rule" resource is deprecated and will be removed in a future major release after January 7, 2026. 
This aligns with the official end-of-life (EOL) for this feature by New Relic. 
Please migrate to its replacement, Pipeline Cloud Rules, using the "newrelic_pipeline_cloud_rule" resource: https://registry.terraform.io/providers/newrelic/newrelic/latest/docs/resources/pipeline_cloud_rule. 
For more details, see the EOL announcement: https://docs.newrelic.com/eol/2025/05/drop-rule-filter/`,

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
			"pipeline_cloud_rule_entity_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The GUID of the Pipeline Cloud Rule that will replace the Drop Rule after its EOL on January 7, 2026.",
			},
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

	var diags diag.Diagnostics
	if len(created.Failures) != 0 || created == nil {
		for _, resp := range created.Failures {
			switch resp.Error.Reason {
			case nrqldroprules.NRQLDropRulesErrorReasonTypes.GENERAL:
				return diag.Errorf("err: drop rule create result wasn't returned. Validate the action value or NRQL query.")
			default:
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  string(resp.Error.Reason),
					Detail:   resp.Error.Description,
				})
			}
		}
		return diags
	}

	rule := created.Successes[0]

	id := fmt.Sprintf("%d:%s", rule.AccountID, rule.ID)

	d.SetId(id)

	return resourceNewRelicNRQLDropRuleRead(ctx, d, meta)
}

func resourceNewRelicNRQLDropRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic NRQL Drop Rule for %s", d.Id())

	accountID, ruleID, err := parseNRQLDropRuleIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rule, err := getNRQLDropRuleByID(ctx, client, accountID, ruleID)
	if err != nil {
		if _, ok := err.(*nrErrors.NotFound); ok {
			d.SetId("")
			return nil
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

	// applying a retry mechanism on `pipeline_cloud_rule_entity_id`, to prevent it being from being lost
	// since we call a read immediately after creation, and the entity could need time to be fully indexed
	// the maximum timeout is set to 1 minute in the resource configuration above, and is not customizable by the customer
	// this is done in order to prevent an infinite timeout, as the default TimeoutRead is 20 minutes
	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		pipelineCloudRuleEntityID := rule.PipelineCloudRuleEntityId
		if pipelineCloudRuleEntityID == "" {
			return resource.RetryableError(fmt.Errorf("waiting for a non-null pipelineCloudRuleEntityId"))
		}

		if err := d.Set("pipeline_cloud_rule_entity_id", rule.PipelineCloudRuleEntityId); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
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

	var diags diag.Diagnostics
	if len(deleted.Failures) != 0 {
		for _, resp := range deleted.Failures {
			switch resp.Error.Reason {
			case nrqldroprules.NRQLDropRulesErrorReasonTypes.GENERAL:
				return diag.Errorf("err: drop rule create result wasn't returned. Validate the action value or NRQL query.")
			default:
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  string(resp.Error.Reason),
					Detail:   resp.Error.Description,
				})
			}
		}
		return diags
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

func getNRQLDropRuleByID(ctx context.Context, client *newrelic.NewRelic, accountID int, ruleID string) (*nrqldroprules.NRQLDropRulesDropRule, error) {
	rules, err := client.Nrqldroprules.GetListWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if rules.Error.Reason == nrqldroprules.NRQLDropRulesErrorReasonTypes.RULE_NOT_FOUND {
		return nil, nrErrors.NewNotFound("RULE_NOT_FOUND\n" + rules.Error.Description)
	}
	for _, v := range rules.Rules {
		if v.ID == ruleID {
			return &v, nil
		}
	}
	return nil, nrErrors.NewNotFound("RULE_NOT_FOUND")
}
