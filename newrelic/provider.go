package newrelic

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider represents a resource provider in Terraform
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_KEY", nil),
				Sensitive:   true,
			},
			"api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_API_URL", "https://api.newrelic.com/v2"),
			},
			"infra_api_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("NEWRELIC_INFRA_API_URL", "https://infra-api.newrelic.com/v2"),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"newrelic_application":        dataSourceNewRelicApplication(),
			"newrelic_key_transaction":    dataSourceNewRelicKeyTransaction(),
			"newrelic_synthetics_monitor": dataSourceNewRelicSyntheticsMonitor(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"newrelic_alert_channel":              resourceNewRelicAlertChannel(),
			"newrelic_alert_condition":            resourceNewRelicAlertCondition(),
			"newrelic_nrql_alert_condition":       resourceNewRelicNrqlAlertCondition(),
			"newrelic_synthetics_alert_condition": resourceNewRelicSyntheticsAlertCondition(),
			"newrelic_infra_alert_condition":      resourceNewRelicInfraAlertCondition(),
			"newrelic_alert_policy":               resourceNewRelicAlertPolicy(),
			"newrelic_alert_policy_channel":       resourceNewRelicAlertPolicyChannel(),
			"newrelic_dashboard":                  resourceNewRelicDashboard(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIKey: data.Get("api_key").(string),
		APIURL: data.Get("api_url").(string),
	}
	log.Println("[INFO] Initializing New Relic client")

	client, err := config.Client()
	if err != nil {
		return nil, fmt.Errorf("Error initializing New Relic client: %s", err)
	}

	log.Println("[INFO] Initializing New Relic Synthetics client")

	clientSynthetics, err := config.ClientSynthetics()
	if err != nil {
		return nil, fmt.Errorf("Error initializing New Relic synthetics client: %s", err)
	}

	infraConfig := Config{
		APIKey: data.Get("api_key").(string),
		APIURL: data.Get("infra_api_url").(string),
	}
	log.Println("[INFO] Initializing New Relic Infra client")

	clientInfra, err := infraConfig.ClientInfra()
	if err != nil {
		return nil, fmt.Errorf("Error initializing New Relic Infra client: %s", err)
	}

	providerConfig := ProviderConfig{
		Client:      client,
		InfraClient: clientInfra,
		Synthetics:  clientSynthetics,
	}

	return &providerConfig, nil
}
