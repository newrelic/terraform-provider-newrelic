package newrelic

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
)

func resourceNewRelicCloudLinkAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudLinkAccountCreate,
		ReadContext:   resourceNewRelicCloudLinkAccountRead,
		UpdateContext: resourceNewRelicCloudLinkAccountUpdate,
		DeleteContext: resourceNewRelicCloudLinkAccountDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"aws": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AWS Provider",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"arn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The AWS role ARN",
						},
						"metricCollectionMode": {
							Type:         schema.TypeString,
							Optional:     true,
							Description:  "How metrics will be collected.",
							ValidateFunc: validation.StringInSlice([]string{"PULL", "PUSH"}, true),
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The linked account name",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicCloudLinkAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountId := selectAccountID(providerConfig, d)

	linkAccountInput, err := expandCloundLinkAccountInput(d)

	if err != nil {
		return diag.FromErr(err)
	}

	_, err := client.Cloud.CloudLinkAccountWithContext(ctx, accountId, input)
	if err != nil {
		return diag.FromErr(err)
	}
	
	return nil
}

func resourceNewRelicCloudLinkAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudLinkAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceNewRelicCloudLinkAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func expandCloundLinkAccountInput(d *schema.ResourceData) (*cloud.CloudLinkCloudAccountsInput, error) {
	input := cloud.CloudLinkCloudAccountsInput{
		Aws: expandAwsAccountsInput(d.Get("aws").([]interface{})),
	}
	return input, nil
}

func expandAwsAccountsInput(cfg []interface{}) ([]cloud.CloudAwsLinkAccountInput) {
	var awsAccounts []cloud.CloudAwsLinkAccountInput

	if len(cfg) == 0 {
		return awsAccounts
	}

	awsAccounts = make([]cloud.CloudAwsLinkAccountInput, 0, len(cfg))
	
	for _, a := range cfg {
		cfgAwsAccounts = a.(map[string]interface{})

		awsAccount := cloud.CloudAwsLinkAccountInput{}

		if arn, ok := cfgAwsAccounts["arn"]; ok {
			awsAccount.Arn = arn.(string)

			if m, ok := cfgAwsAccounts["metricCollectionMode"]: ok {
				awsAccount.MetricCollectionMode = m.(string)

				if n, ok := cfgAwsAccounts["name"]: ok {
					awsAccount.Name = n.(string)
				}
			}
		}
		awsAccounts = append(awsAccounts, awsAccount)
	}
	return nil
}
