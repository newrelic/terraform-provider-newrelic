package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicApplicationSettings_Basic(t *testing.T) {
	resourceName := "newrelic_application.app"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicApplicationConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicApplicationExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicApplicationConfigUpdated(testAccExpectedApplicationName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicApplicationExists(resourceName),
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

func testAccCheckNewRelicApplicationDestroy(s *terraform.State) error {
	// We expect the application to still exist
	return nil
}

// The test application for this data source is created in provider_test.go
func testAccNewRelicApplicationConfig() string {
	return fmt.Sprintf(`
resource "newrelic_application" "app" {
	name = "%s"
	app_apdex_threshold = "0.9"
	end_user_apdex_threshold = "0.8"
	enable_real_user_monitoring = true
}
`, testAccExpectedApplicationName)
}

func testAccNewRelicApplicationConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_application" "app" {
	name = "%s-updated"
	app_apdex_threshold = "0.8"
	end_user_apdex_threshold = "0.7"
	enable_real_user_monitoring = false
}
`, name)
}

func testAccCheckNewRelicApplicationExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no application ID is set")
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return nil
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		_, err = client.APM.GetApplication(int(id))
		if err != nil {
			return err
		}

		return nil
	}
}
