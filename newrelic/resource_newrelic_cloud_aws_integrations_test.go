//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAwsIntegrations_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_aws_integrations.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsIntegrationsDestroy,
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAwsIntegrationsConfig(109698),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsIntegrationsExist(resourceName),
				),
			},
			//Test: Update
			{
				Config: testAccNewRelicAwsIntegrationsConfigUpdated(109698),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsIntegrationsExist(resourceName),
				),
			},
		},
	})
}

func testAccCheckNewRelicCloudAwsIntegrationsExist(n string) resource.TestCheckFunc {
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
			fmt.Errorf("An error occurred creating AWS integrations")
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAwsIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_aws_integrations" {
			continue
		}

		resourceId, err := strconv.Atoi(r.Primary.ID)

		if err != nil {
			fmt.Errorf("error converting string id to int")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, resourceId)

		if len(linkedAccount.Integrations) != 0 && err == nil {
			return fmt.Errorf("AWS integrations were not unlinked: #{err}")
		}
	}

	return nil
}

func testAccNewRelicAwsIntegrationsConfig(linkedAccountID int) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_integrations" "foo" {
		linked_account_id = %[1]d

		billing {
			metrics_polling_interval = 6000
		}

		cloudtrail {
			metrics_polling_interval = 6000
		}

		health {
			metrics_polling_interval = 6000
		}

		trusted_advisor {
			metrics_polling_interval = 6000
		}

		x_ray {
			metrics_polling_interval = 6000
		}
	}
`, linkedAccountID)
}

func testAccNewRelicAwsIntegrationsConfigUpdated(linkedAccountID int) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_integrations" "foo" {
		linked_account_id = %[1]d
		billing {
			metrics_polling_interval = 6000
		}
		cloudtrail {
			metrics_polling_interval = 6000
		}
		health {
			metrics_polling_interval = 6000
		}
		trusted_advisor {
			metrics_polling_interval = 6000
		}
	}
`, linkedAccountID)
}
