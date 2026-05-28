//go:build integration || CLOUD

package newrelic

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAwsEuSovereignIntegrations_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_aws_eu_sovereign_integrations.foo"
	rName := generateNameForIntegrationTestResource()

	testAwsEuSovereignAccountId := os.Getenv("NEW_RELIC_AWS_EU_SOVEREIGN_ACCOUNT_ID")
	if testAwsEuSovereignAccountId == "" {
		t.Skipf("NEW_RELIC_AWS_EU_SOVEREIGN_ACCOUNT_ID must be set for acceptance test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsEuSovereignIntegrationsDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicCloudAwsEuSovereignIntegrationsConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsEuSovereignIntegrationsExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicCloudAwsEuSovereignIntegrationsConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsEuSovereignIntegrationsExists(resourceName),
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

func testAccCheckNewRelicCloudAwsEuSovereignIntegrationsExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		linkedAccountID, err := strconv.Atoi(rs.Primary.Attributes["linked_account_id"])
		if err != nil {
			return fmt.Errorf("unable to parse linked account ID")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, linkedAccountID)
		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) == 0 {
			return fmt.Errorf("no integrations found")
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAwsEuSovereignIntegrationsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_aws_eu_sovereign_integrations" {
			continue
		}

		linkedAccountID, err := strconv.Atoi(r.Primary.Attributes["linked_account_id"])
		if err != nil {
			return fmt.Errorf("unable to parse linked account ID")
		}

		linkedAccount, err := client.Cloud.GetLinkedAccount(testAccountID, linkedAccountID)
		if err != nil {
			return err
		}

		if len(linkedAccount.Integrations) > 0 {
			return fmt.Errorf("integrations still exist")
		}
	}

	return nil
}

func testAccNewRelicCloudAwsEuSovereignIntegrationsConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_eu_sovereign_link_account" "foo" {
  name = "%[1]s"
  arn  = "%[2]s"
}

resource "newrelic_cloud_aws_eu_sovereign_integrations" "foo" {
  linked_account_id = newrelic_cloud_aws_eu_sovereign_link_account.foo.id

  billing {
    metrics_polling_interval = 3600
  }

  cloudtrail {
    metrics_polling_interval = 3600
    aws_regions              = ["eusc-de-east-1"]
  }

  x_ray {
    metrics_polling_interval = 3600
    aws_regions              = ["eusc-de-east-1"]
  }
}
`, rName, testAccExpectedAwsEuSovereignArn())
}

func testAccNewRelicCloudAwsEuSovereignIntegrationsConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_eu_sovereign_link_account" "foo" {
  name = "%[1]s"
  arn  = "%[2]s"
}

resource "newrelic_cloud_aws_eu_sovereign_integrations" "foo" {
  linked_account_id = newrelic_cloud_aws_eu_sovereign_link_account.foo.id

  billing {
    metrics_polling_interval = 7200
  }

  cloudtrail {
    metrics_polling_interval = 7200
    aws_regions              = ["eusc-de-east-1"]
  }

  x_ray {
    metrics_polling_interval = 7200
    aws_regions              = ["eusc-de-east-1"]
  }
}
`, rName, testAccExpectedAwsEuSovereignArn())
}
