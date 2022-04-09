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

func TestAccNewRelicCloudAzureLinkAccount_Basic(t *testing.T) {
	t.Skipf("Skipping test until integrations work is finished")
	randName := acctest.RandString(5)
	resourceName := "newrelic_cloud_azure_link_account.foo"

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
				Config: testAccNewRelicAzureIntegartionsConfig(testAzureApplicationID, testAzureClientSecretID, testAzureSubscriptionID, testAzureTenantID, randName),
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
	resource "newrelic_cloud_azure_link_account" "foo1"{
	account_id = 2520528
	application_id = "%[1]s"
	client_secret = "%[2]s"
	subscription_id = "%[3]s"
	tenant_id = "%[4]s"
	name  = "%[5]s"
	}
	
	resource "newrelic_cloud_azure_integrations" "foo" {
		linked_account_id = newrelic_cloud_azure_link_account.foo.id
		account_id = 2520528
		azure_api_management {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
		}
		
		azure_app_gateway {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
		}
		azure_app_service"{
				metrics_polling_interval = 1200
				resource_groups = "beyond"
		}
		azure_containers {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
					
		azure_cosmos_db {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_cost_management {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		
		azure_data_factory {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		
		azure_event_hub {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		
		azure_express_route {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_event_hub {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_firewalls {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_front_door {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_functions {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_key_vault {
				metrics_polling_interval = 1200
				resource_groups = "beyond"
				}
		azure_load_balancer { 
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
				}
		azure_logic_apps {
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		azure_machine_learning {
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		azure_maria_db {
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		
		azure_mysql {
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		azure_postgresql{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		azure_power_bi_dedicated{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
				}
		azure_power_bi_dedicated{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
				}
		azure_redis_cache{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"	
				}
		
		azure_service_bus{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"	
				}
		
		azure_sql{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		
		azure_storage{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		
		azure_virtual_machine{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		
		azure_virtual_networks{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		
		azure_vms{
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
		
		azure_vpn_gateway{	
				metrics_polling_interval = 1200	
				resource_groups = "beyond"
		
				}
	`, applicationID, clientSecretID, subscriptionID, tenantID, name)

}

func testAccNewRelicAzureIntegrationsConfigUpdated(applicationID string, clientSecretID string, subscriptionID string, tenantID string, name string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_azure_link_account" "foo1"{
	account_id = 2520528
	application_id = "%[1]s"
	client_secret = "%[2]s"
	subscription_id = "%[3]s"
	tenant_id = "%[4]s"
	name  = "%[5]s"
	}
	
	resource "newrelic_cloud_azure_integrations" "foo" {
		linked_account_id = newrelic_cloud_azure_link_account.foo.id
		account_id = 2520528
		azure_api_management {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
		}
		
		azure_app_gateway {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
		}
		azure_app_service"{
				metrics_polling_interval = 1000
				resource_groups = "beyond"
		}
		azure_containers {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
					
		azure_cosmos_db {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_cost_management {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		
		azure_data_factory {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		
		azure_event_hub {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		
		azure_express_route {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_event_hub {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_firewalls {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_front_door {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_functions {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_key_vault {
				metrics_polling_interval = 1000
				resource_groups = "beyond"
				}
		azure_load_balancer { 
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
				}
		azure_logic_apps {
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		azure_machine_learning {
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		azure_maria_db {
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		
		azure_mysql {
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		azure_postgresql{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		azure_power_bi_dedicated{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
				}
		azure_power_bi_dedicated{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
				}
		azure_redis_cache{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"	
				}
		
		azure_service_bus{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"	
				}
		
		azure_sql{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		
		azure_storage{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		
		azure_virtual_machine{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		
		azure_virtual_networks{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		
		azure_vms{
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
		
		azure_vpn_gateway{	
				metrics_polling_interval = 1000	
				resource_groups = "beyond"
		
				}
	`, applicationID, clientSecretID, subscriptionID, tenantID, name)

}
