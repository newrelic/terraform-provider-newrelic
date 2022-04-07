//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudGcpIntegrations_Basic(t *testing.T) {
	t.Skipf("Skipping test until enviroment variables are added")
	resourceName := "newrelic_cloud_gcp_integrations.foo"
	testGcpLinkedAccountID := os.Getenv("INTEGRATION_TESTING_GCP_LINKED_ACCOUNT_ID")
	if testGcpLinkedAccountID == "" {
		t.Skipf("INTEGRATION_TESTING_GCP_LINKED_ACCOUNT_IDINTEGRATION_TESTING_GCP_LINKED_ACCOUNT_ID must be set for acceptance test")
	}
	linkedAccountID, convErr := strconv.Atoi(testGcpLinkedAccountID)

	if convErr != nil {
		fmt.Errorf("error converting string to int")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicCloudGcpIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicCloudGcpIntegrationsConfig(linkedAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpIntegrationsExists(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicCloudGcpIntegrationsConfigUpdated(linkedAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCloudGcpIntegrationsExists(resourceName),
				),
			},
		},
	})
}

func testAccNewRelicCloudGcpIntegrationsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			fmt.Errorf("not found %s", n)
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resourceId, err := strconv.Atoi(rs.Primary.ID)

		if err != nil {
			fmt.Errorf("error converting string to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)
		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			fmt.Errorf("An error occurred creating GCP integrations")
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

		if (linkedAccount.Integrations) != nil && err == nil {
			return fmt.Errorf("GCP integrations were not unlinked: #{err}")
		}
	}
	return nil
}

func testAccNewRelicCloudGcpIntegrationsConfig(linkedAccountId int) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_gcp_integrations" "foo" {
		  account_id = 2520528
		  linked_account_id = 111492
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
	`, linkedAccountId)
}

func testAccNewRelicCloudGcpIntegrationsConfigUpdated(linkedAccountId int) string {
	return fmt.Sprintf(`
	resource "newrelic_cloud_gcp_integrations" "foo" {
		  account_id = 2520528
		  linked_account_id = 111492
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
	`, linkedAccountId)
}
