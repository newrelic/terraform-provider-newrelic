package newrelic

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	_ "github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/logconfigurations"
	"log"
)

//
func resourceNewRelicObfuscationExpression() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicObfuscationExpressionCreate,
		ReadContext:   resourceNewRelicObfuscationExpressionRead,
		UpdateContext: resourceNewRelicObfuscationExpressionUpdate,
		DeleteContext: resourceNewRelicObfuscationExpressionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Description: "The account id associated with the obfuscation expression.",
				Computed:    true,
				Optional:    true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of expression.",
				Optional:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Name of expression.",
				Required:    true,
			},
			"regex": {
				Type:        schema.TypeString,
				Description: "Regex of expression.",
				Required:    true,
			},
		},
	}
}

//Create the obfuscation expression
func resourceNewRelicObfuscationExpressionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	createInput := logconfigurations.LogConfigurationsCreateObfuscationExpressionInput{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Regex:       d.Get("regex").(string),
	}

	created, err := client.Logconfigurations.LogConfigurationsCreateObfuscationExpressionWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if created == nil {
		return diag.Errorf("err: obfuscation expression create result wasn't returned or expression was not created.")
	}

	expressionID := created.ID

	d.SetId(expressionID)

	return resourceNewRelicObfuscationExpressionRead(ctx, d, meta)
}

//Read the obfuscation expression
func resourceNewRelicObfuscationExpressionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	expressionID := d.Id()
	expression, err := getObfuscationExpressionByID(ctx, client, accountID, expressionID)

	if err != nil && expression == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("account_id", accountID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", expression.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("description", expression.Description); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("regex", expression.Regex); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

//Update the obfuscation expression
func resourceNewRelicObfuscationExpressionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandObfuscationExpressionUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One obfuscation expression %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)

	_, err := client.Logconfigurations.LogConfigurationsUpdateObfuscationExpressionWithContext(ctx, accountID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicObfuscationExpressionRead(ctx, d, meta)
}

func expandObfuscationExpressionUpdateInput(d *schema.ResourceData) logconfigurations.LogConfigurationsUpdateObfuscationExpressionInput {
	updateInp := logconfigurations.LogConfigurationsUpdateObfuscationExpressionInput{
		ID: d.Id(),
	}

	if e, ok := d.GetOk("name"); ok {
		updateInp.Name = e.(string)
	}

	if e, ok := d.GetOk("regex"); ok {
		updateInp.Regex = e.(string)
	}

	if e, ok := d.GetOk("description"); ok {
		updateInp.Description = e.(string)
	}

	return updateInp
}

//Delete the obfuscation expression
func resourceNewRelicObfuscationExpressionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic obfuscation expression id %s", d.Id())

	accountID := selectAccountID(meta.(*ProviderConfig), d)
	expressionID := d.Id()

	_, err := client.Logconfigurations.LogConfigurationsDeleteObfuscationExpressionWithContext(ctx, accountID, expressionID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getObfuscationExpressionByID(ctx context.Context, client *newrelic.NewRelic, accountID int, expressonID string) (*logconfigurations.LogConfigurationsObfuscationExpression, error) {
	expressions, err := client.Logconfigurations.GetObfuscationExpressionsWithContext(ctx, accountID)
	if err != nil {
		return nil, err
	}

	for _, v := range *expressions {
		log.Printf("values %s:%T", v.ID, (v.ID))
		if v.ID == expressonID {
			return &v, nil
		}
	}
	return nil, errors.New("obfuscation expression not found")

}
