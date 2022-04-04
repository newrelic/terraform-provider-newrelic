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
				Elem:        cloudAzureIntegrationAzureAPIManagement(),
				MaxItems:    1,
			},
			"azure_app_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app gateway integration",
				Elem:        cloudAzureIntegrationAzureAppGateway(),
				MaxItems:    1,
			},
			"azure_app_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app services",
				Elem:        cloudAzureIntegrationAzureAppService(),
				MaxItems:    1,
			},
			"azure_containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure containers",
				Elem:        cloudAzureIntegrationAzureContainers,
				MaxItems:    1,
			},
			"azure_cosmos_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure cosmoDB",
				Elem:        cloudAzureIntegrationAzureCosmosDB,
				MaxItems:    1,
			},

			"azure_data_factory": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure data factory",
				Elem:        cloudAzureIntegrationAzureDataFactory,
				MaxItems:    1,
			},
			"azure_event_hub": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure event hub",
				Elem:        cloudAzureIntegrationAzureEventHub,
				MaxItems:    1,
			},
			"azure_express_route": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure express route",
				Elem:        cloudAzureIntegrationAzureExpressRoute,
				MaxItems:    1,
			},
			"azure_firewalls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure firewalls",
				Elem:        cloudAzureIntegrationAzureFirewalls,
				MaxItems:    1,
			},
			"azure_front_door": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure front door",
				Elem:        cloudAzureIntegrationAzureFrontDoor,
				MaxItems:    1,
			},
			"azure_functions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure functions",
				Elem:        cloudAzureIntegrationAzureFunctions,
				MaxItems:    1,
			},
			"azure_key_vault": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure key vault",
				Elem:        cloudAzureIntegrationAzureKeyVault,
				MaxItems:    1,
			},
			"azure_load_balancer": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure load balancer",
				Elem:        cloudAzureIntegrationAzureLoadBalancer,
				MaxItems:    1,
			},
			"azure_logic_apps": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure logic apps",
				Elem:        cloudAzureIntegrationAzureLogicApps,
				MaxItems:    1,
			},
			"azure_machine_learning": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure machine learning",
				Elem:        cloudAzureIntegrationAzureMachineLearning,
				MaxItems:    1,
			},
			"azure_maria_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Maria DB",
				Elem:        cloudAzureIntegrationAzureMariadb,
				MaxItems:    1,
			},
			"azure_mysql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure mysql",
				Elem:        cloudAzureIntegrationAzureMysql,
				MaxItems:    1,
			},
			"azure_postgresql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure postgresql",
				Elem:        cloudAzureIntegrationAzurePostgresql,
				MaxItems:    1,
			},
			"azure_power_bi_dedicated": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure powerBI dedicated",
				Elem:        cloudAzureIntegrationAzurePowerBiDedicated,
				MaxItems:    1,
			},
			"azure_redis_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure redis cache",
				Elem:        cloudAzureIntegrationAzureRedisCache,
				MaxItems:    1,
			},
			"azure_service_bus": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure service bus",
				Elem:        cloudAzureIntegrationAzureServiceBus,
				MaxItems:    1,
			},
			"azure_service_fabric": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The azure services fabric",
				Elem:        cloudAzureIntegrationAzureServiceFabric,
				MaxItems:    1,
			},
			"azure_sql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql",
				Elem:        cloudAzureIntegrationAzureSql,
				MaxItems:    1,
			},
			"azure_sql_managed": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql managed",
				Elem:        cloudAzureIntegrationAzureSqlManaged,
				MaxItems:    1,
			},
			"azure_storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure storage",
				Elem:        cloudAzureIntegrationAzureStorage,
				MaxItems:    1,
			},
			"azure_virtual_machine": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual machine",
				Elem:        cloudAzureIntegrationAzureVirtualMachine,
				MaxItems:    1,
			},
			"azure_virtual_networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual networks",
				Elem:        cloudAzureIntegrationAzureVirtualNetworks,
				MaxItems:    1,
			},
			"azure_vms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Vms",
				Elem:        cloudAzureIntegrationAzureVms,
				MaxItems:    1,
			},
			"azure_vpn_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure vpn gateway",
				Elem:        cloudAzureIntegrationAzureVPNGateway,
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

func cloudAzureIntegrationAzureAPIManagement() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAzureIntegrationAzureAppGateway() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}

func cloudAzureIntegrationAzureAppService() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureContainers() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureCosmosDB() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureDataFactory() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureEventHub() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureExpressRoute() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureFirewalls() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureFrontDoor() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureFunctions() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureKeyVault() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureLoadBalancer() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureLogicApps() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}
}
func cloudAzureIntegrationAzureMachineLearning() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureMariadb() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureMysql() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzurePostgresql() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzurePowerBiDedicated() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureRedisCache() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}

func cloudAzureIntegrationAzureServiceBus() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureServiceFabric() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureSql() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureSqlManaged() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureStorage() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureVirtualMachine() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}
func cloudAzureIntegrationAzureVirtualNetworks() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureVms() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()

	return &schema.Resource{
		Schema: s,
	}

}func cloudAzureIntegrationAzureVPNGateway() *schema.Resource {
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

	if v, ok := d.GetOk("azure_cosmosdb"); ok {
		cloudAzureIntegration.AzureCosmosdb = expandCloudAzureIntegrationAzureCosmosdbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_cosmosdb"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
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
		if r, ok := in["resource_groups"]; ok {
			azureAPIManagementInput.ResourceGroups[0] = r.(string)
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

		if azureAppGatewayInput == nil {
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

		if azureAppServiceInput == nil {
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

		case *cloud.CloudAzureKeyvaultIntegration
			_ = d.Set("azure_key_vault", flattenCloudAzureKeyVaultIntegration(t))

		case *cloud.CloudAzureLoadbalancerIntegration:
			_ = d.Set("azure_load_balancer", flattenCloudAzureLoadBalancerIntegration(t))

		case *cloud.CloudAzureLogicappsIntegration:
			_ = d.Set("azure_logic_apps", flattenCloudAzureLogicAppsIntegration(t))

			case *cloud.CloudAzureMachinelearningIntegration:
			_ = d.Set("azure_machine_learning", flattenCloudAzureMachineLearningIntegration(t))
		case *cloud.CloudAzureMariadbIntegration:
			_ = d.Set("azure_maria_db", flattenCloudAzureMariadbIntegration))
		case *cloud.CloudAzureMysqlIntegration:
			_ = d.Set("azure_mysql", flattenCloudAzureMysqlIntegration(t))

		case *cloud.CloudAzurePostgresqlIntegration:
			_ = d.Set("azure_postgresql", flattenCloudAzurePostgresqlIntegration(t))
		case *cloud.:
			_ = d.Set("azure_power_bi_dedicated", flattenCloudAzurePowerBIDedicatedIntegration(t))
		case *cloud.CloudAzureRediscacheIntegration:
			_ = d.Set("azure_redis_cache", flattenCloudAzureRedisCacheIntegration(t))
		case *cloud.CloudAzureServicebusIntegration:
			_ = d.Set("azure_service_bus", flattenCloudAzureServiceBusIntegration(t))
		case *cloud.CloudAzureServicefabricIntegration:
			_ = d.Set("azure_service_fabric", flattenCloudAzureServiceFabricIntegration(t))
		case *cloud.CloudAzureSqlIntegration:
			_ = d.Set("azure_sql", flattenCloudAzureSqlIntegration(t))
		case *cloud.CloudAzureSqlmanagedIntegration:
			_ = d.Set("azure_sql_managed", flattenCloudAzureSqlManagedIntegration(t))
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

	flattened[0] = out

	return flattened
}

// flatten for App Gateway

func flattenCloudAzureAppGatewayIntegration(in *cloud.CloudAzureAppgatewayIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for App service

func flattenCloudAzureAppServiceIntegration(in *cloud.CloudAzureAppserviceIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Containers

func flattenCloudAzureContainersIntegration(in *cloud.CloudAzureContainersIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Cosmosdb

func flattenCloudAzureCosmosdbIntegration(in *cloud.CloudAzureCosmosdbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for data factory

func flattenCloudAzureDataFactoryIntegration(in *cloud.CloudAzureDatafactoryIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for Event hub

func flattenCloudAzureEventhubIntegration(in *cloud.CloudAzureEventhubIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for express route

func flattenCloudAzureExpressRouteIntegration(in *cloud.CloudAzureExpressrouteIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure firewalls

func flattenCloudAzureFirewallsIntegration(in *cloud.CloudAzureFirewallsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_front_door

func flattenCloudAzureFrontDoorIntegration(in *cloud.CloudAzureFrontdoorIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_functions

func flattenCloudAzureFunctionsIntegration(in *cloud.CloudAzureFunctionsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_key_vault

func flattenCloudAzureKeyVaultIntegration(in *cloud.CloudAzureKeyvaultIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_load_balancer

func flattenCloudAzureLoadBalancerIntegration(in *cloud.CloudAzureLoadbalancerIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_logic_apps

func flattenCloudAzureLogicAppsIntegration(in *cloud.CloudAzureLogicappsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_machine_learning

func flattenCloudAzureMachineLearningIntegration(in *cloud.CloudAzureMachinelearningIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_maria_db

func flattenCloudAzureMariadbIntegration(in *cloud.CloudAzureMariadbIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_mysql

func flattenCloudAzureMysqlIntegration(in *cloud.CloudAzureMysqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_postgresql

func flattenCloudAzurePostgresqlIntegration(in *cloud.CloudAzurePostgresqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_power_bi_dedicated

func flattenCloudAzurePowerBIDedicatedIntegration(in *cloud.CloudAzurePowerbidedicatedIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_redis_cache

func flattenCloudAzureRedisCacheIntegration(in *cloud.CloudAzureRediscacheIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_service_bus

func flattenCloudAzureServiceBusIntegration(in *cloud.CloudAzureServicebusIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_service_fabric

func flattenCloudAzureServiceFabricIntegration(in *cloud.CloudAzureServicefabricIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_sql
func flattenCloudAzureSqlIntegration(in *cloud.CloudAzureSqlIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_sql_managed

func flattenCloudAzureSqlManagedIntegration(in *cloud.CloudAzureSqlmanagedIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_storage

func flattenCloudAzureStorageIntegration(in *cloud.CloudAzureStorageIntegration) []interface{}{
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_virtual_machine

func flattenCloudAzureVirtualMachineIntegration(in *cloud.CloudAzureVirtualmachineIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_virtual_networks

func flattenCloudAzureVirtualNetworksIntegration(in *cloud.CloudAzureVirtualnetworksIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_vms

func flattenCloudAzureVmsIntegration(in *cloud.CloudAzureVmsIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

	flattened[0] = out

	return flattened
}

// flatten for azure_vpn_gateway

func flattenCloudAzureVpnGatewaysIntegration(in *cloud.CloudAzureVpngatewaysIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval

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

	deleteInput := disableInput(d)
	cloudDisableIntegrationsPayload, err := client.Cloud.CloudDisableIntegrationWithContext(ctx, accountID, deleteInput)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	if len(cloudDisableIntegrationsPayload.Errors) > 0 {
		for _, err := range cloudDisableIntegrationsPayload.Errors {
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

func disableInput(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudDisableAzureIntegration := cloud.CloudAzureDisableIntegrationsInput{}
	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("azure_api_management"); ok {
		cloudDisableAzureIntegration.AzureAPImanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	deleteInput := cloud.CloudDisableIntegrationsInput{
		Azure: cloudDisableAzureIntegration,
	}
	return deleteInput
}
