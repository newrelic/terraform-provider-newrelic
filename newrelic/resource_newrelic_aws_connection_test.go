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

// TestAccNewRelicAwsConnection_Basic exercises create / read / update / import /
// destroy for the newrelic_aws_connection resource.
func TestAccNewRelicAwsConnection_Basic(t *testing.T) {
	resourceName := "newrelic_aws_connection.foo"
	rName := generateNameForIntegrationTestResource()

	roleArn := "arn:aws:iam::123456789012:role/tf-test-role"
	roleArnUpdated := "arn:aws:iam::123456789012:role/tf-test-role-rotated"
	externalID := "tf-test-external-id"
	externalIDUpdated := "tf-test-external-id-updated"
	orgID := "fb33fea3-4d7e-4736-9701-acb59a634fdf"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAwsConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAwsConnectionConfig(rName, "Acceptance test AWS connection", roleArn, externalID, true, "us-east-1", orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Acceptance test AWS connection"),
					resource.TestCheckResourceAttr(resourceName, "credential.0.assume_role.0.role_arn", roleArn),
					resource.TestCheckResourceAttr(resourceName, "credential.0.assume_role.0.external_id", externalID),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-east-1"),
					resource.TestCheckResourceAttr(resourceName, "scope_id", orgID),
					resource.TestCheckResourceAttr(resourceName, "scope_type", "ORGANIZATION"),
				),
			},
			// In-place update: change role_arn / assume-role external_id /
			// description / enabled / region. All are in EntityManagementAwsConnectionEntityUpdateInput.
			{
				Config: testAccNewRelicAwsConnectionConfig(rName, "Updated AWS connection", roleArnUpdated, externalIDUpdated, false, "us-west-2", orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAwsConnectionExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated AWS connection"),
					resource.TestCheckResourceAttr(resourceName, "credential.0.assume_role.0.role_arn", roleArnUpdated),
					resource.TestCheckResourceAttr(resourceName, "credential.0.assume_role.0.external_id", externalIDUpdated),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "region", "us-west-2"),
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

func testAccNewRelicAwsConnectionConfig(name, description, roleArn, externalID string, enabled bool, region, orgID string) string {
	return fmt.Sprintf(`
resource "newrelic_aws_connection" "foo" {
  name        = "%[1]s"
  description = "%[2]s"
  enabled     = %[3]t
  region      = "%[4]s"
  scope_id    = "%[6]s"
  scope_type  = "ORGANIZATION"

  credential {
    assume_role {
      role_arn    = "%[5]s"
      external_id = "%[7]s"
    }
  }
}
`, name, description, enabled, region, roleArn, orgID, externalID)
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
