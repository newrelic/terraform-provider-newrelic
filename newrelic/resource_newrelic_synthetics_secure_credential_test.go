// +build integration

package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicSyntheticsSecureCredential_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_secure_credential.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsSecureCredentialDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsSecureCredentialConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsSecureCredentialExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsSecureCredentialConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsSecureCredentialExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"value",
				},
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsSecureCredentialExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics secure credential key is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		found, err := client.Synthetics.GetSecureCredential(rs.Primary.ID)
		if err != nil {
			return err
		}

		if !strings.EqualFold(found.Key, rs.Primary.ID) {
			return fmt.Errorf("synthetics secure credential not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsSecureCredentialDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_secure_credential" {
			continue
		}

		_, err := client.Synthetics.GetSecureCredential(r.Primary.ID)
		if err == nil {
			return fmt.Errorf("synthetics secure credential still exists")
		}

	}
	return nil
}

func testAccNewRelicSyntheticsSecureCredentialConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_secure_credential" "foo" {
	key          = "tf_test_%[1]s"
	value        = "Test Value"
	description  = "Test Description"
}
`, name)
}

func testAccNewRelicSyntheticsSecureCredentialConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_secure_credential" "foo" {
	key         = "tf_test_%[1]s"
	value        = "Test Value Updated"
	description  = "Test Description"
}
`, name)
}
