package newrelic

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

func resourceNewRelicCloudAzureIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAzureIntegrationsCreate,
		ReadContext:   resourceNewRelicCloudAzureIntegrationsRead,
		UpdateContext: resourceNewRelicCloudAzureIntegrationsUpdate,
		DeleteContext: resourceNewRelicCloudAzureIntegrationsDelete,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the account in New Relic.",
			},
			"linked_account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The ID of the linked Azure account in New Relic",
			},

			// List of Integrations with Azure

			"azure_api_management": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure api management integration",
				Elem:        cloudAzureIntegrationAzureApiManagement,
				MaxItems:    1,
			},
			"azure_app_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app gateway integration",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_app_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app services",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure containers",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_cosmos_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure cosmoDB",
				Elem:        " ",
				MaxItems:    1,
			},

			"azure_data_factory": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure data factory",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_event_hub": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure event hub",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_express_route": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure express route",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_firewalls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure firewalls",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_front_door": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure front door",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_functions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure functions",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_key_vault": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure key vault",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_load_balancer": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure load balancer",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_logic_apps": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure logic apps",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_machine_learning": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure machine learning",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_maria_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Maria DB",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_mysql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure mysql",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_postgresql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure postgresql",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_power_bi_dedicated": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure powerBI dedicated",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_redis_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure redis cache",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_service_bus": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure service bus",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_service_fabric": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The azure services fabric",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_sql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_sql_managed": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql managed",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure storage",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_virtual_machine": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual machine",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_virtual_networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual networks",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_vms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure vms",
				Elem:        " ",
				MaxItems:    1,
			},
			"azure_vpn_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure vpn gateway",
				Elem:        " ",
				MaxItems:    1,
			},
		},
	}
}

func cloudAzureIntegrationSchemaBase() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metrics_polling_interval": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The data polling interval in seconds.",
		},
	}
}

func cloudAzureIntegrationAzureApiManagement() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

/////
func resourceNewRelicCloudAzureIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	cloudAzureIntegrationsInput, _ := cloudAzureIntegrationsInput(d)

	cloudAzureIntegrationsPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, cloudAzureIntegrationsInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudAzureIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudAzureIntegrationsPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	if len(cloudAzureIntegrationsPayload.Integrations) > 0 {
		d.SetId(strconv.Itoa(d.Get("linked_account_id").(int)))
	}

	return resourceNewRelicCloudAzureIntegrationsRead(ctx, d, meta)
}

///

func resourceNewRelicCloudAzureIntegrationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	linkedAccountID, convErr := strconv.Atoi(d.Id())

	if convErr != nil {
		return diag.FromErr(convErr)
	}

	linkedAccount, err := client.Cloud.GetLinkedAccountWithContext(ctx, accountID, linkedAccountID)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenCloudAwsLinkedAccount(d, linkedAccount)

	return nil
}
}
func resourceNewRelicCloudAzureIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

}
func resourceNewRelicCloudAzureIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

}
