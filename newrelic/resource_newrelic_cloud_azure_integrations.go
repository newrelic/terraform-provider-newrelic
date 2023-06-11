package newrelic

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/cloud"
)

func resourceNewRelicCloudAzureIntegrations() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicCloudAzureIntegrationsCreate,
		ReadContext:   resourceNewRelicCloudAzureIntegrationsRead,
		UpdateContext: resourceNewRelicCloudAzureIntegrationsUpdate,
		DeleteContext: resourceNewRelicCloudAzureIntegrationsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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

			"api_management": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure api management integration",
				Elem:        cloudAzureIntegrationAPIManagementElem(),
				MaxItems:    1,
			},
			"app_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app gateway integration",
				Elem:        cloudAzureIntegrationAppGatewayElem(),
				MaxItems:    1,
			},
			"app_service": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure app services",
				Elem:        cloudAzureIntegrationAppServiceElem(),
				MaxItems:    1,
			},
			"containers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure containers",
				Elem:        cloudAzureIntegrationContainersElem(),
				MaxItems:    1,
			},
			"cosmos_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure cosmoDB",
				Elem:        cloudAzureIntegrationCosmosDBElem(),
				MaxItems:    1,
			},
			"cost_management": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure cost management",
				Elem:        cloudAzureIntegrationCostManagementElem(),
				MaxItems:    1,
			},
			"data_factory": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure data factory",
				Elem:        cloudAzureIntegrationDataFactoryElem(),
				MaxItems:    1,
			},
			"event_hub": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure event hub",
				Elem:        cloudAzureIntegrationEventHubElem(),
				MaxItems:    1,
			},
			"express_route": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure express route",
				Elem:        cloudAzureIntegrationExpressRouteElem(),
				MaxItems:    1,
			},
			"firewalls": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure firewalls",
				Elem:        cloudAzureIntegrationFirewallsElem(),
				MaxItems:    1,
			},
			"front_door": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure front door",
				Elem:        cloudAzureIntegrationFrontDoorElem(),
				MaxItems:    1,
			},
			"functions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure functions",
				Elem:        cloudAzureIntegrationFunctionsElem(),
				MaxItems:    1,
			},
			"key_vault": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure key vault",
				Elem:        cloudAzureIntegrationKeyVaultElem(),
				MaxItems:    1,
			},
			"load_balancer": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure load balancer",
				Elem:        cloudAzureIntegrationLoadBalancerElem(),
				MaxItems:    1,
			},
			"logic_apps": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure logic apps",
				Elem:        cloudAzureIntegrationLogicAppsElem(),
				MaxItems:    1,
			},
			"machine_learning": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure machine learning",
				Elem:        cloudAzureIntegrationMachineLearningElem(),
				MaxItems:    1,
			},
			"maria_db": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Maria DB",
				Elem:        cloudAzureIntegrationMariadbElem(),
				MaxItems:    1,
			},
			"monitor": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Monitor",
				Elem:        cloudAzureIntegrationMonitorElem(),
				MaxItems:    1,
			},
			"mysql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure mysql",
				Elem:        cloudAzureIntegrationMysqlElem(),
				MaxItems:    1,
			},
			"mysql_flexible": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure mysql flexible service integration",
				Elem:        cloudAzureIntegrationMysqlFlexElem(),
				MaxItems:    1,
			},
			"postgresql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure postgresql",
				Elem:        cloudAzureIntegrationPostgresqlElem(),
				MaxItems:    1,
			},
			"postgresql_flexible": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure postgresql flexible service integration",
				Elem:        cloudAzureIntegrationPostgresqlFlexElem(),
				MaxItems:    1,
			},
			"power_bi_dedicated": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure powerBI dedicated",
				Elem:        cloudAzureIntegrationPowerBiDedicatedElem(),
				MaxItems:    1,
			},
			"redis_cache": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure redis cache",
				Elem:        cloudAzureIntegrationRedisCacheElem(),
				MaxItems:    1,
			},
			"service_bus": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure service bus",
				Elem:        cloudAzureIntegrationServiceBusElem(),
				MaxItems:    1,
			},
			"sql": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql",
				Elem:        cloudAzureIntegrationSQLElem(),
				MaxItems:    1,
			},
			"sql_managed": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure sql managed",
				Elem:        cloudAzureIntegrationSQLManagedElem(),
				MaxItems:    1,
			},
			"storage": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure storage",
				Elem:        cloudAzureIntegrationStorageElem(),
				MaxItems:    1,
			},
			"virtual_machine": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual machine",
				Elem:        cloudAzureIntegrationVirtualMachineElem(),
				MaxItems:    1,
			},
			"virtual_networks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure virtual networks",
				Elem:        cloudAzureIntegrationVirtualNetworksElem(),
				MaxItems:    1,
			},
			"vms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure Vms",
				Elem:        cloudAzureIntegrationVmsElem(),
				MaxItems:    1,
			},
			"vpn_gateway": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The Azure vpn gateway",
				Elem:        cloudAzureIntegrationVPNGatewayElem(),
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
	}
}

// function to add schema for azure API management
func cloudAzureIntegrationAPIManagementElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

//function to add schema for azure Gateway

func cloudAzureIntegrationAppGatewayElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for azure app service
func cloudAzureIntegrationAppServiceElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure containers
func cloudAzureIntegrationContainersElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure cosmo database
func cloudAzureIntegrationCosmosDBElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure cost management

func cloudAzureIntegrationCostManagementElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["tag_keys"] = &schema.Schema{
		Type:        schema.TypeList,
		Description: "Specify if additional cost data per tag should be collected. This field is case sensitive.\n\n",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for azure data factory

func cloudAzureIntegrationDataFactoryElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure event hub

func cloudAzureIntegrationEventHubElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure express route
func cloudAzureIntegrationExpressRouteElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure firewalls
func cloudAzureIntegrationFirewallsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure front door
func cloudAzureIntegrationFrontDoorElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure functions
func cloudAzureIntegrationFunctionsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure key vault
func cloudAzureIntegrationKeyVaultElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure load balancer
func cloudAzureIntegrationLoadBalancerElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure logic apps
func cloudAzureIntegrationLogicAppsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}
}

// function to add schema for azure machine learning
func cloudAzureIntegrationMachineLearningElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure maria database
func cloudAzureIntegrationMariadbElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure mysql
func cloudAzureIntegrationMysqlElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// cloudAzureIntegrationMonitorElem defines the schema of elements in the "monitor" Azure integration.
func cloudAzureIntegrationMonitorElem() *schema.Resource {
	s := mergeSchemas(
		cloudAzureIntegrationSchemaBase(),
		cloudAzureIntegrationMonitorSchema())

	return &schema.Resource{
		Schema: s,
	}
}

// cloudAzureIntegrationMonitorSchema defines the schema of elements specific to the "monitor" Azure integration.
func cloudAzureIntegrationMonitorSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "A flag that specifies if the integration is active",
		},
		"exclude_tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify resource tags in 'key:value' form to be excluded from monitoring",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"include_tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify resource tags in 'key:value' form to be monitored",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"resource_types": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify each Azure resource type that needs to be monitored",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"resource_groups": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// function to add schema for azure mysql
func cloudAzureIntegrationMysqlFlexElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure postgresql
func cloudAzureIntegrationPostgresqlElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure postgresql
func cloudAzureIntegrationPostgresqlFlexElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure power bi dedicated
func cloudAzureIntegrationPowerBiDedicatedElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure redis cache
func cloudAzureIntegrationRedisCacheElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure service bus
func cloudAzureIntegrationServiceBusElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure sql
func cloudAzureIntegrationSQLElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure sql managed
func cloudAzureIntegrationSQLManagedElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure storage
func cloudAzureIntegrationStorageElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure virtual machine
func cloudAzureIntegrationVirtualMachineElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure virtual networks
func cloudAzureIntegrationVirtualNetworksElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure vms
func cloudAzureIntegrationVmsElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

// function to add schema for azure VPN gateway
func cloudAzureIntegrationVPNGatewayElem() *schema.Resource {
	s := cloudAzureIntegrationSchemaBase()
	s["resource_groups"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Specify each Resource group associated with the resources that you want to monitor. Filter values are case-sensitive",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}
	return &schema.Resource{
		Schema: s,
	}

}

func resourceNewRelicCloudAzureIntegrationsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)

	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	cloudAzureIntegrationsInput, _ := expandCloudAzureIntegrationsInput(d)

	//cloudLinkAccountWithContext func which integrates azure account with Newrelic
	//which returns payload and error

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

// expand function to extract inputs from the schema.
// It takes ResourceData as input and returns CloudDisableIntegrationsInput.
// TODO: Reduce the cyclomatic complexity of this func
// nolint: gocyclo
func expandCloudAzureIntegrationsInput(d *schema.ResourceData) (cloud.CloudIntegrationsInput, cloud.CloudDisableIntegrationsInput) {
	cloudAzureIntegration := cloud.CloudAzureIntegrationsInput{}
	cloudDisableAzureIntegration := cloud.CloudAzureDisableIntegrationsInput{}

	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if v, ok := d.GetOk("api_management"); ok {
		cloudAzureIntegration.AzureAPImanagement = expandCloudAzureIntegrationAPIManagementInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("api_management"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAPImanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("app_gateway"); ok {
		cloudAzureIntegration.AzureAppgateway = expandCloudAzureIntegrationAppGatewayInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("app_gateway"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAppgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("app_service"); ok {
		cloudAzureIntegration.AzureAppservice = expandCloudAzureIntegrationAppServiceInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("app_service"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAppservice = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("containers"); ok {
		cloudAzureIntegration.AzureContainers = expandCloudAzureIntegrationContainersInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("containers"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureContainers = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("cosmos_db"); ok {
		cloudAzureIntegration.AzureCosmosdb = expandCloudAzureIntegrationCosmosdbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("cosmos_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("cost_management"); ok {
		cloudAzureIntegration.AzureCostmanagement = expandCloudAzureIntegrationCostManagementInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("cost_management"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureCostmanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if v, ok := d.GetOk("data_factory"); ok {
		cloudAzureIntegration.AzureDatafactory = expandCloudAzureIntegrationDataFactoryInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("data_factory"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureDatafactory = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("event_hub"); ok {
		cloudAzureIntegration.AzureEventhub = expandCloudAzureIntegrationCloudEventHubInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("event_hub"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureEventhub = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("express_route"); ok {
		cloudAzureIntegration.AzureExpressroute = expandCloudAzureIntegrationExpressRouteInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("express_route"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureExpressroute = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("firewalls"); ok {
		cloudAzureIntegration.AzureFirewalls = expandCloudAzureIntegrationFirewallsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("firewalls"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureFirewalls = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("front_door"); ok {
		cloudAzureIntegration.AzureFrontdoor = expandCloudAzureIntegrationFrontDoorInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("front_door"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureFrontdoor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("functions"); ok {
		cloudAzureIntegration.AzureFunctions = expandCloudAzureIntegrationFunctionsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("functions"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureFunctions = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("key_vault"); ok {
		cloudAzureIntegration.AzureKeyvault = expandCloudAzureIntegrationKeyVaultInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("key_vault"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureKeyvault = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("load_balancer"); ok {
		cloudAzureIntegration.AzureLoadbalancer = expandCloudAzureIntegrationLoadBalancerInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("load_balancer"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureLoadbalancer = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("logic_apps"); ok {
		cloudAzureIntegration.AzureLogicapps = expandCloudAzureIntegrationLogicAppsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("logic_apps"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureLogicapps = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	//
	if v, ok := d.GetOk("machine_learning"); ok {
		cloudAzureIntegration.AzureMachinelearning = expandCloudAzureIntegrationMachineLearningInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("machine_learning"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMachinelearning = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("maria_db"); ok {
		cloudAzureIntegration.AzureMariadb = expandCloudAzureIntegrationMariadbInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("maria_db"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMariadb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("monitor"); ok {
		cloudAzureIntegration.AzureMonitor = expandCloudAzureIntegrationMonitorInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("monitor"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMonitor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("mysql"); ok {
		cloudAzureIntegration.AzureMysql = expandCloudAzureIntegrationMysqlInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("mysql"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMysql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("mysql_flexible"); ok {
		cloudAzureIntegration.AzureMysqlflexible = expandCloudAzureIntegrationMysqlFlexibleInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("mysql_flexible"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureMysqlflexible = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("postgresql"); ok {
		cloudAzureIntegration.AzurePostgresql = expandCloudAzureIntegrationPostgresqlInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("postgresql"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzurePostgresql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("postgresql_flexible"); ok {
		cloudAzureIntegration.AzurePostgresqlflexible = expandCloudAzureIntegrationPostgresqlFlexibleInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("postgresql_flexible"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzurePostgresqlflexible = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("power_bi_dedicated"); ok {
		cloudAzureIntegration.AzurePowerbidedicated = expandCloudAzureIntegrationPowerBiDedicatedInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("power_bi_dedicated"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzurePowerbidedicated = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("redis_cache"); ok {
		cloudAzureIntegration.AzureRediscache = expandCloudAzureIntegrationRedisCacheInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("redis_cache"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureRediscache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("service_bus"); ok {
		cloudAzureIntegration.AzureServicebus = expandCloudAzureIntegrationServiceBusInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("service_bus"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureServicebus = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("sql"); ok {
		cloudAzureIntegration.AzureSql = expandCloudAzureIntegrationSQLInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("sql"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureSql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("sql_managed"); ok {
		cloudAzureIntegration.AzureSqlmanaged = expandCloudAzureIntegrationSQLManagedInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("sql_managed"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureSqlmanaged = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("storage"); ok {
		cloudAzureIntegration.AzureStorage = expandCloudAzureIntegrationStorageInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("storage"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureStorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("virtual_machine"); ok {
		cloudAzureIntegration.AzureVirtualmachine = expandCloudAzureIntegrationVirtualMachineInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("virtual_machine"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVirtualmachine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("virtual_networks"); ok {
		cloudAzureIntegration.AzureVirtualnetworks = expandCloudAzureIntegrationVirtualNetworksInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("virtual_networks"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVirtualnetworks = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("vms"); ok {
		cloudAzureIntegration.AzureVms = expandCloudAzureIntegrationVmsInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("vms"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVms = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if v, ok := d.GetOk("vpn_gateway"); ok {
		cloudAzureIntegration.AzureVpngateways = expandCloudAzureIntegrationVpnGatewayInput(v.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("vpn_gateway"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureVpngateways = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	configureInput := cloud.CloudIntegrationsInput{
		Azure: cloudAzureIntegration,
	}

	disableInput := cloud.CloudDisableIntegrationsInput{
		Azure: cloudDisableAzureIntegration,
	}

	return configureInput, disableInput
}

// Expanding the AzureAPIManagement

func expandCloudAzureIntegrationAPIManagementInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureAPImanagementIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureAPIManagementInput.ResourceGroups = groups
		}
		expanded[i] = azureAPIManagementInput
	}

	return expanded
}

// Expanding the Azure App Gateway

func expandCloudAzureIntegrationAppGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureAppgatewayIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureAppGatewayInput.ResourceGroups = groups
		}
		expanded[i] = azureAppGatewayInput
	}

	return expanded
}

// Expanding the Azure App service

func expandCloudAzureIntegrationAppServiceInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureAppserviceIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureAppServiceInput.ResourceGroups = groups
		}
		expanded[i] = azureAppServiceInput
	}

	return expanded
}

// Expanding the Azure Containers

func expandCloudAzureIntegrationContainersInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureContainersIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureContainersInput.ResourceGroups = groups
		}

		expanded[i] = azureContainersInput
	}

	return expanded
}

// Expanding the Azure Cosmosdb

func expandCloudAzureIntegrationCosmosdbInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureCosmosdbIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureCosmosdbInput.ResourceGroups = groups
		}
		expanded[i] = azureCosmosdbInput
	}

	return expanded
}

// Expanding the Azure Cost_management

func expandCloudAzureIntegrationCostManagementInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureCostmanagementIntegrationInput {
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
			tagKeys := r.([]interface{})
			var keys []string

			for _, key := range tagKeys {
				keys = append(keys, key.(string))
			}
			azureCostManagementInput.TagKeys = keys
		}
		expanded[i] = azureCostManagementInput
	}

	return expanded
}

// Expanding the Azure Data Factory

func expandCloudAzureIntegrationDataFactoryInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureDatafactoryIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureDataFactoryInput.ResourceGroups = groups
		}
		expanded[i] = azureDataFactoryInput
	}

	return expanded
}

// Expanding the Azure Event Hub

func expandCloudAzureIntegrationCloudEventHubInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureEventhubIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureEventHubInput.ResourceGroups = groups
		}
		expanded[i] = azureEventHubInput
	}

	return expanded
}

// Expanding the Azure Express Route

func expandCloudAzureIntegrationExpressRouteInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureExpressrouteIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureExpressRouteInput.ResourceGroups = groups
		}
		expanded[i] = azureExpressRouteInput
	}

	return expanded
}

// Expanding the azure_firewalls

func expandCloudAzureIntegrationFirewallsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureFirewallsIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureFirewallsInput.ResourceGroups = groups
		}
		expanded[i] = azureFirewallsInput
	}

	return expanded
}

// Expanding the Azure front_door

func expandCloudAzureIntegrationFrontDoorInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureFrontdoorIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureFrontDoorInput.ResourceGroups = groups
		}
		expanded[i] = azureFrontDoorInput
	}

	return expanded
}

// Expanding the Azure Functions

func expandCloudAzureIntegrationFunctionsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureFunctionsIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureFunctionsInput.ResourceGroups = groups
		}
		expanded[i] = azureFunctionsInput
	}

	return expanded
}

// Expanding the Azure Key Vault

func expandCloudAzureIntegrationKeyVaultInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureKeyvaultIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureKeyVaultInput.ResourceGroups = groups
		}
		expanded[i] = azureKeyVaultInput
	}

	return expanded
}

// Expanding the Azure Load Balancer

func expandCloudAzureIntegrationLoadBalancerInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureLoadbalancerIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureLoadBalancerInput.ResourceGroups = groups
		}
		expanded[i] = azureLoadBalancerInput
	}

	return expanded
}

// Expanding the Azure Cosmosdb

func expandCloudAzureIntegrationLogicAppsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureLogicappsIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureLogicAppsInput.ResourceGroups = groups
		}
		expanded[i] = azureLogicAppsInput
	}

	return expanded
}

// Expanding the azure_machine_learning

func expandCloudAzureIntegrationMachineLearningInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMachinelearningIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureMachineLearningInput.ResourceGroups = groups
		}
		expanded[i] = azureMachineLearningInput
	}

	return expanded
}

// Expanding the azure_maria_db

func expandCloudAzureIntegrationMariadbInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMariadbIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureMariadbInput.ResourceGroups = groups
		}
		expanded[i] = azureMariadbInput
	}

	return expanded
}

// Expanding the input for the azureMonitor integration
func expandCloudAzureIntegrationMonitorInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMonitorIntegrationInput {
	expanded := make([]cloud.CloudAzureMonitorIntegrationInput, len(b))

	for i, azureMonitor := range b {
		var azureMonitorInput cloud.CloudAzureMonitorIntegrationInput

		if azureMonitor == nil {
			azureMonitorInput.LinkedAccountId = linkedAccountID
			expanded[i] = azureMonitorInput
			return expanded
		}

		in := azureMonitor.(map[string]interface{})

		azureMonitorInput.LinkedAccountId = linkedAccountID

		if m, ok := in["metrics_polling_interval"]; ok {
			azureMonitorInput.MetricsPollingInterval = m.(int)
		}
		if r, ok := in["resource_groups"]; ok {
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureMonitorInput.ResourceGroups = groups
		}
		if rt, ok := in["resource_types"]; ok {
			resourceTypes := rt.([]interface{})
			var rTypes []string

			for _, rType := range resourceTypes {
				rTypes = append(rTypes, rType.(string))
			}
			azureMonitorInput.ResourceTypes = rTypes
		}
		if et, ok := in["exclude_tags"]; ok {
			excludeTags := et.([]interface{})
			var eTags []string

			for _, eTag := range excludeTags {
				eTags = append(eTags, eTag.(string))
			}
			azureMonitorInput.ExcludeTags = eTags
		}

		if it, ok := in["include_tags"]; ok {
			includeTags := it.([]interface{})
			var iTags []string

			for _, iTag := range includeTags {
				iTags = append(iTags, iTag.(string))
			}
			azureMonitorInput.IncludeTags = iTags
		}

		if enabled, ok := in["enabled"]; ok {
			azureMonitorInput.Enabled = enabled.(bool)
		}
		expanded[i] = azureMonitorInput
	}

	return expanded
}

// Expanding the Azure_mysql

func expandCloudAzureIntegrationMysqlInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMysqlIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureMysqlInput.ResourceGroups = groups
		}
		expanded[i] = azureMysqlInput
	}

	return expanded
}

func expandCloudAzureIntegrationMysqlFlexibleInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureMysqlflexibleIntegrationInput {
	expanded := make([]cloud.CloudAzureMysqlflexibleIntegrationInput, len(b))

	for i, azureMysql := range b {
		var azureMysqlInput cloud.CloudAzureMysqlflexibleIntegrationInput

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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureMysqlInput.ResourceGroups = groups
		}
		expanded[i] = azureMysqlInput
	}

	return expanded
}

// Expanding the azure_postgresql

func expandCloudAzureIntegrationPostgresqlInput(b []interface{}, linkedAccountID int) []cloud.CloudAzurePostgresqlIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azurePostgresqlInput.ResourceGroups = groups
		}
		expanded[i] = azurePostgresqlInput
	}

	return expanded
}

func expandCloudAzureIntegrationPostgresqlFlexibleInput(b []interface{}, linkedAccountID int) []cloud.CloudAzurePostgresqlflexibleIntegrationInput {
	expanded := make([]cloud.CloudAzurePostgresqlflexibleIntegrationInput, len(b))

	for i, azurePostgresql := range b {
		var azurePostgresqlInput cloud.CloudAzurePostgresqlflexibleIntegrationInput

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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azurePostgresqlInput.ResourceGroups = groups
		}
		expanded[i] = azurePostgresqlInput
	}

	return expanded
}

// Expanding the azure_power_bi_dedicated

func expandCloudAzureIntegrationPowerBiDedicatedInput(b []interface{}, linkedAccountID int) []cloud.CloudAzurePowerbidedicatedIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azurePowerBiDedicatedInput.ResourceGroups = groups
		}
		expanded[i] = azurePowerBiDedicatedInput
	}

	return expanded
}

// Expanding the azure_redis_cache

func expandCloudAzureIntegrationRedisCacheInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureRediscacheIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureRedisCacheInput.ResourceGroups = groups
		}
		expanded[i] = azureRedisCacheInput
	}

	return expanded
}

// Expanding the azure_service_bus

func expandCloudAzureIntegrationServiceBusInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureServicebusIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureServiceBusInput.ResourceGroups = groups
		}
		expanded[i] = azureServiceBusInput
	}

	return expanded
}

// Expanding the azure_sql

func expandCloudAzureIntegrationSQLInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureSqlIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureSQLInput.ResourceGroups = groups
		}
		expanded[i] = azureSQLInput
	}

	return expanded
}

// Expanding the azure_sql_managed

func expandCloudAzureIntegrationSQLManagedInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureSqlmanagedIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureSQLManagedInput.ResourceGroups = groups
		}
		expanded[i] = azureSQLManagedInput
	}

	return expanded
}

// Expanding the azure_storage

func expandCloudAzureIntegrationStorageInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureStorageIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureStorageInput.ResourceGroups = groups
		}
		expanded[i] = azureStorageInput
	}

	return expanded
}

// Expanding the azure_virtual_machine

func expandCloudAzureIntegrationVirtualMachineInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVirtualmachineIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureVirtualMachineInput.ResourceGroups = groups
		}
		expanded[i] = azureVirtualMachineInput
	}

	return expanded
}

// Expanding the azure_virtual_networks

func expandCloudAzureIntegrationVirtualNetworksInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVirtualnetworksIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureVirtualNetworksInput.ResourceGroups = groups
		}
		expanded[i] = azureVirtualNetworksInput
	}

	return expanded
}

// Expanding the Azure vms

func expandCloudAzureIntegrationVmsInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVmsIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureVmsInput.ResourceGroups = groups
		}
		expanded[i] = azureVmsInput
	}

	return expanded
}

// Expanding the azure_vpn_gateway

func expandCloudAzureIntegrationVpnGatewayInput(b []interface{}, linkedAccountID int) []cloud.CloudAzureVpngatewaysIntegrationInput {
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
			resourceGroups := r.([]interface{})
			var groups []string

			for _, group := range resourceGroups {
				groups = append(groups, group.(string))
			}
			azureVpnGatewayInput.ResourceGroups = groups
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
		if strings.Contains(err.Error(), "not found") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	flattenCloudAzureLinkedAccount(d, linkedAccount)

	return nil
}

/// flatten

// nolint: gocyclo
func flattenCloudAzureLinkedAccount(d *schema.ResourceData, result *cloud.CloudLinkedAccount) {
	_ = d.Set("account_id", result.NrAccountId)
	_ = d.Set("linked_account_id", result.ID)

	for _, i := range result.Integrations {
		switch t := i.(type) {
		case *cloud.CloudAzureAPImanagementIntegration:
			_ = d.Set("api_management", flattenCloudAPIManagementIntegration(t))
		case *cloud.CloudAzureAppgatewayIntegration:
			_ = d.Set("app_gateway", flattenCloudAzureAppGatewayIntegration(t))
		case *cloud.CloudAzureAppserviceIntegration:
			_ = d.Set("app_service", flattenCloudAzureAppServiceIntegration(t))
		case *cloud.CloudAzureContainersIntegration:
			_ = d.Set("containers", flattenCloudAzureContainersIntegration(t))
		case *cloud.CloudAzureCosmosdbIntegration:
			_ = d.Set("cosmos_db", flattenCloudAzureCosmosdbIntegration(t))
		case *cloud.CloudAzureCostmanagementIntegration:
			_ = d.Set("cost_management", flattenCloudAzureCostManagementIntegration(t))
		case *cloud.CloudAzureDatafactoryIntegration:
			_ = d.Set("data_factory", flattenCloudAzureDataFactoryIntegration(t))
		case *cloud.CloudAzureEventhubIntegration:
			_ = d.Set("event_hub", flattenCloudAzureEventhubIntegration(t))
		case *cloud.CloudAzureExpressrouteIntegration:
			_ = d.Set("express_route", flattenCloudAzureExpressRouteIntegration(t))
		case *cloud.CloudAzureFirewallsIntegration:
			_ = d.Set("firewalls", flattenCloudAzureFirewallsIntegration(t))
		case *cloud.CloudAzureFrontdoorIntegration:
			_ = d.Set("front_door", flattenCloudAzureFrontDoorIntegration(t))
		case *cloud.CloudAzureFunctionsIntegration:
			_ = d.Set("functions", flattenCloudAzureFunctionsIntegration(t))
		case *cloud.CloudAzureKeyvaultIntegration:
			_ = d.Set("key_vault", flattenCloudAzureKeyVaultIntegration(t))
		case *cloud.CloudAzureLoadbalancerIntegration:
			_ = d.Set("load_balancer", flattenCloudAzureLoadBalancerIntegration(t))
		case *cloud.CloudAzureLogicappsIntegration:
			_ = d.Set("logic_apps", flattenCloudAzureLogicAppsIntegration(t))
		case *cloud.CloudAzureMachinelearningIntegration:
			_ = d.Set("machine_learning", flattenCloudAzureMachineLearningIntegration(t))
		case *cloud.CloudAzureMariadbIntegration:
			_ = d.Set("maria_db", flattenCloudAzureMariadbIntegration(t))
		case *cloud.CloudAzureMonitorIntegration:
			_ = d.Set("monitor", flattenCloudAzureMonitorIntegration(t))
		case *cloud.CloudAzureMysqlIntegration:
			_ = d.Set("mysql", flattenCloudAzureMysqlIntegration(t))
		case *cloud.CloudAzureMysqlflexibleIntegration:
			_ = d.Set("mysql_flexible", flattenCloudAzureMysqlFlexibleIntegration(t))
		case *cloud.CloudAzurePostgresqlIntegration:
			_ = d.Set("postgresql", flattenCloudAzurePostgresqlIntegration(t))
		case *cloud.CloudAzurePostgresqlflexibleIntegration:
			_ = d.Set("postgresql_flexible", flattenCloudAzurePostgresqlFlexibleIntegration(t))
		case *cloud.CloudAzurePowerbidedicatedIntegration:
			_ = d.Set("power_bi_dedicated", flattenCloudAzurePowerBIDedicatedIntegration(t))
		case *cloud.CloudAzureRediscacheIntegration:
			_ = d.Set("redis_cache", flattenCloudAzureRedisCacheIntegration(t))
		case *cloud.CloudAzureServicebusIntegration:
			_ = d.Set("service_bus", flattenCloudAzureServiceBusIntegration(t))
		case *cloud.CloudAzureSqlIntegration:
			_ = d.Set("sql", flattenCloudAzureSQLIntegration(t))
		case *cloud.CloudAzureSqlmanagedIntegration:
			_ = d.Set("sql_managed", flattenCloudAzureSQLManagedIntegration(t))
		case *cloud.CloudAzureStorageIntegration:
			_ = d.Set("storage", flattenCloudAzureStorageIntegration(t))
		case *cloud.CloudAzureVirtualmachineIntegration:
			_ = d.Set("virtual_machine", flattenCloudAzureVirtualMachineIntegration(t))
		case *cloud.CloudAzureVirtualnetworksIntegration:
			_ = d.Set("virtual_networks", flattenCloudAzureVirtualNetworksIntegration(t))
		case *cloud.CloudAzureVmsIntegration:
			_ = d.Set("vms", flattenCloudAzureVmsIntegration(t))
		case *cloud.CloudAzureVpngatewaysIntegration:
			_ = d.Set("vpn_gateway", flattenCloudAzureVpnGatewaysIntegration(t))

		}

	}
}

// flatten for API Management
func flattenCloudAPIManagementIntegration(in *cloud.CloudAzureAPImanagementIntegration) []interface{} {
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
	out["tag_keys"] = in.TagKeys

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

// Flatten values for the azureMonitor integration
func flattenCloudAzureMonitorIntegration(in *cloud.CloudAzureMonitorIntegration) []interface{} {
	flattened := make([]interface{}, 1)

	out := make(map[string]interface{})

	out["metrics_polling_interval"] = in.MetricsPollingInterval
	out["resource_groups"] = in.ResourceGroups
	out["exclude_tags"] = in.ExcludeTags
	out["include_tags"] = in.IncludeTags
	out["resource_types"] = in.ResourceTypes
	out["enabled"] = in.Enabled

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

func flattenCloudAzureMysqlFlexibleIntegration(in *cloud.CloudAzureMysqlflexibleIntegration) []interface{} {
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

func flattenCloudAzurePostgresqlFlexibleIntegration(in *cloud.CloudAzurePostgresqlflexibleIntegration) []interface{} {
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

// / Delete
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

// nolint: gocyclo
func expandCloudAzureDisableInputs(d *schema.ResourceData) cloud.CloudDisableIntegrationsInput {
	cloudAzureDisableInput := cloud.CloudAzureDisableIntegrationsInput{}
	var linkedAccountID int

	if l, ok := d.GetOk("linked_account_id"); ok {
		linkedAccountID = l.(int)
	}
	if _, ok := d.GetOk("api_management"); ok {
		cloudAzureDisableInput.AzureAPImanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("app_gateway"); ok {
		cloudAzureDisableInput.AzureAppgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("app_service"); ok {
		cloudAzureDisableInput.AzureAppservice = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("containers"); ok {
		cloudAzureDisableInput.AzureContainers = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("cosmos_db"); ok {
		cloudAzureDisableInput.AzureCosmosdb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("cost_management"); ok {
		cloudAzureDisableInput.AzureCostmanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("data_factory"); ok {
		cloudAzureDisableInput.AzureDatafactory = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("event_hub"); ok {
		cloudAzureDisableInput.AzureEventhub = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("express_route"); ok {
		cloudAzureDisableInput.AzureExpressroute = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("firewalls"); ok {
		cloudAzureDisableInput.AzureFirewalls = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("front_door"); ok {
		cloudAzureDisableInput.AzureFrontdoor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("functions"); ok {
		cloudAzureDisableInput.AzureFunctions = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("key_vault"); ok {
		cloudAzureDisableInput.AzureKeyvault = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("load_balancer"); ok {
		cloudAzureDisableInput.AzureLoadbalancer = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("logic_apps"); ok {
		cloudAzureDisableInput.AzureLogicapps = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("machine_learning"); ok {
		cloudAzureDisableInput.AzureMachinelearning = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("maria_db"); ok {
		cloudAzureDisableInput.AzureMariadb = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("monitor"); ok {
		cloudAzureDisableInput.AzureMonitor = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("mysql"); ok {
		cloudAzureDisableInput.AzureMysql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("mysql_flexible"); ok {
		cloudAzureDisableInput.AzureMysqlflexible = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("postgresql"); ok {
		cloudAzureDisableInput.AzurePostgresql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("postgresql_flexible"); ok {
		cloudAzureDisableInput.AzurePostgresqlflexible = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("power_bi_dedicated"); ok {
		cloudAzureDisableInput.AzurePowerbidedicated = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("redis_cache"); ok {
		cloudAzureDisableInput.AzureRediscache = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("service_bus"); ok {
		cloudAzureDisableInput.AzureServicebus = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("sql"); ok {
		cloudAzureDisableInput.AzureSql = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("sql_managed"); ok {
		cloudAzureDisableInput.AzureSqlmanaged = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("storage"); ok {
		cloudAzureDisableInput.AzureStorage = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("virtual_machine"); ok {
		cloudAzureDisableInput.AzureVirtualmachine = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("virtual_networks"); ok {
		cloudAzureDisableInput.AzureVirtualnetworks = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("vms"); ok {
		cloudAzureDisableInput.AzureVms = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if _, ok := d.GetOk("vpn_gateway"); ok {
		cloudAzureDisableInput.AzureVpngateways = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	deleteInput := cloud.CloudDisableIntegrationsInput{
		Azure: cloudAzureDisableInput,
	}
	return deleteInput
}
