package newrelic

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/cloud"
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
				Elem:        cloudAzureIntegrationAzureAPIManagementElem(),
				MaxItems:    1,
			},
			"azure_app_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app gateway integration",
				Elem:        cloudAzureIntegrationAzureAppGatewayElem(),
				MaxItems:    1,
			},
			"azure_app_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app services",
				Elem:        cloudAzureIntegrationAzureAppServiceElem(),
				MaxItems:    1,
			},
			"azure_containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure containers",
				Elem:        cloudAzureIntegrationAzureContainersElem(),
				MaxItems:    1,
			},
			"azure_cosmos_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure cosmoDB",
				Elem:        cloudAzureIntegrationAzureCosmosDBElem(),
				MaxItems:    1,
			},
			"azure_cost_management": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure cost management",
				Elem:        cloudAzureIntegrationCostManagementElem(),
				MaxItems:    1,
			},
			"azure_data_factory": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure data factory",
				Elem:        cloudAzureIntegrationAzureDataFactoryElem(),
				MaxItems:    1,
			},
			"azure_event_hub": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure event hub",
				Elem:        cloudAzureIntegrationAzureEventHubElem(),
				MaxItems:    1,
			},
			"azure_express_route": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure express route",
				Elem:        cloudAzureIntegrationAzureExpressRouteElem(),
				MaxItems:    1,
			},
			"azure_firewalls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure firewalls",
				Elem:        cloudAzureIntegrationAzureFirewallsElem(),
				MaxItems:    1,
			},
			"azure_front_door": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure front door",
				Elem:        cloudAzureIntegrationAzureFrontDoorElem(),
				MaxItems:    1,
			},
			"azure_functions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure functions",
				Elem:        cloudAzureIntegrationAzureFunctionsElem(),
				MaxItems:    1,
			},
			"azure_key_vault": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure key vault",
				Elem:        cloudAzureIntegrationAzureKeyVaultElem(),
				MaxItems:    1,
			},
			"azure_load_balancer": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure load balancer",
				Elem:        cloudAzureIntegrationAzureLoadBalancerElem(),
				MaxItems:    1,
			},
			"azure_logic_apps": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure logic apps",
				Elem:        cloudAzureIntegrationAzureLogicAppsElem(),
				MaxItems:    1,
			},
			"azure_machine_learning": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure machine learning",
				Elem:        cloudAzureIntegrationAzureMachineLearningElem(),
				MaxItems:    1,
			},
			"azure_maria_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Maria DB",
				Elem:        cloudAzureIntegrationAzureMariadbElem(),
				MaxItems:    1,
			},
			"azure_mysql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure mysql",
				Elem:        cloudAzureIntegrationAzureMysqlElem(),
				MaxItems:    1,
			},
			"azure_postgresql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure postgresql",
				Elem:        cloudAzureIntegrationAzurePostgresqlElem(),
				MaxItems:    1,
			},
			"azure_power_bi_dedicated": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure powerBI dedicated",
				Elem:        cloudAzureIntegrationAzurePowerBiDedicatedElem(),
				MaxItems:    1,
			},
			"azure_redis_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure redis cache",
				Elem:        cloudAzureIntegrationAzureRedisCacheElem(),
				MaxItems:    1,
			},
			"azure_service_bus": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure service bus",
				Elem:        cloudAzureIntegrationAzureServiceBusElem(),
				MaxItems:    1,
			},
			"azure_service_fabric": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The azure services fabric",
				Elem:        cloudAzureIntegrationAzureServiceFabricElem(),
				MaxItems:    1,
			},
			"azure_sql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql",
				Elem:        cloudAzureIntegrationAzureSQLElem(),
				MaxItems:    1,
			},
			"azure_sql_managed": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql managed",
				Elem:        cloudAzureIntegrationAzureSQLManagedElem(),
				MaxItems:    1,
			},
			"azure_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure storage",
				Elem:        cloudAzureIntegrationAzureStorageElem(),
				MaxItems:    1,
			},
			"azure_virtual_machine": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual machine",
				Elem:        cloudAzureIntegrationAzureVirtualMachineElem(),
				MaxItems:    1,
			},
			"azure_virtual_networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual networks",
				Elem:        cloudAzureIntegrationAzureVirtualNetworksElem(),
				MaxItems:    1,
			},
			"azure_vms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Vms",
				Elem:        cloudAzureIntegrationAzureVmsElem(),
				MaxItems:    1,
			},
			"azure_vpn_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure vpn gateway",
				Elem:        cloudAzureIntegrationAzureVPNGatewayElem(),
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
			Description: "The data polling interval in seconds",
		},
		"resource_groups": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		},
	}
}

func cloudAzureIntegrationAzureAPIManagementElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

func cloudAzureIntegrationAzureAppGatewayElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}

func cloudAzureIntegrationAzureAppServiceElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureContainersElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureCosmosDBElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationCostManagementElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["tag_keys"] = &schema.Schema{
		Type:        schema.TypeString,
		Description: "",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

func cloudAzureIntegrationAzureDataFactoryElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureEventHubElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureExpressRouteElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureFirewallsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureFrontDoorElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureFunctionsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureKeyVaultElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureLoadBalancerElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureLogicAppsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}
}
func cloudAzureIntegrationAzureMachineLearningElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureMariadbElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureMysqlElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzurePostgresqlElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzurePowerBiDedicatedElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureRedisCacheElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureServiceBusElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureServiceFabricElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureSQLElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureSQLManagedElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureStorageElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureVirtualMachineElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureVirtualNetworksElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureVmsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureVPNGatewayElem() *schema.Resource {
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

	cloudAzureIntegrationsInput, _ := expandCloudAzureIntegrationsInput(d)

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

	return nil
}

///// Inputs
//nolint: gocyclo
func expandCloudAzureIntegrationsInput(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	cloudAzureIntegration := cloud.CloudAzureIntegrationsInput{}
	cloudDisableAzureIntegration := cloud.CloudAzureDisableIntegrationsInput{}

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if v, ok := d.GetOk("azure_api_management"); ok {
		cloudAzureIntegration.AzureAPImanagement = expandCloudAzureIntegrationAzureAPIManagementInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_api_management"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAPImanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_app_gateway"); ok {
		cloudAzureIntegration.AzureAppgateway = expandCloudAzureIntegrationAzureAppGatewayInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_app_gateway"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAppgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("azure_app_service"); ok {
		cloudAzureIntegration.AzureAppservice = expandCloudAzureIntegrationAzureAppServiceInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_app_service"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAppservice = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_containers"); ok {
		cloudAzureIntegration.AzureContainers = expandCloudAzureIntegrationAzureContainersInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_containers"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureContainers = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_cosmos_db"); ok {
		cloudAzureIntegration.AzureCosmosdb = expandCloudAzureIntegrationAzureCosmosdbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_cosmos_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_cost_management"); ok {
		cloudAzureIntegration.AzureCostmanagement = expandCloudAzureIntegrationAzureCostManagementInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_cost_management"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("azure_data_factory"); ok {
		cloudAzureIntegration.AzureDatafactory = expandCloudAzureIntegrationAzureDataFactoryInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_data_factory"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureDatafactory = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_event_hub"); ok {
		cloudAzureIntegration.AzureEventhub = expandCloudAzureIntegrationCloudAzureEventHubInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_event_hub"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureEventhub = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_express_route"); ok {
		cloudAzureIntegration.AzureExpressroute = expandCloudAzureIntegrationAzureExpressRouteInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_express_route"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureExpressroute = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_firewalls"); ok {
		cloudAzureIntegration.AzureFirewalls = expandCloudAzureIntegrationAzureFirewallsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_firewalls"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_front_door"); ok {
		cloudAzureIntegration.AzureFrontdoor = expandCloudAzureIntegrationAzureFrontDoorInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_front_door"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureFrontdoor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_functions"); ok {
		cloudAzureIntegration.AzureFunctions = expandCloudAzureIntegrationAzureFunctionsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_functions"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureFunctions = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_key_vault"); ok {
		cloudAzureIntegration.AzureKeyvault = expandCloudAzureIntegrationAzureKeyVaultInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_key_vault"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureKeyvault = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_load_balancer"); ok {
		cloudAzureIntegration.AzureLoadbalancer = expandCloudAzureIntegrationAzureLoadBalancerInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_load_balancer"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureLoadbalancer = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_logic_apps"); ok {
		cloudAzureIntegration.AzureLogicapps = expandCloudAzureIntegrationAzureLogicAppsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_logic_apps"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureLogicapps = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	//
	if v, ok := d.GetOk("azure_machine_learning"); ok {
		cloudAzureIntegration.AzureMachinelearning = expandCloudAzureIntegrationAzureMachineLearningInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_machine_learning"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMachinelearning = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_maria_db"); ok {
		cloudAzureIntegration.AzureMariadb = expandCloudAzureIntegrationAzureMariadbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_maria_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMariadb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_mysql"); ok {
		cloudAzureIntegration.AzureMysql = expandCloudAzureIntegrationAzureMysqlInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_mysql"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMysql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_postgresql"); ok {
		cloudAzureIntegration.AzurePostgresql = expandCloudAzureIntegrationAzurePostgresqlInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_postgresql"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzurePostgresql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_power_bi_dedicated"); ok {
		cloudAzureIntegration.AzurePowerbidedicated = expandCloudAzureIntegrationAzurePowerBiDedicatedInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_power_bi_dedicated"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzurePowerbidedicated = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_redis_cache"); ok {
		cloudAzureIntegration.AzureRediscache = expandCloudAzureIntegrationAzureRedisCacheInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_redis_cache"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureRediscache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_service_bus"); ok {
		cloudAzureIntegration.AzureServicebus = expandCloudAzureIntegrationAzureServiceBusInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_service_bus"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureServicebus = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_service_fabric"); ok {
		cloudAzureIntegration.AzureServicefabric = expandCloudAzureIntegrationAzureServiceFabricInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_service_fabric"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureServicefabric = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_sql"); ok {
		cloudAzureIntegration.AzureSql = expandCloudAzureIntegrationAzureSQLInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_sql"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureSql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_sql_managed"); ok {
		cloudAzureIntegration.AzureSqlmanaged = expandCloudAzureIntegrationAzureSQLManagedInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_sql_managed"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureSqlmanaged = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_storage"); ok {
		cloudAzureIntegration.AzureStorage = expandCloudAzureIntegrationAzureStorageInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_storage"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureStorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_virtual_machine"); ok {
		cloudAzureIntegration.AzureVirtualmachine = expandCloudAzureIntegrationAzureVirtualMachineInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_virtual_machine"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVirtualmachine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_virtual_networks"); ok {
		cloudAzureIntegration.AzureVirtualnetworks = expandCloudAzureIntegrationAzureVirtualNetworksInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_virtual_networks"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVirtualnetworks = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_vms"); ok {
		cloudAzureIntegration.AzureVms = expandCloudAzureIntegrationAzureVmsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_vms"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVms = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("azure_vpn_gateway"); ok {
		cloudAzureIntegration.AzureVpngateways = expandCloudAzureIntegrationAzureVpnGatewayInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_vpn_gateway"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVpngateways = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	//// Unique for all resources

	configureInput := cloud.CloudIntegrationsInput{
		Azure: cloudAzureIntegration,
	}

	disableInput := cloud.CloudDisableIntegrationsInput{
		Azure: cloudDisableAzureIntegration,
	}

	return configureInput, disableInput
}

// Expanding the Azure API management

func expandCloudAzureIntegrationAzureAPIManagementInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureAPImanagementIntegrationInput {
	expanded := make([]cloud.CloudAzureAPImanagementIntegrationInput, len(b))

	for i, azureAPIManagement := range b {
		var azureAPIManagementInput cloud.CloudAzureAPImanagementIntegrationInput

		if azureAPIManagement == nil {
			azureAPIManagementInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureAPIManagementInput
			return expanded
		}

		in := azureAPIManagement.(map[string]interface{})

		azureAPIManagementInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureAPIManagementInput.MetricsPollingInterval = m.(int)
		}

		expanded[i] = azureAPIManagementInput
	}

	return expanded
}

// Expanding the Azure App Gateway

func expandCloudAzureIntegrationAzureAppGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureAppgatewayIntegrationInput {
	expanded := make([]cloud.CloudAzureAppgatewayIntegrationInput, len(b))

	for i, azureAppGateway := range b {
		var azureAppGatewayInput cloud.CloudAzureAppgatewayIntegrationInput

		if azureAppGateway == nil {
			azureAppGatewayInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureAppGatewayInput
			return expanded
		}

		in := azureAppGateway.(map[string]interface{})

		azureAppGatewayInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureAppGatewayInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureAppGatewayInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureAppGatewayInput
	}

	return expanded
}

// Expanding the Azure App service

func expandCloudAzureIntegrationAzureAppServiceInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureAppserviceIntegrationInput {
	expanded := make([]cloud.CloudAzureAppserviceIntegrationInput, len(b))

	for i, azureAppService := range b {
		var azureAppServiceInput cloud.CloudAzureAppserviceIntegrationInput

		if azureAppService == nil {
			azureAppServiceInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureAppServiceInput
			return expanded
		}

		in := azureAppService.(map[string]interface{})

		azureAppServiceInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureAppServiceInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureAppServiceInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureAppServiceInput
	}

	return expanded
}

// Expanding the Azure Containers

func expandCloudAzureIntegrationAzureContainersInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureContainersIntegrationInput {
	expanded := make([]cloud.CloudAzureContainersIntegrationInput, len(b))

	for i, azureContainers := range b {
		var azureContainersInput cloud.CloudAzureContainersIntegrationInput

		if azureContainers == nil {
			azureContainersInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureContainersInput
			return expanded
		}

		in := azureContainers.(map[string]interface{})

		azureContainersInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureContainersInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureContainersInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureContainersInput
	}

	return expanded
}

// Expanding the Azure Cosmosdb

func expandCloudAzureIntegrationAzureCosmosdbInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureCosmosdbIntegrationInput {
	expanded := make([]cloud.CloudAzureCosmosdbIntegrationInput, len(b))

	for i, azureCosmosdb := range b {
		var azureCosmosdbInput cloud.CloudAzureCosmosdbIntegrationInput

		if azureCosmosdb == nil {
			azureCosmosdbInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureCosmosdbInput
			return expanded
		}

		in := azureCosmosdb.(map[string]interface{})

		azureCosmosdbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureCosmosdbInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureCosmosdbInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureCosmosdbInput
	}

	return expanded
}

// Expanding the Azure Cost_management

func expandCloudAzureIntegrationAzureCostManagementInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureCostmanagementIntegrationInput {
	expanded := make([]cloud.CloudAzureCostmanagementIntegrationInput, len(b))

	for i, azureCostManagement := range b {
		var azureCostManagementInput cloud.CloudAzureCostmanagementIntegrationInput

		if azureCostManagement == nil {
			azureCostManagementInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureCostManagementInput
			return expanded
		}

		in := azureCostManagement.(map[string]interface{})

		azureCostManagementInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureCostManagementInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["tag_keys"]; ok {
			azureCostManagementInput.TagKeys[0] = r.(string)
		}
		expanded[i] = azureCostManagementInput
	}

	return expanded
}

// Expanding the Azure Data Factory

func expandCloudAzureIntegrationAzureDataFactoryInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureDatafactoryIntegrationInput {
	expanded := make([]cloud.CloudAzureDatafactoryIntegrationInput, len(b))

	for i, azureDataFactory := range b {
		var azureDataFactoryInput cloud.CloudAzureDatafactoryIntegrationInput

		if azureDataFactory == nil {
			azureDataFactoryInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureDataFactoryInput
			return expanded
		}

		in := azureDataFactory.(map[string]interface{})

		azureDataFactoryInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureDataFactoryInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureDataFactoryInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureDataFactoryInput
	}

	return expanded
}

// Expanding the Azure Event Hub

func expandCloudAzureIntegrationCloudAzureEventHubInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureEventhubIntegrationInput {
	expanded := make([]cloud.CloudAzureEventhubIntegrationInput, len(b))

	for i, azureEventHub := range b {
		var azureEventHubInput cloud.CloudAzureEventhubIntegrationInput

		if azureEventHub == nil {
			azureEventHubInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureEventHubInput
			return expanded
		}

		in := azureEventHub.(map[string]interface{})

		azureEventHubInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureEventHubInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureEventHubInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureEventHubInput
	}

	return expanded
}

// Expanding the Azure Express Route

func expandCloudAzureIntegrationAzureExpressRouteInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureExpressrouteIntegrationInput {
	expanded := make([]cloud.CloudAzureExpressrouteIntegrationInput, len(b))

	for i, azureExpressRoute := range b {
		var azureExpressRouteInput cloud.CloudAzureExpressrouteIntegrationInput

		if azureExpressRoute == nil {
			azureExpressRouteInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureExpressRouteInput
			return expanded
		}

		in := azureExpressRoute.(map[string]interface{})

		azureExpressRouteInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureExpressRouteInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureExpressRouteInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureExpressRouteInput
	}

	return expanded
}

// Expanding the azure_firewalls

func expandCloudAzureIntegrationAzureFirewallsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureFirewallsIntegrationInput {
	expanded := make([]cloud.CloudAzureFirewallsIntegrationInput, len(b))

	for i, azureFirewalls := range b {
		var azureFirewallsInput cloud.CloudAzureFirewallsIntegrationInput

		if azureFirewalls == nil {
			azureFirewallsInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureFirewallsInput
			return expanded
		}

		in := azureFirewalls.(map[string]interface{})

		azureFirewallsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureFirewallsInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureFirewallsInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureFirewallsInput
	}

	return expanded
}

// Expanding the Azure front_door

func expandCloudAzureIntegrationAzureFrontDoorInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureFrontdoorIntegrationInput {
	expanded := make([]cloud.CloudAzureFrontdoorIntegrationInput, len(b))

	for i, azureFrontDoor := range b {
		var azureFrontDoorInput cloud.CloudAzureFrontdoorIntegrationInput

		if azureFrontDoor == nil {
			azureFrontDoorInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureFrontDoorInput
			return expanded
		}

		in := azureFrontDoor.(map[string]interface{})

		azureFrontDoorInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureFrontDoorInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureFrontDoorInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureFrontDoorInput
	}

	return expanded
}

// Expanding the Azure Functions

func expandCloudAzureIntegrationAzureFunctionsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureFunctionsIntegrationInput {
	expanded := make([]cloud.CloudAzureFunctionsIntegrationInput, len(b))

	for i, azureFunctions := range b {
		var azureFunctionsInput cloud.CloudAzureFunctionsIntegrationInput

		if azureFunctions == nil {
			azureFunctionsInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureFunctionsInput
			return expanded
		}

		in := azureFunctions.(map[string]interface{})

		azureFunctionsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureFunctionsInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureFunctionsInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureFunctionsInput
	}

	return expanded
}

// Expanding the Azure Key Vault

func expandCloudAzureIntegrationAzureKeyVaultInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureKeyvaultIntegrationInput {
	expanded := make([]cloud.CloudAzureKeyvaultIntegrationInput, len(b))

	for i, azureKeyVault := range b {
		var azureKeyVaultInput cloud.CloudAzureKeyvaultIntegrationInput

		if azureKeyVault == nil {
			azureKeyVaultInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureKeyVaultInput
			return expanded
		}

		in := azureKeyVault.(map[string]interface{})

		azureKeyVaultInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureKeyVaultInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureKeyVaultInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureKeyVaultInput
	}

	return expanded
}

// Expanding the Azure Load Balancer

func expandCloudAzureIntegrationAzureLoadBalancerInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureLoadbalancerIntegrationInput {
	expanded := make([]cloud.CloudAzureLoadbalancerIntegrationInput, len(b))

	for i, azureLoadBalancer := range b {
		var azureLoadBalancerInput cloud.CloudAzureLoadbalancerIntegrationInput

		if azureLoadBalancer == nil {
			azureLoadBalancerInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureLoadBalancerInput
			return expanded
		}

		in := azureLoadBalancer.(map[string]interface{})

		azureLoadBalancerInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureLoadBalancerInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureLoadBalancerInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureLoadBalancerInput
	}

	return expanded
}

// Expanding the Azure Cosmosdb

func expandCloudAzureIntegrationAzureLogicAppsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureLogicappsIntegrationInput {
	expanded := make([]cloud.CloudAzureLogicappsIntegrationInput, len(b))

	for i, azureLogicApps := range b {
		var azureLogicAppsInput cloud.CloudAzureLogicappsIntegrationInput

		if azureLogicApps == nil {
			azureLogicAppsInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureLogicAppsInput
			return expanded
		}

		in := azureLogicApps.(map[string]interface{})

		azureLogicAppsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureLogicAppsInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureLogicAppsInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureLogicAppsInput
	}

	return expanded
}

// Expanding the azure_machine_learning

func expandCloudAzureIntegrationAzureMachineLearningInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMachinelearningIntegrationInput {
	expanded := make([]cloud.CloudAzureMachinelearningIntegrationInput, len(b))

	for i, azureMachineLearning := range b {
		var azureMachineLearningInput cloud.CloudAzureMachinelearningIntegrationInput

		if azureMachineLearning == nil {
			azureMachineLearningInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureMachineLearningInput
			return expanded
		}

		in := azureMachineLearning.(map[string]interface{})

		azureMachineLearningInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureMachineLearningInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureMachineLearningInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureMachineLearningInput
	}

	return expanded
}

// Expanding the azure_maria_db

func expandCloudAzureIntegrationAzureMariadbInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMariadbIntegrationInput {
	expanded := make([]cloud.CloudAzureMariadbIntegrationInput, len(b))

	for i, azureMariadb := range b {
		var azureMariadbInput cloud.CloudAzureMariadbIntegrationInput

		if azureMariadb == nil {
			azureMariadbInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureMariadbInput
			return expanded
		}

		in := azureMariadb.(map[string]interface{})

		azureMariadbInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureMariadbInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureMariadbInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureMariadbInput
	}

	return expanded
}

// Expanding the Azure_mysql

func expandCloudAzureIntegrationAzureMysqlInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMysqlIntegrationInput {
	expanded := make([]cloud.CloudAzureMysqlIntegrationInput, len(b))

	for i, azureMysql := range b {
		var azureMysqlInput cloud.CloudAzureMysqlIntegrationInput

		if azureMysql == nil {
			azureMysqlInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureMysqlInput
			return expanded
		}

		in := azureMysql.(map[string]interface{})

		azureMysqlInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureMysqlInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureMysqlInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureMysqlInput
	}

	return expanded
}

// Expanding the azure_postgresql

func expandCloudAzureIntegrationAzurePostgresqlInput(b []interface{}, linkedAccountID int) []cloud.CloudAzurePostgresqlIntegrationInput {
	expanded := make([]cloud.CloudAzurePostgresqlIntegrationInput, len(b))

	for i, azurePostgresql := range b {
		var azurePostgresqlInput cloud.CloudAzurePostgresqlIntegrationInput

		if azurePostgresql == nil {
			azurePostgresqlInput.LinkedAccountId = linkedAccountID
			expanded[i] = azurePostgresqlInput
			return expanded
		}

		in := azurePostgresql.(map[string]interface{})

		azurePostgresqlInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azurePostgresqlInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azurePostgresqlInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azurePostgresqlInput
	}

	return expanded
}

// Expanding the azure_power_bi_dedicated

func expandCloudAzureIntegrationAzurePowerBiDedicatedInput(b []interface{}, linkedAccountID int) []cloud.CloudAzurePowerbidedicatedIntegrationInput {
	expanded := make([]cloud.CloudAzurePowerbidedicatedIntegrationInput, len(b))

	for i, azurePowerBiDedicated := range b {
		var azurePowerBiDedicatedInput cloud.CloudAzurePowerbidedicatedIntegrationInput

		if azurePowerBiDedicated == nil {
			azurePowerBiDedicatedInput.LinkedAccountId = linkedAccountID
			expanded[i] = azurePowerBiDedicatedInput
			return expanded
		}

		in := azurePowerBiDedicated.(map[string]interface{})

		azurePowerBiDedicatedInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azurePowerBiDedicatedInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azurePowerBiDedicatedInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azurePowerBiDedicatedInput
	}

	return expanded
}

// Expanding the azure_redis_cache

func expandCloudAzureIntegrationAzureRedisCacheInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureRediscacheIntegrationInput {
	expanded := make([]cloud.CloudAzureRediscacheIntegrationInput, len(b))

	for i, azureRedisCache := range b {
		var azureRedisCacheInput cloud.CloudAzureRediscacheIntegrationInput

		if azureRedisCache == nil {
			azureRedisCacheInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureRedisCacheInput
			return expanded
		}

		in := azureRedisCache.(map[string]interface{})

		azureRedisCacheInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureRedisCacheInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureRedisCacheInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureRedisCacheInput
	}

	return expanded
}

// Expanding the azure_service_bus

func expandCloudAzureIntegrationAzureServiceBusInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureServicebusIntegrationInput {
	expanded := make([]cloud.CloudAzureServicebusIntegrationInput, len(b))

	for i, azureServiceBus := range b {
		var azureServiceBusInput cloud.CloudAzureServicebusIntegrationInput

		if azureServiceBus == nil {
			azureServiceBusInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureServiceBusInput
			return expanded
		}

		in := azureServiceBus.(map[string]interface{})

		azureServiceBusInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureServiceBusInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureServiceBusInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureServiceBusInput
	}

	return expanded
}

// Expanding the azure_service_fabric

func expandCloudAzureIntegrationAzureServiceFabricInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureServicefabricIntegrationInput {
	expanded := make([]cloud.CloudAzureServicefabricIntegrationInput, len(b))

	for i, azureServiceFabric := range b {
		var azureServiceFabricInput cloud.CloudAzureServicefabricIntegrationInput

		if azureServiceFabric == nil {
			azureServiceFabricInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureServiceFabricInput
			return expanded
		}

		in := azureServiceFabric.(map[string]interface{})

		azureServiceFabricInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureServiceFabricInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureServiceFabricInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureServiceFabricInput
	}

	return expanded
}

// Expanding the azure_sql

func expandCloudAzureIntegrationAzureSQLInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureSqlIntegrationInput {
	expanded := make([]cloud.CloudAzureSqlIntegrationInput, len(b))

	for i, azureSQL := range b {
		var azureSQLInput cloud.CloudAzureSqlIntegrationInput

		if azureSQL == nil {
			azureSQLInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureSQLInput
			return expanded
		}

		in := azureSQL.(map[string]interface{})

		azureSQLInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureSQLInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureSQLInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureSQLInput
	}

	return expanded
}

// Expanding the azure_sql_managed

func expandCloudAzureIntegrationAzureSQLManagedInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureSqlmanagedIntegrationInput {
	expanded := make([]cloud.CloudAzureSqlmanagedIntegrationInput, len(b))

	for i, azureSQLManaged := range b {
		var azureSQLManagedInput cloud.CloudAzureSqlmanagedIntegrationInput

		if azureSQLManaged == nil {
			azureSQLManagedInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureSQLManagedInput
			return expanded
		}

		in := azureSQLManaged.(map[string]interface{})

		azureSQLManagedInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureSQLManagedInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureSQLManagedInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureSQLManagedInput
	}

	return expanded
}

// Expanding the azure_storage

func expandCloudAzureIntegrationAzureStorageInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureStorageIntegrationInput {
	expanded := make([]cloud.CloudAzureStorageIntegrationInput, len(b))

	for i, azureStorage := range b {
		var azureStorageInput cloud.CloudAzureStorageIntegrationInput

		if azureStorage == nil {
			azureStorageInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureStorageInput
			return expanded
		}

		in := azureStorage.(map[string]interface{})

		azureStorageInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureStorageInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureStorageInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureStorageInput
	}

	return expanded
}

// Expanding the azure_virtual_machine

func expandCloudAzureIntegrationAzureVirtualMachineInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVirtualmachineIntegrationInput {
	expanded := make([]cloud.CloudAzureVirtualmachineIntegrationInput, len(b))

	for i, azureVirtualMachine := range b {
		var azureVirtualMachineInput cloud.CloudAzureVirtualmachineIntegrationInput

		if azureVirtualMachine == nil {
			azureVirtualMachineInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureVirtualMachineInput
			return expanded
		}

		in := azureVirtualMachine.(map[string]interface{})

		azureVirtualMachineInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureVirtualMachineInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureVirtualMachineInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureVirtualMachineInput
	}

	return expanded
}

// Expanding the azure_virtual_networks

func expandCloudAzureIntegrationAzureVirtualNetworksInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVirtualnetworksIntegrationInput {
	expanded := make([]cloud.CloudAzureVirtualnetworksIntegrationInput, len(b))

	for i, azureVirtualNetworks := range b {
		var azureVirtualNetworksInput cloud.CloudAzureVirtualnetworksIntegrationInput

		if azureVirtualNetworks == nil {
			azureVirtualNetworksInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureVirtualNetworksInput
			return expanded
		}

		in := azureVirtualNetworks.(map[string]interface{})

		azureVirtualNetworksInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureVirtualNetworksInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureVirtualNetworksInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureVirtualNetworksInput
	}

	return expanded
}

// Expanding the Azure vms

func expandCloudAzureIntegrationAzureVmsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVmsIntegrationInput {
	expanded := make([]cloud.CloudAzureVmsIntegrationInput, len(b))

	for i, azureVms := range b {
		var azureVmsInput cloud.CloudAzureVmsIntegrationInput

		if azureVms == nil {
			azureVmsInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureVmsInput
			return expanded
		}

		in := azureVms.(map[string]interface{})

		azureVmsInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureVmsInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureVmsInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureVmsInput
	}

	return expanded
}

// Expanding the azure_vpn_gateway

func expandCloudAzureIntegrationAzureVpnGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVpngatewaysIntegrationInput {
	expanded := make([]cloud.CloudAzureVpngatewaysIntegrationInput, len(b))

	for i, azureVpnGateway := range b {
		var azureVpnGatewayInput cloud.CloudAzureVpngatewaysIntegrationInput

		if azureVpnGateway == nil {
			azureVpnGatewayInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureVpnGatewayInput
			return expanded
		}

		in := azureVpnGateway.(map[string]interface{})

		azureVpnGatewayInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureVpnGatewayInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			azureVpnGatewayInput.ResourceGroups[0] = r.(string)
		}
		expanded[i] = azureVpnGatewayInput
	}

	return expanded
}

/// Read

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

	flattenCloudAzureLinkedAccount(d, linkedAccount)

	return nil
}

/// flatten

//nolint: gocyclo
func flattenCloudAzureLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("linked_account_id", result.ID)

	for _, i := range result.Integrations {
		switch t := i.(type) {
		case *cloud.CloudAzureAPImanagementIntegration:
			_ = d.Set("azure_api_management", flattenCloudAzureAPIManagementIntegration(t))
		case *cloud.CloudAzureAppgatewayIntegration:
			_ = d.Set("azure_app_gateway", flattenCloudAzureAppGatewayIntegration(t))
		case *cloud.CloudAzureAppserviceIntegration:
			_ = d.Set("azure_app_service", flattenCloudAzureAppServiceIntegration(t))
		case *cloud.CloudAzureContainersIntegration:
			_ = d.Set("azure_containers", flattenCloudAzureContainersIntegration(t))
		case *cloud.CloudAzureCosmosdbIntegration:
			_ = d.Set("azure_cosmosdb", flattenCloudAzureCosmosdbIntegration(t))
		case *cloud.CloudAzureCostmanagementIntegration:
			_ = d.Set("azure_data_factory", flattenCloudAzureCostManagementIntegration(t))
		case *cloud.CloudAzureDatafactoryIntegration:
			_ = d.Set("azure_data_factory", flattenCloudAzureDataFactoryIntegration(t))
		case *cloud.CloudAzureEventhubIntegration:
			_ = d.Set("azure_event_hub", flattenCloudAzureEventhubIntegration(t))
		case *cloud.CloudAzureExpressrouteIntegration:
			_ = d.Set("azure_express_route", flattenCloudAzureExpressRouteIntegration(t))
		case *cloud.CloudAzureFirewallsIntegration:
			_ = d.Set("azure_firewalls", flattenCloudAzureFirewallsIntegration(t))
		case *cloud.CloudAzureFrontdoorIntegration:
			_ = d.Set("azure_front_door", flattenCloudAzureFrontDoorIntegration(t))
		case *cloud.CloudAzureFunctionsIntegration:
			_ = d.Set("azure_functions", flattenCloudAzureFunctionsIntegration(t))
		case *cloud.CloudAzureKeyvaultIntegration:
			_ = d.Set("azure_key_vault", flattenCloudAzureKeyVaultIntegration(t))
		case *cloud.CloudAzureLoadbalancerIntegration:
			_ = d.Set("azure_load_balancer", flattenCloudAzureLoadBalancerIntegration(t))
		case *cloud.CloudAzureLogicappsIntegration:
			_ = d.Set("azure_logic_apps", flattenCloudAzureLogicAppsIntegration(t))
		case *cloud.CloudAzureMachinelearningIntegration:
			_ = d.Set("azure_machine_learning", flattenCloudAzureMachineLearningIntegration(t))
		case *cloud.CloudAzureMariadbIntegration:
			_ = d.Set("azure_maria_db", flattenCloudAzureMariadbIntegration(t))
		case *cloud.CloudAzureMysqlIntegration:
			_ = d.Set("azure_mysql", flattenCloudAzureMysqlIntegration(t))
		case *cloud.CloudAzurePostgresqlIntegration:
			_ = d.Set("azure_postgresql", flattenCloudAzurePostgresqlIntegration(t))
		case *cloud.CloudAzurePowerbidedicatedIntegration:
			_ = d.Set("azure_power_bi_dedicated", flattenCloudAzurePowerBIDedicatedIntegration(t))
		case *cloud.CloudAzureRediscacheIntegration:
			_ = d.Set("azure_redis_cache", flattenCloudAzureRedisCacheIntegration(t))
		case *cloud.CloudAzureServicebusIntegration:
			_ = d.Set("azure_service_bus", flattenCloudAzureServiceBusIntegration(t))
		case *cloud.CloudAzureServicefabricIntegration:
			_ = d.Set("azure_service_fabric", flattenCloudAzureServiceFabricIntegration(t))
		case *cloud.CloudAzureSqlIntegration:
			_ = d.Set("azure_sql", flattenCloudAzureSQLIntegration(t))
		case *cloud.CloudAzureSqlmanagedIntegration:
			_ = d.Set("azure_sql_managed", flattenCloudAzureSQLManagedIntegration(t))
		case *cloud.CloudAzureStorageIntegration:
			_ = d.Set("azure_storage", flattenCloudAzureStorageIntegration(t))
		case *cloud.CloudAzureVirtualmachineIntegration:
			_ = d.Set("azure_virtual_machine", flattenCloudAzureVirtualMachineIntegration(t))
		case *cloud.CloudAzureVirtualnetworksIntegration:
			_ = d.Set("azure_virtual_networks", flattenCloudAzureVirtualNetworksIntegration(t))
		case *cloud.CloudAzureVmsIntegration:
			_ = d.Set("azure_vms", flattenCloudAzureVmsIntegration(t))
		case *cloud.CloudAzureVpngatewaysIntegration:
			_ = d.Set("azure_vpn_gateway", flattenCloudAzureVpnGatewaysIntegration(t))

		}

	}
}

// flatten for API Management
func flattenCloudAzureAPIManagementIntegration(in *cloud.CloudAzureAPImanagementIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for App Gateway

func flattenCloudAzureAppGatewayIntegration(in *cloud.CloudAzureAppgatewayIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for App service

func flattenCloudAzureAppServiceIntegration(in *cloud.CloudAzureAppserviceIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for Containers

func flattenCloudAzureContainersIntegration(in *cloud.CloudAzureContainersIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for Cosmosdb

func flattenCloudAzureCosmosdbIntegration(in *cloud.CloudAzureCosmosdbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for cost management

func flattenCloudAzureCostManagementIntegration(in *cloud.CloudAzureCostmanagementIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["tags_keys"] = in.TagKeys

	flattened[0] = out

	return flattened
}

// flatten for data factory

func flattenCloudAzureDataFactoryIntegration(in *cloud.CloudAzureDatafactoryIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for Event hub

func flattenCloudAzureEventhubIntegration(in *cloud.CloudAzureEventhubIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for express route

func flattenCloudAzureExpressRouteIntegration(in *cloud.CloudAzureExpressrouteIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure firewalls

func flattenCloudAzureFirewallsIntegration(in *cloud.CloudAzureFirewallsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_front_door

func flattenCloudAzureFrontDoorIntegration(in *cloud.CloudAzureFrontdoorIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_functions

func flattenCloudAzureFunctionsIntegration(in *cloud.CloudAzureFunctionsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_key_vault

func flattenCloudAzureKeyVaultIntegration(in *cloud.CloudAzureKeyvaultIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_load_balancer

func flattenCloudAzureLoadBalancerIntegration(in *cloud.CloudAzureLoadbalancerIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_logic_apps

func flattenCloudAzureLogicAppsIntegration(in *cloud.CloudAzureLogicappsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_machine_learning

func flattenCloudAzureMachineLearningIntegration(in *cloud.CloudAzureMachinelearningIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_maria_db

func flattenCloudAzureMariadbIntegration(in *cloud.CloudAzureMariadbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_mysql

func flattenCloudAzureMysqlIntegration(in *cloud.CloudAzureMysqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_postgresql

func flattenCloudAzurePostgresqlIntegration(in *cloud.CloudAzurePostgresqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_power_bi_dedicated

func flattenCloudAzurePowerBIDedicatedIntegration(in *cloud.CloudAzurePowerbidedicatedIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_redis_cache

func flattenCloudAzureRedisCacheIntegration(in *cloud.CloudAzureRediscacheIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_service_bus

func flattenCloudAzureServiceBusIntegration(in *cloud.CloudAzureServicebusIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_service_fabric

func flattenCloudAzureServiceFabricIntegration(in *cloud.CloudAzureServicefabricIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_sql
func flattenCloudAzureSQLIntegration(in *cloud.CloudAzureSqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_sql_managed

func flattenCloudAzureSQLManagedIntegration(in *cloud.CloudAzureSqlmanagedIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_storage

func flattenCloudAzureStorageIntegration(in *cloud.CloudAzureStorageIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_virtual_machine

func flattenCloudAzureVirtualMachineIntegration(in *cloud.CloudAzureVirtualmachineIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_virtual_networks

func flattenCloudAzureVirtualNetworksIntegration(in *cloud.CloudAzureVirtualnetworksIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_vms

func flattenCloudAzureVmsIntegration(in *cloud.CloudAzureVmsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

// flatten for azure_vpn_gateway

func flattenCloudAzureVpnGatewaysIntegration(in *cloud.CloudAzureVpngatewaysIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups

	flattened[0] = out

	return flattened
}

/// update
func resourceNewRelicCloudAzureIntegrationsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	integrateInput, disableInput := expandCloudAzureIntegrationsInput(d)

	azureDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, disableInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(azureDisablePayload.Errors) > 0 {
		for _, err := range azureDisablePayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	azureIntegrationPayload, err := client.Cloud.CloudConfigureIntegrationWithContext(ctx, accountID, integrateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(azureIntegrationPayload.Errors) > 0 {
		for _, err := range azureIntegrationPayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}
	return nil
}

/// Delete
func resourceNewRelicCloudAzureIntegrationsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)
	deleteInput := expandCloudAzureDisableInputs(d)
	azureDisablePayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(azureDisablePayload.Errors) > 0 {
		for _, err := range azureDisablePayload.Errors {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  err.Type + " " + err.Message,
			})
		}
		return diags
	}

	d.SetId("")

	return nil
}

//nolint: gocyclo
func expandCloudAzureDisableInputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudAzureDisableInput := cloud.CloudAzureDisableIntegrationsInput{}
	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("azure_api_management"); ok {
		cloudAzureDisableInput.AzureAPImanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_app_gateway"); ok {
		cloudAzureDisableInput.AzureAppgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_app_service"); ok {
		cloudAzureDisableInput.AzureAppservice = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_containers"); ok {
		cloudAzureDisableInput.AzureContainers = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_cosmos_db"); ok {
		cloudAzureDisableInput.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_cost_management"); ok {
		cloudAzureDisableInput.AzureCostmanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_data_factory"); ok {
		cloudAzureDisableInput.AzureDatafactory = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_event_hub"); ok {
		cloudAzureDisableInput.AzureEventhub = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_express_route"); ok {
		cloudAzureDisableInput.AzureExpressroute = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_firewalls"); ok {
		cloudAzureDisableInput.AzureFirewalls = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_front_door"); ok {
		cloudAzureDisableInput.AzureFrontdoor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_functions"); ok {
		cloudAzureDisableInput.AzureFunctions = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_key_vault"); ok {
		cloudAzureDisableInput.AzureKeyvault = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_load_balancer"); ok {
		cloudAzureDisableInput.AzureLoadbalancer = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_logic_apps"); ok {
		cloudAzureDisableInput.AzureLogicapps = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_machine_learning"); ok {
		cloudAzureDisableInput.AzureMachinelearning = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_maria_db"); ok {
		cloudAzureDisableInput.AzureMariadb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_mysql"); ok {
		cloudAzureDisableInput.AzureMysql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_postgresql"); ok {
		cloudAzureDisableInput.AzurePostgresql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_power_bi_dedicated"); ok {
		cloudAzureDisableInput.AzurePowerbidedicated = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_redis_cache"); ok {
		cloudAzureDisableInput.AzureRediscache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_service_bus"); ok {
		cloudAzureDisableInput.AzureServicebus = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_service_fabric"); ok {
		cloudAzureDisableInput.AzureServicefabric = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_sql"); ok {
		cloudAzureDisableInput.AzureSql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_sql_managed"); ok {
		cloudAzureDisableInput.AzureSqlmanaged = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_storage"); ok {
		cloudAzureDisableInput.AzureStorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_virtual_machine"); ok {
		cloudAzureDisableInput.AzureVirtualmachine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_virtual_networks"); ok {
		cloudAzureDisableInput.AzureVirtualnetworks = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_vms"); ok {
		cloudAzureDisableInput.AzureVms = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("azure_vpn_gateway"); ok {
		cloudAzureDisableInput.AzureVpngateways = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	deleteInput := cloud.CloudDisableIntegrationsInput{
		Azure: cloudAzureDisableInput,
	}
	return deleteInput
}
