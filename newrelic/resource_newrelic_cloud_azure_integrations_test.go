//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAzureIntegration_Basic(t *testing.T) {
	randName := acctest.RandString(5)
	resourceName := "newrelic_cloud_azure_integrations.bar"

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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAzureIntegrationsDestroy,
		Steps: []resource.TestStep{

			//Test: Create
			{
				Config: testAccNewRelicAzureIntegrationsConfig(testAzureApplicationID, testAzureClientSecretID, testAzureSubscriptionID, testAzureTenantID, randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAzureIntegrationsExist(resourceName),
				),
			},

			// Test: Update
			{
				Config: testAccNewRelicAzureIntegrationsConfigUpdated(testAzureApplicationID, testAzureClientSecretID, testAzureSubscriptionID, testAzureSubscriptionID, randName),
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
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			fmt.Errorf("An error occurred creating Azure integrations")
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAzureIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_azure_integrations" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("linked azure account still exists: #{err}")
		}
	}

	return nil
}

func testAccNewRelicAzureIntegrationsConfig(applicationID string, clientSecretID string, subscriptionID string, tenantID string, name string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_azure_link_account" "foo" {
  application_id = "%[1]s"
  client_secret = "%[2]s"
  subscription_id = "%[3]s"
  tenant_id = "%[4]s"
  name  = "%[5]s"
}
 
resource "newrelic_cloud_azure_integrations" "bar" {
  linked_account_id = newrelic_cloud_azure_link_account.foo.id
  
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

  service_fabric {
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

  vpn_gateway {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }
}
  `, applicationID, clientSecretID, subscriptionID, tenantID, name)
}

func testAccNewRelicAzureIntegrationsConfigUpdated(applicationID string, clientSecretID string, subscriptionID string, tenantID string, name string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_azure_link_account" "foo" {
  application_id  = "%[1]s"
  client_secret   = "%[2]s"
  subscription_id = "%[3]s"
  tenant_id       = "%[4]s"
  name            = "%[5]s"
}

resource "newrelic_cloud_azure_integrations" "bar" {
  linked_account_id = newrelic_cloud_azure_link_account.foo.id

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
  service_fabric {
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

  vpn_gateway {
    metrics_polling_interval = 3600
    resource_groups          = ["beyond"]
  }
}
  `, applicationID, clientSecretID, subscriptionID, tenantID, name)

}
