//go:build integration || CLOUD
// +build integration CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAzureIntegration_Basic(t *testing.T) {
	testAzureIntegrationName := fmt.Sprintf("tf_cloud_integrations_test_azure_%s", acctest.RandString(5))
	resourceName := "newrelic_cloud_azure_integrations.bar"

	// t.Skipf("Skipping test until we can get a better Azure test account")

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testAzureApplicationID := os.Getenv("INTEGRATION_TESTING_AZURE_APPLICATION_ID")
	if testAzureApplicationID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_APPLICATION_ID must be set for acceptance test")
	}

	testAzureClientSecretID := os.Getenv("INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID")
	if testAzureClientSecretID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_CLIENT_SECRET_ID must be set for acceptance test")
	}

	testAzureSubscriptionID := os.Getenv("INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID")
	if testAzureSubscriptionID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_SUBSCRIPTION_ID must be set for acceptance test")
	}

	testAzureTenantID := os.Getenv("INTEGRATION_TESTING_AZURE_TENANT_ID")
	if testAzureTenantID == "" {
		t.Skip("INTEGRATION_TESTING_AZURE_TENANT_ID must be set for acceptance test")
	}

	azureIntegrationsTestConfig := map[string]string{
		"name":            testAzureIntegrationName,
		"account_id":      strconv.Itoa(testSubAccountID),
		"application_id":  testAzureApplicationID,
		"client_secret":   testAzureClientSecretID,
		"subscription_id": testAzureSubscriptionID,
		"tenant_id":       testAzureTenantID,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "azure") },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAzureIntegrationsDestroy,
		Steps: []resource.TestStep{

			//Test: Create
			{
				Config: testAccNewRelicAzureIntegrationsConfig(azureIntegrationsTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAzureIntegrationsExist(resourceName),
				),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
			},

			// Test: Update
			{
				Config: testAccNewRelicAzureIntegrationsConfigUpdated(azureIntegrationsTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAzureIntegrationsExist(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNewRelicCloudAzureIntegrationsExist(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			return fmt.Errorf("An error occurred creating Azure integrations")
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAzureIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_azure_integrations" && r.Type != "newrelic_cloud_azure_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked azure account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicAzureIntegrationsCommonConfig(azureIntegrationsTestConfig map[string]string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  account_id = "` + azureIntegrationsTestConfig["account_id"] + `"
  alias      = "cloud-integration-provider"
}

resource "newrelic_cloud_azure_link_account" "foo" {
  provider        = newrelic.cloud-integration-provider
  application_id  = "` + azureIntegrationsTestConfig["application_id"] + `"
  client_secret   = "` + azureIntegrationsTestConfig["client_secret"] + `"
  subscription_id = "` + azureIntegrationsTestConfig["subscription_id"] + `"
  tenant_id       = "` + azureIntegrationsTestConfig["tenant_id"] + `"
  name            = "` + azureIntegrationsTestConfig["name"] + `"
  account_id      = "` + azureIntegrationsTestConfig["account_id"] + `"
}
`)
}

func testAccNewRelicAzureIntegrationsConfig(azureIntegrationsTestConfig map[string]string) string {
	return fmt.Sprintf(`
` + testAccNewRelicAzureIntegrationsCommonConfig(azureIntegrationsTestConfig) + `

resource "newrelic_cloud_azure_integrations" "bar" {
  provider          = newrelic.cloud-integration-provider
  linked_account_id = newrelic_cloud_azure_link_account.foo.id
  account_id        = "` + azureIntegrationsTestConfig["account_id"] + `"

  api_management {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  app_gateway {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  app_service {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  containers {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  cosmos_db {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }
  cost_management {
    metrics_polling_interval = 3600
  }

  data_factory {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  event_hub {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  express_route {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  firewalls {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]

  }
  front_door {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]

  }
  functions {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]

  }
  key_vault {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]

  }
  load_balancer {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  logic_apps {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  machine_learning {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  maria_db {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  monitor {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
    include_tags             = ["env:testing", "env:production"]
    exclude_tags             = ["env:staging"]
    enabled                  = true
    resource_types           = ["microsoft.datashare/accounts"]
  }

  mysql {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  postgresql {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  power_bi_dedicated {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  redis_cache {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  service_bus {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  sql {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  sql_managed {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  storage {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  virtual_machine {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  virtual_networks {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  vms {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  auto_discovery {
    metrics_polling_interval = 28800
  }

  vpn_gateway {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

}`)
}

func testAccNewRelicAzureIntegrationsConfigUpdated(azureIntegrationsTestConfig map[string]string) string {
	return fmt.Sprintf(`
` + testAccNewRelicAzureIntegrationsCommonConfig(azureIntegrationsTestConfig) + `

resource "newrelic_cloud_azure_integrations" "bar" {
  provider          = newrelic.cloud-integration-provider
  linked_account_id = newrelic_cloud_azure_link_account.foo.id
  account_id        = "` + azureIntegrationsTestConfig["account_id"] + `"

  api_management {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  app_gateway {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  app_service {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  containers {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  cosmos_db {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  cost_management {
    metrics_polling_interval = 3600
  }

  data_factory {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  express_route {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  event_hub {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  firewalls {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  front_door {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  functions {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  key_vault {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  load_balancer {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  logic_apps {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  machine_learning {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  maria_db {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  monitor {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
    include_tags             = ["env:production"]
    exclude_tags             = ["env:staging"]
    enabled                  = true
    resource_types           = ["microsoft.datashare/accounts", "microsoft.eventhub/clusters"]
  }

  mysql {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  mysql_flexible {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  postgresql {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  postgresql_flexible {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  power_bi_dedicated {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  redis_cache {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  service_bus {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]

  }

  sql {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  sql_managed {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  storage {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  virtual_machine {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  virtual_networks {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  vms {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

  auto_discovery {
    metrics_polling_interval = 28800
   }

  vpn_gateway {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }

}`)
}
