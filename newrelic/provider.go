package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/httpclient"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-newrelic/version"
)

// TerraformProviderProductUserAgent string used to identify this provider in User Agent requests
const TerraformProviderProductUserAgent = "terraform-provider-newrelic"

// Provider represents a resource provider in Terraform
func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_KEY", nil),
				Sensitive:   true,
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_URL", "https://api.newrelic.com/v2"),
			},
			"synthetics_api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_SYNTHETICS_API_URL", "https://synthetics.newrelic.com/synthetics/api/v3"),
			},
			"insights_account_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INSIGHTS_ACCOUNT_ID", nil),
				Sensitive:   true,
			},
			"insights_insert_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INSIGHTS_INSERT_KEY", nil),
				Sensitive:   true,
			},
			"insights_insert_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INSIGHTS_INSERT_URL", "https://insights-collector.newrelic.com/v1/accounts"),
			},
			"insights_query_key": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INSIGHTS_QUERY_KEY", nil),
				Sensitive:   true,
			},
			"insights_query_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INSIGHTS_QUERY_URL", "https://insights-api.newrelic.com/v1/accounts"),
			},
			"infra_api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INFRA_API_URL", "https://infra-api.newrelic.com/v2"),
			},
			"insecure_skip_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_SKIP_VERIFY", false),
			},
			"cacert_file": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_CACERT", ""),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"newrelic_alert_channel":      dataSourceNewRelicAlertChannel(),
			"newrelic_alert_policy":       dataSourceNewRelicAlertPolicy(),
			"newrelic_application":        dataSourceNewRelicApplication(),
			"newrelic_key_transaction":    dataSourceNewRelicKeyTransaction(),
			"newrelic_plugin":             dataSourceNewRelicPlugin(),
			"newrelic_plugin_component":   dataSourceNewRelicPluginComponent(),
			"newrelic_synthetics_monitor": dataSourceNewRelicSyntheticsMonitor(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"newrelic_alert_channel":              resourceNewRelicAlertChannel(),
			"newrelic_alert_condition":            resourceNewRelicAlertCondition(),
			"newrelic_alert_policy_channel":       resourceNewRelicAlertPolicyChannel(),
			"newrelic_alert_policy":               resourceNewRelicAlertPolicy(),
			"newrelic_plugins_alert_condition":    resourceNewRelicPluginsAlertCondition(),
			"newrelic_dashboard":                  resourceNewRelicDashboard(),
			"newrelic_infra_alert_condition":      resourceNewRelicInfraAlertCondition(),
			"newrelic_insights_event":             resourceNewRelicInsightsEvent(),
			"newrelic_nrql_alert_condition":       resourceNewRelicNrqlAlertCondition(),
			"newrelic_synthetics_alert_condition": resourceNewRelicSyntheticsAlertCondition(),
			"newrelic_synthetics_monitor":         resourceNewRelicSyntheticsMonitor(),
			"newrelic_synthetics_monitor_script":  resourceNewRelicSyntheticsMonitorScript(),
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
	config := Config{
		APIKey:             data.Get("api_key").(string),
		APIURL:             data.Get("api_url").(string),
		userAgent:          fmt.Sprintf("%s %s/%s", httpclient.TerraformUserAgent(terraformVersion), TerraformProviderProductUserAgent, version.ProviderVersion),
		InsecureSkipVerify: data.Get("insecure_skip_verify").(bool),
		CACertFile:         data.Get("cacert_file").(string),
	}
	log.Println("[INFO] Initializing New Relic client")

	client, err := config.Client()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic client: %s", err)
	}

	log.Println("[INFO] Initializing New Relic Synthetics client")

	syntheticsConfig := Config{
		APIURL:             data.Get("synthetics_api_url").(string),
		APIKey:             data.Get("api_key").(string),
		userAgent:          fmt.Sprintf("%s %s/%s", httpclient.TerraformUserAgent(terraformVersion), TerraformProviderProductUserAgent, version.ProviderVersion),
		InsecureSkipVerify: data.Get("insecure_skip_verify").(bool),
		CACertFile:         data.Get("cacert_file").(string),
	}

	clientSynthetics, err := syntheticsConfig.ClientSynthetics()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic Synthetics client: %s", err)
	}

	insightsInsertConfig := Config{
		InsightsAccountID: data.Get("insights_account_id").(string),
		InsightsInsertKey: data.Get("insights_insert_key").(string),
		InsightsInsertURL: data.Get("insights_insert_url").(string),
	}
	clientInsightsInsert, err := insightsInsertConfig.ClientInsightsInsert()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic Insights insert client: %s", err)
	}

	insightsQueryConfig := Config{
		InsightsAccountID: data.Get("insights_account_id").(string),
		InsightsQueryKey:  data.Get("insights_query_key").(string),
		InsightsQueryURL:  data.Get("insights_query_url").(string),
	}
	clientInsightsQuery, err := insightsQueryConfig.ClientInsightsQuery()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic Insights query client: %s", err)
	}

	infraConfig := Config{
		APIKey: data.Get("api_key").(string),
		APIURL: data.Get("infra_api_url").(string),
	}
	log.Println("[INFO] Initializing New Relic Infra client")

	clientInfra, err := infraConfig.ClientInfra()
	if err != nil {
		return nil, fmt.Errorf("error initializing New Relic Infra client: %s", err)
	}

	providerConfig := ProviderConfig{
		Client:               client,
		InfraClient:          clientInfra,
		InsightsInsertClient: clientInsightsInsert,
		InsightsQueryClient:  clientInsightsQuery,
		Synthetics:           clientSynthetics,
	}

	return &providerConfig, nil
}
