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
				Elem:        "",
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
	if aam, ok := d.GetOk("azure_api_management"); ok {
		cloudAzureIntegration.AzureAPImanagement = expandCloudAzureIntegrationAzureAPIManagementInput(aam.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_api_management"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAPImanagement = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}

	if aag, ok := d.GetOk("azure_app_gateway"); ok {
		cloudAzureIntegration.AzureAppgateway = expandCloudAzureIntegrationAzureAppGatewayInput(aag.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_app_gateway"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAppgateway = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
	}
	if aps, ok := d.GetOk("azure_app_service"); ok {
		cloudAzureIntegration.AzureAppservice = expandCloudAzureIntegrationAzureAppServiceInput(aps.([]interface{}), linkedAccountID)
	} else if o, n := d.GetChange("azure_app_service"); len(n.([]interface{})) < len(o.([]interface{})) {
		cloudDisableAzureIntegration.AzureAppservice = []cloud.CloudDisableAccountIntegrationInput{{LinkedAccountId: linkedAccountID}}
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

func expandCloudAzureIntegrationAzureAPIManagementInput(a []interface{}, linkedAccountID int) []cloud.CloudAzureAPImanagementIntegrationInput {
	expanded := make([]cloud.CloudAzureAPImanagementIntegrationInput, len(a))

	for i, azureAPIManagement := range a {
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

func expandCloudAzureIntegrationAzureAppServiceInput(a []interface{}, linkedAccountID int) []cloud.CloudAzureAppserviceIntegrationInput {
	expanded := make([]cloud.CloudAzureAppserviceIntegrationInput, len(a))

	for i, azureAppService := range a {
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
			_ = d.Set("azure_app_gateway", flattenCloudAzureAppgatewayIntegration(t))

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

func flattenCloudAzureAppgatewayIntegration(in *cloud.CloudAzureAppgatewayIntegration) []interface{} {
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
