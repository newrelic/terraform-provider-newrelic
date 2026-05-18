//go:build integration || LOGGING_INTEGRATIONS

package newrelic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/newrelic/newrelic-client-go/v2/pkg/federatedlogs"
)

// TestAccNewRelicAwsConnection_Basic exercises create / read / import / destroy
// for the newrelic_aws_connection resource. Update is not supported by the
// underlying entity-management mutation, so no update step is performed.
//
// Required env vars (in addition to the standard testAccPreCheck set):
//   - TEST_AWS_ROLE_ARN: a valid IAM role ARN that the integration test account
//     can pass through to the EntityManagementCreateAwsConnection mutation.
func TestAccNewRelicAwsConnection_Basic(t *testing.T) {
	resourceName := "newrelic_aws_connection.foo"
	rName := generateNameForIntegrationTestResource()

	roleArn := "arn:aws:iam::123456789012:role/tf-test-role"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAwsConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAwsConnectionConfig(rName, roleArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "role_arn", roleArn),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east-1"),
					resource.TestCheckResourceAttrSet(resourceName, "scope_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNewRelicAwsConnectionConfig(name, roleArn string) string {
	return fmt.Sprintf(`
resource "newrelic_aws_connection" "foo" {
  name        = "%[1]s"
  description = "Acceptance test AWS connection"
  enabled     = true
  region      = "us-east-1"
  role_arn    = "%[2]s"
}
`, name, roleArn)
}

func testAccCheckNewRelicAwsConnectionExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no AWS connection ID set in state")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		resp, err := client.Federatedlogs.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching AWS connection %s: %w", rs.Primary.ID, err)
		}
		if resp == nil {
			return fmt.Errorf("AWS connection %s returned nil from API", rs.Primary.ID)
		}
		if _, ok := (*resp).(*federatedlogs.EntityManagementAwsConnectionEntity); !ok {
			return fmt.Errorf("entity %s is not an EntityManagementAwsConnectionEntity (got %T)", rs.Primary.ID, *resp)
		}
		return nil
	}
}

func testAccCheckNewRelicAwsConnectionDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_aws_connection" {
			continue
		}
		resp, err := client.Federatedlogs.GetEntityWithContext(context.Background(), rs.Primary.ID)
		if err != nil {
			return nil
		}
		if resp != nil {
			if _, ok := (*resp).(*federatedlogs.EntityManagementAwsConnectionEntity); ok {
				return fmt.Errorf("AWS connection %s still exists after destroy", rs.Primary.ID)
			}
		}
	}
	return nil
}
