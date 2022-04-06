//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"testing"
)

func TestAccNewRelicCloudGcpIntegrations(t *testing.T) {
	t.Skipf("skipping test until integrations work is finished")
	resourceName := "newrelic_cloud_gcp_link_account.foo"

	testGcpProjectID := os.Getenv("INTEGRATION_TESTING_GCP_PROJECT_ID")
	if testGcpProjectID == "" {
		t.Skipf("INTEGRATION_TESTING_GCP_PROJECT_ID must be set for acceptance test")
	}

	testGcpAccountName := os.Getenv("INTEGRATION_TESTING_GCP_ACCOUNT_NAME")
	if testGcpAccountName == "" {
		t.Skipf("INTEGRATION_TESTING_GCP_ACCOUNT_NAME must be set for acceptance test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudGcpIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudGcpIntegrationsConfig(testGcpAccountName, testGcpProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpIntegrationsExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicCloudGcpIntegrationsConfigUpdated(testGcpAccountName, testGcpProjectID),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpIntegrationsExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicCloudGcpIntegrationsExists(n string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			fmt.Errorf("not found %s", n)
		}
		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resourceId, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			fmt.Errorf("error converting string to int")
		}
		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)
		if err != nil && linkedAccount == nil {
			return err
		}
		return nil
	}
}

func testAccNewRelicCloudGcpIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_gcp_integrations" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if linkedAccount != nil && err == nil {
			return fmt.Errorf("Linked gcp account still exists: #{err}")
		}
	}
	return nil
}

func testAccNewRelicCloudGcpIntegrationsConfig(name string, projectId string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_gcp_link_account" "name" {
		  account_id=2520528
		  name= "%[1]s"
		  project_id = "%[2]s"
	}
	
	resource "newrelic_cloud_gcp_integrations" "foo1" {
		  account_id = 2520528
		  linked_account_id = newrelic_cloud_gcp_link_account.name.id
		  app_engine {
			metrics_polling_interval = 400
		  }
		  big_query {
			metrics_polling_interval = 400
			fetch_tags = true
		  }
		  big_table {
			metrics_polling_interval = 400
		  }
		  composer{
			metrics_polling_interval = 400
		  }
		  data_flow {
			metrics_polling_interval = 400
		  }
		  data_proc{
			metrics_polling_interval = 400
		  }
		  data_store{
			metrics_polling_interval = 400
		  }
		  fire_base_database{
			metrics_polling_interval = 400
		  }
		  fire_base_hosting{
			metrics_polling_interval = 400
		  }
		  fire_base_storage{
			metrics_polling_interval = 400
		  }
		  fire_store{
			metrics_polling_interval = 400
		  }
		  functions{
			metrics_polling_interval = 400
		  }
		  interconnect{
			metrics_polling_interval = 400
		  }
		  kubernetes{
			metrics_polling_interval = 400
		  }
		  load_balancing{
			metrics_polling_interval = 400
		  }
		  mem_cache{
			metrics_polling_interval = 400
		  }
		  pub_sub{
			metrics_polling_interval = 400
			fetch_tags=true
		  }
		  redis{
			metrics_polling_interval = 400
		  }
		  router{
			metrics_polling_interval = 400
		  }
		  run{
			metrics_polling_interval = 400
		  }
		  spanner{
			metrics_polling_interval = 400
			fetch_tags=true
		  }
		  sql{
			metrics_polling_interval = 400
		  }
		  storage{
			metrics_polling_interval = 400
			fetch_tags=true
		  }
		  virtual_machines{
			metrics_polling_interval = 400
		  }
		  vpc_access{
			metrics_polling_interval = 400
		  }
	}
	`, name, projectId)
}

func testAccNewRelicCloudGcpIntegrationsConfigUpdated(name string, projectId string) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_gcp_link_account" "name" {
		  account_id=2520528
		  name= "%[1]s-updated"
		  project_id = "%[2]s"
	}
	
	resource "newrelic_cloud_gcp_integrations" "foo1" {
		  account_id = 2520528
		  linked_account_id = newrelic_cloud_gcp_link_account.name.id
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
	`, name, projectId)
}
