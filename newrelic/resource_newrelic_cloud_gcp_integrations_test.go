//go:build integration || CLOUD

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

func TestAccNewRelicCloudGcpIntegrations_Basic(t *testing.T) {
	t.Skipf("Skipping test until GCP Environment Variables are fixed")
	resourceName := "newrelic_cloud_gcp_integrations.foo1"
	testGCPIntegrationName := fmt.Sprintf("tf_cloud_integrations_test_gcp_%s", acctest.RandString(5))

	if subAccountIDExists := os.Getenv("NEW_RELIC_SUBACCOUNT_ID"); subAccountIDExists == "" {
		t.Skipf("Skipping this test, as NEW_RELIC_SUBACCOUNT_ID must be set for this test to run.")
	}

	testGCPProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	if testGCPProjectID == "" {
		t.Skipf("INTEGRATION_TESTING_GCP_PROJECT_ID must be set for acceptance test")
	}

	GCPIntegrationTestConfig := map[string]string{
		"name":       testGCPIntegrationName,
		"account_id": strconv.Itoa(testSubAccountID),
		"project_id": testGCPProjectID,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccCloudLinkedAccountsCleanup(t, "gcp") },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudGcpIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudGcpIntegrationsConfig(GCPIntegrationTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpIntegrationsExists(resourceName),
				),
				PreConfig: func() {
					time.Sleep(10 * time.Second)
				},
			},
			//Test: Update
			{
				Config: testAccNewRelicCloudGcpIntegrationsConfigUpdated(GCPIntegrationTestConfig),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpIntegrationsExists(resourceName),
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

func testAccNewRelicCloudGcpIntegrationsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found %s", n)
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)
		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			return fmt.Errorf("An error occurred creating GCP integrations")
		}

		return nil
	}
}

func testAccNewRelicCloudGcpIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_gcp_integrations" && r.Type != "newrelic_cloud_gcp_link_account" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			return fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testSubAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("GCP Linked account is not unlinked: #{err}")
		}
	}
	return nil
}

func testAccNewRelicCloudGcpIntegrationsConfig(GCPIntegrationTestConfig map[string]string) string {
	return fmt.Sprintf(`
	provider "newrelic" {
  		account_id = "` + GCPIntegrationTestConfig["account_id"] + `"
  		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_gcp_link_account" "foo" {
        provider        = newrelic.cloud-integration-provider
  		name       = "` + GCPIntegrationTestConfig["name"] + `"
  		account_id = "` + GCPIntegrationTestConfig["account_id"] + `"
  		project_id = "` + GCPIntegrationTestConfig["project_id"] + `"
	}

resource "newrelic_cloud_gcp_integrations" "foo1" {
  provider        = newrelic.cloud-integration-provider
  account_id 		= "` + GCPIntegrationTestConfig["account_id"] + `"
  linked_account_id = newrelic_cloud_gcp_link_account.foo.id
  alloy_db {
    metrics_polling_interval = 400
  }
  app_engine {
    metrics_polling_interval = 400
  }
  big_query {
    metrics_polling_interval = 400
    fetch_tags               = true
  }
  big_table {
    metrics_polling_interval = 400
  }
  composer {
    metrics_polling_interval = 400
  }
  data_flow {
    metrics_polling_interval = 400
  }
  data_proc {
    metrics_polling_interval = 400
  }
  data_store {
    metrics_polling_interval = 400
  }
  fire_base_database {
    metrics_polling_interval = 400
  }
  fire_base_hosting {
    metrics_polling_interval = 400
  }
  fire_base_storage {
    metrics_polling_interval = 400
  }
  fire_store {
    metrics_polling_interval = 400
  }
  functions {
    metrics_polling_interval = 400
  }
  interconnect {
    metrics_polling_interval = 400
  }
  kubernetes {
    metrics_polling_interval = 400
  }
  load_balancing {
    metrics_polling_interval = 400
  }
  mem_cache {
    metrics_polling_interval = 400
  }
  pub_sub {
    metrics_polling_interval = 400
    fetch_tags               = true
  }
  redis {
    metrics_polling_interval = 400
  }
  router {
    metrics_polling_interval = 400
  }
  run {
    metrics_polling_interval = 400
  }
  spanner {
    metrics_polling_interval = 400
    fetch_tags               = true
  }
  sql {
    metrics_polling_interval = 400
  }
  storage {
    metrics_polling_interval = 400
    fetch_tags               = true
  }
  virtual_machines {
    metrics_polling_interval = 400
  }
  vpc_access {
    metrics_polling_interval = 400
  }
}`)
}

func testAccNewRelicCloudGcpIntegrationsConfigUpdated(GCPIntegrationTestConfig map[string]string) string {
	return fmt.Sprintf(`
	provider "newrelic" {
  		account_id = "` + GCPIntegrationTestConfig["account_id"] + `"
  		alias      = "cloud-integration-provider"
	}

	resource "newrelic_cloud_gcp_link_account" "foo" {
        provider        = newrelic.cloud-integration-provider
  		name       = "` + GCPIntegrationTestConfig["name"] + `"
  		account_id = "` + GCPIntegrationTestConfig["account_id"] + `"
  		project_id = "` + GCPIntegrationTestConfig["project_id"] + `"
	}

	resource "newrelic_cloud_gcp_integrations" "foo1" {
  		provider        = newrelic.cloud-integration-provider
  		account_id 		= "` + GCPIntegrationTestConfig["account_id"] + `"
  		linked_account_id = newrelic_cloud_gcp_link_account.foo.id
		  app_engine {
			metrics_polling_interval = 1400
		  }
		  big_query {
			metrics_polling_interval = 1400
			fetch_tags = true
		  }
		  big_table {
			metrics_polling_interval = 1400
		  }
		  composer{
			metrics_polling_interval = 1400
		  }
		  data_flow {
			metrics_polling_interval = 1400
		  }
		  data_proc{
			metrics_polling_interval = 1400
		  }
		  data_store{
			metrics_polling_interval = 1400
		  }
		  fire_base_database{
			metrics_polling_interval = 1400
		  }
		  fire_base_hosting{
			metrics_polling_interval = 1400
		  }
		  fire_base_storage{
			metrics_polling_interval = 1400
		  }
		  fire_store{
			metrics_polling_interval = 1400
		  }
		  functions{
			metrics_polling_interval = 1400
		  }
		  interconnect{
			metrics_polling_interval = 1400
		  }
		  kubernetes{
			metrics_polling_interval = 1400
		  }
		  load_balancing{
			metrics_polling_interval = 1400
		  }
		  mem_cache{
			metrics_polling_interval = 1400
		  }
		  pub_sub{
			metrics_polling_interval = 1400
			fetch_tags=true
		  }
		  redis{
			metrics_polling_interval = 1400
		  }
		  router{
			metrics_polling_interval = 1400
		  }
		  run{
			metrics_polling_interval = 1400
		  }
		  spanner{
			metrics_polling_interval = 1400
			fetch_tags=true
		  }
		  sql{
			metrics_polling_interval = 1400
		  }
		  storage{
			metrics_polling_interval = 1400
			fetch_tags=true
		  }
		  virtual_machines{
			metrics_polling_interval = 1400
		  }
		  vpc_access{
			metrics_polling_interval = 1400
		  }
	}
`)
}
