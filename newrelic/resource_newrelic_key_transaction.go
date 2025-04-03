package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/keytransaction"
)

func resourceNewRelicKeyTransaction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicKeyTransactionCreate,
		ReadContext:   resourceNewRelicKeyTransactionRead,
		UpdateContext: resourceNewRelicKeyTransactionUpdate,
		DeleteContext: resourceNewRelicKeyTransactionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"apdex_index": {
				Type:        schema.TypeFloat,
				Description: "The acceptable amount of the time spent in the backend before customers get frustrated (Apdex target)",
				Required:    true,
			},
			"application_guid": {
				Type:        schema.TypeString,
				Description: "The GUID of the application.",
				Required:    true,
				ForceNew:    true,
			},
			"browser_apdex_target": {
				Type:        schema.TypeFloat,
				Description: "The acceptable amount of time for rendering a page in a browser before customers get frustrated (browser Apdex target).",
				Required:    true,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Description: "The name of the metric underlying this key transaction",
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The name of the key transaction.",
				Required:    true,
			},
			"domain": {
				Type:        schema.TypeString,
				Description: "Domain of the entity.",
				Computed:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the entity.",
				Computed:    true,
			},
		},
	}
}

func resourceNewRelicKeyTransactionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	apdexIndex,
		applicationGUID,
		browserApdexTarget,
		metricName,
		name := resourceNewRelicKeyTransactionFetchValuesFromConfig(d)

	log.Printf("[INFO] Creating New Relic Key Transaction %s", name)
	createKeyTransactionResult, err := client.KeyTransaction.KeyTransactionCreate(
		apdexIndex,
		keytransaction.EntityGUID(applicationGUID),
		browserApdexTarget,
		metricName,
		name,
	)

	if err != nil {
		return diag.FromErr(err)
	}
	if createKeyTransactionResult == nil {
		return diag.FromErr(fmt.Errorf("something went wrong while creating the key transaction"))
	}

	resourceNewRelicKeyTransactionSetValuesToState(d, *createKeyTransactionResult, keytransaction.KeyTransactionUpdateResult{})

	return nil
}

func resourceNewRelicKeyTransactionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	guid := d.Id()

	resp, err := client.Entities.GetEntityWithContext(ctx, common.EntityGUID(guid))
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("no key transaction found with given guid %s", guid))
	}

	switch (*resp).(type) {
	case *entities.KeyTransactionEntity:
		entity := (*resp).(*entities.KeyTransactionEntity)
		_ = d.Set("apdex_index", entity.ApdexTarget)
		_ = d.Set("application_guid", entity.Application.GUID)
		_ = d.Set("browser_apdex_target", entity.BrowserApdexTarget)
		_ = d.Set("metric_name", entity.MetricName)
		_ = d.Set("name", entity.Name)
	default:
		return diag.FromErr(fmt.Errorf("entity with GUID %s was not a key transaction", guid))
	}
	return nil
}

func resourceNewRelicKeyTransactionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	apdexIndex,
		_,
		browserApdexTarget,
		_,
		name := resourceNewRelicKeyTransactionFetchValuesFromConfig(d)

	guid := keytransaction.EntityGUID(d.Id())
	log.Printf("[INFO] Updating New Relic Key Transaction %s", name)

	updateKeyTransactionResult, err := client.KeyTransaction.KeyTransactionUpdate(apdexIndex, browserApdexTarget, guid, name)

	if err != nil {

		return diag.FromErr(err)
	}
	if updateKeyTransactionResult == nil {
		return diag.FromErr(fmt.Errorf("something went wrong while updating the key transaction"))
	}

	resourceNewRelicKeyTransactionSetValuesToState(d, keytransaction.KeyTransactionCreateResult{}, *updateKeyTransactionResult)
	return nil
}

func resourceNewRelicKeyTransactionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	keyTransactionGUID := keytransaction.EntityGUID(d.Id())

	log.Printf("[INFO] Deleting New Relic Key Transaction %s", d.Id())
	_, err := client.KeyTransaction.KeyTransactionDelete(keyTransactionGUID)

	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNewRelicKeyTransactionFetchValuesFromConfig(d *schema.ResourceData) (
	apdexIndex float64,
	applicationGUID string,
	browserApdexTarget float64,
	metricName string,
	name string,
) {
	apdexIndex = d.Get("apdex_index").(float64)
	applicationGUID = d.Get("application_guid").(string)
	browserApdexTarget = d.Get("browser_apdex_target").(float64)
	metricName = d.Get("metric_name").(string)
	name = d.Get("name").(string)

	return apdexIndex, applicationGUID, browserApdexTarget, metricName, name
}

func resourceNewRelicKeyTransactionSetValuesToState(
	d *schema.ResourceData,
	createKeyTransactionResult keytransaction.KeyTransactionCreateResult,
	updateKeyTransactionResult keytransaction.KeyTransactionUpdateResult,
) {
	if createKeyTransactionResult.GUID != "" && updateKeyTransactionResult.Name == "" {
		d.SetId(string(createKeyTransactionResult.GUID))
		_ = d.Set("apdex_index", createKeyTransactionResult.ApdexTarget)
		_ = d.Set("application_guid", createKeyTransactionResult.Application.GUID)
		_ = d.Set("browser_apdex_target", createKeyTransactionResult.BrowserApdexTarget)
		_ = d.Set("metric_name", createKeyTransactionResult.MetricName)
		_ = d.Set("name", createKeyTransactionResult.Name)
		entity := createKeyTransactionResult.Application.Entity.(*keytransaction.ApmApplicationEntityOutline)
		_ = d.Set("domain", entity.Domain)
		_ = d.Set("type", entity.Type)
	} else if createKeyTransactionResult.GUID == "" && updateKeyTransactionResult.Name != "" {
		_ = d.Set("apdex_index", updateKeyTransactionResult.ApdexTarget)
		_ = d.Set("application_guid", updateKeyTransactionResult.Application.GUID)
		_ = d.Set("browser_apdex_target", updateKeyTransactionResult.BrowserApdexTarget)
		_ = d.Set("name", updateKeyTransactionResult.Name)
		entity := updateKeyTransactionResult.Application.Entity.(*keytransaction.ApmApplicationEntityOutline)
		_ = d.Set("domain", entity.Domain)
		_ = d.Set("type", entity.Type)
	}
}
