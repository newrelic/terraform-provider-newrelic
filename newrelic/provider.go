package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-newrelic/version"
)

// TerraformProviderProductUserAgent string used to identify this provider in User Agent requests
const TerraformProviderProductUserAgent = "terraform-provider-newrelic"

const (
	insightsInsertURL = "https://insights-collector.newrelic.com/v1/accounts"
	insightsQueryURL  = "https://insights-api.newrelic.com/v1/accounts"
)

// Provider represents a resource provider in Terraform
func Provider() terraform.ResourceProvider {
	deprecationMsgBaseURLs := "New Relic internal use only. API URLs are now configured based on the configured region."

	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_ACCOUNT_ID", nil),
				Sensitive:   true,
			},
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_API_KEY", nil),
				Sensitive:   true,
			},
			"admin_api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_ADMIN_API_KEY", nil),
				Sensitive:   true,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("NEW_RELIC_REGION", "US"),
				Description:  "The data center for which your New Relic account is configured. Only one region per provider block is permitted.",
				ValidateFunc: validation.StringInSlice([]string{"US", "EU", "Staging"}, true),
			},
			// New Relic internal use only
			"api_url": {
				Deprecated:  deprecationMsgBaseURLs,
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_API_URL", nil),
			},
			// New Relic internal use only
			"synthetics_api_url": {
				Deprecated:  deprecationMsgBaseURLs,
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_SYNTHETICS_API_URL", nil),
			},
			// New Relic internal use only
			"infrastructure_api_url": {
				Deprecated:  deprecationMsgBaseURLs,
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_INFRASTRUCTURE_API_URL", nil),
			},
			// New Relic internal use only
			"nerdgraph_api_url": {
				Deprecated:  deprecationMsgBaseURLs,
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_NERDGRAPH_API_URL", nil),
			},
			"insights_insert_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_INSIGHTS_INSERT_KEY", nil),
				Sensitive:   true,
			},
			"insights_insert_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_INSIGHTS_INSERT_URL", insightsInsertURL),
			},
			"insights_query_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_INSIGHTS_QUERY_KEY", nil),
				Sensitive:   true,
			},
			"insights_query_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_INSIGHTS_QUERY_URL", insightsQueryURL),
			},
			"insecure_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_API_SKIP_VERIFY", false),
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEW_RELIC_API_CACERT", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"newrelic_alert_channel":                dataSourceNewRelicAlertChannel(),
			"newrelic_alert_policy":                 dataSourceNewRelicAlertPolicy(),
			"newrelic_application":                  dataSourceNewRelicApplication(),
			"newrelic_entity":                       dataSourceNewRelicEntity(),
			"newrelic_key_transaction":              dataSourceNewRelicKeyTransaction(),
			"newrelic_plugin":                       dataSourceNewRelicPlugin(),
			"newrelic_plugin_component":             dataSourceNewRelicPluginComponent(),
			"newrelic_synthetics_monitor":           dataSourceNewRelicSyntheticsMonitor(),
			"newrelic_synthetics_secure_credential": dataSourceNewRelicSyntheticsSecureCredential(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"newrelic_alert_channel":                resourceNewRelicAlertChannel(),
			"newrelic_alert_condition":              resourceNewRelicAlertCondition(),
			"newrelic_alert_policy":                 resourceNewRelicAlertPolicy(),
			"newrelic_alert_policy_channel":         resourceNewRelicAlertPolicyChannel(),
			"newrelic_application_settings":         resourceNewRelicApplicationSettings(),
			"newrelic_dashboard":                    resourceNewRelicDashboard(),
			"newrelic_entity_tags":                  resourceNewRelicEntityTags(),
			"newrelic_infra_alert_condition":        resourceNewRelicInfraAlertCondition(),
			"newrelic_insights_event":               resourceNewRelicInsightsEvent(),
			"newrelic_nrql_alert_condition":         resourceNewRelicNrqlAlertCondition(),
			"newrelic_plugins_alert_condition":      resourceNewRelicPluginsAlertCondition(),
			"newrelic_synthetics_alert_condition":   resourceNewRelicSyntheticsAlertCondition(),
			"newrelic_synthetics_label":             resourceNewRelicSyntheticsLabel(),
			"newrelic_synthetics_monitor":           resourceNewRelicSyntheticsMonitor(),
			"newrelic_synthetics_monitor_script":    resourceNewRelicSyntheticsMonitorScript(),
			"newrelic_synthetics_secure_credential": resourceNewRelicSyntheticsSecureCredential(),
			"newrelic_workload":                     resourceNewRelicWorkload(),
		},
	}

	provider.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		terraformVersion := provider.TerraformVersion
		if terraformVersion == "" {
			// Catch for versions < 0.12
			terraformVersion = "0.11+compatible"
		}
		return providerConfigure(d, terraformVersion)
	}

	return provider
}

func providerConfigure(data *schema.ResourceData, terraformVersion string) (interface{}, error) {
	adminAPIKey := data.Get("admin_api_key").(string)
	personalAPIKey := data.Get("api_key").(string)
	terraformUA := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", terraformVersion, meta.SDKVersionString())
	userAgent := fmt.Sprintf("%s %s/%s", terraformUA, TerraformProviderProductUserAgent, version.ProviderVersion)
	accountID := data.Get("account_id").(int)

	cfg := Config{
		AdminAPIKey:          adminAPIKey,
		PersonalAPIKey:       personalAPIKey,
		Region:               data.Get("region").(string),
		APIURL:               data.Get("api_url").(string),
		SyntheticsAPIURL:     data.Get("synthetics_api_url").(string),
		NerdGraphAPIURL:      data.Get("nerdgraph_api_url").(string),
		InfrastructureAPIURL: getInfraAPIURL(data),
		userAgent:            userAgent,
		InsecureSkipVerify:   data.Get("insecure_skip_verify").(bool),
		CACertFile:           data.Get("cacert_file").(string),
	}
	log.Println("[INFO] Initializing newrelic-client-go")

	client, err := cfg.Client()
	if err != nil {
		return nil, fmt.Errorf("error initializing newrelic-client-go: %w", err)
	}

	insightsInsertConfig := Config{
		InsightsAccountID: strconv.Itoa(accountID),
		InsightsInsertKey: data.Get("insights_insert_key").(string),
		InsightsInsertURL: data.Get("insights_insert_url").(string),
	}
	clientInsightsInsert, err := insightsInsertConfig.ClientInsightsInsert()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic Insights insert client: %w", err)
	}

	insightsQueryConfig := Config{
		InsightsAccountID: strconv.Itoa(accountID),
		InsightsQueryKey:  data.Get("insights_query_key").(string),
		InsightsQueryURL:  data.Get("insights_query_url").(string),
	}
	clientInsightsQuery, err := insightsQueryConfig.ClientInsightsQuery()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic Insights query client: %s", err)
	}

	providerConfig := ProviderConfig{
		NewClient:            client,
		InsightsInsertClient: clientInsightsInsert,
		InsightsQueryClient:  clientInsightsQuery,
		PersonalAPIKey:       personalAPIKey,
		AccountID:            accountID,
	}

	return &providerConfig, nil
}

func getInfraAPIURL(data *schema.ResourceData) string {
	newURL, newURLOk := data.GetOk("infrastructure_api_url")

	if newURLOk {
		return newURL.(string)
	}

	return ""
}
