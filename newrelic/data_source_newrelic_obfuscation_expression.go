package newrelic

import (
	"context"
	"errors"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
)

func dataSourceNewRelicObfuscationExpression() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNewRelicObfuscationExpressionRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The account id associated with the obfuscation expression.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of expression.",
			},
		},
	}
}

func dataSourceNewRelicObfuscationExpressionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading obfuscation expression.")

	if name, ok := d.GetOk("name"); ok {
		expression, err := getObfuscationExpressionByName(ctx, client, accountID, name.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		d.SetId(expression.ID)
		return nil
	}
	return diag.FromErr(errors.New("`name` is required"))

}

func getObfuscationExpressionByName(ctx context.Context, client *newrelic.NewRelic, accountID int, name string) (*logconfigurations.LogConfigurationsObfuscationExpression, error) {
	expressions, err := client.Logconfigurations.GetObfuscationExpressionsWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range *expressions {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, errors.New("obfuscation expression not found")
}
