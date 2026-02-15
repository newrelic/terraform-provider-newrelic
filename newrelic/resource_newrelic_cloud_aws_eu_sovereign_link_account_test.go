//go:build integration || CLOUD

package newrelic

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicCloudAwsEuSovereignLinkAccount_Basic(t *testing.T) {
	resourceName := "newrelic_cloud_aws_eu_sovereign_link_account.foo"
	rName := generateNameForIntegrationTestResource()

	testAwsEuSovereignAccountId := os.Getenv("NEW_RELIC_AWS_EU_SOVEREIGN_ACCOUNT_ID")
	if testAwsEuSovereignAccountId == "" {
		t.Skipf("NEW_RELIC_AWS_EU_SOVEREIGN_ACCOUNT_ID must be set for acceptance test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicCloudAwsEuSovereignLinkAccountDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicCloudAwsEuSovereignLinkAccountConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsEuSovereignLinkAccountExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicCloudAwsEuSovereignLinkAccountConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicCloudAwsEuSovereignLinkAccountExists(resourceName),
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

func testAccCheckNewRelicCloudAwsEuSovereignLinkAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		linkedAccountID, parseErr := parseIDs(rs.Primary.ID, 1)
		if parseErr != nil {
			return fmt.Errorf("unable to parse linked account ID")
		}

		_, err := client.Cloud.GetLinkedAccount(testAccountID, linkedAccountID[0])
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckNewRelicCloudAwsEuSovereignLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_cloud_aws_eu_sovereign_link_account" {
			continue
		}

		linkedAccountID, parseErr := parseIDs(r.Primary.ID, 1)
		if parseErr != nil {
			return fmt.Errorf("unable to parse linked account ID")
		}

		_, err := client.Cloud.GetLinkedAccount(testAccountID, linkedAccountID[0])
		if err == nil {
			return fmt.Errorf("linked account still exists")
		}
	}

	return nil
}

func testAccNewRelicCloudAwsEuSovereignLinkAccountConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_eu_sovereign_link_account" "foo" {
  name = "%[1]s"
  arn  = "%[2]s"
}
`, rName, testAccExpectedAwsEuSovereignArn())
}

func testAccNewRelicCloudAwsEuSovereignLinkAccountConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_cloud_aws_eu_sovereign_link_account" "foo" {
  name = "%[1]s-updated"
  arn  = "%[2]s"
  metric_collection_mode = "PUSH"
}
`, rName, testAccExpectedAwsEuSovereignArn())
}

func testAccExpectedAwsEuSovereignArn() string {
	return fmt.Sprintf("arn:aws-eusc:iam::%s:role/NewRelicInfrastructure-Integrations", testAccExpectedAwsEuSovereignAccountId())
}

func testAccExpectedAwsEuSovereignAccountId() string {
	return os.Getenv("NEW_RELIC_AWS_EU_SOVEREIGN_ACCOUNT_ID")
}
