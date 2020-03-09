package newrelic

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicApplication_Basic(t *testing.T) {
	resourceName := "newrelic_application.app"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicApplicationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicApplicationConfig(),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckNewRelicApplicationExists2(testAccExpectedApplicationName),
					resource.TestCheckResourceAttr(
						resourceName, "name", testAccExpectedApplicationName),
					resource.TestCheckResourceAttr(
						resourceName, "app_apdex_threshold", "0.9"),
					resource.TestCheckResourceAttr(
						resourceName, "end_user_apdex_threshold", "0.8"),
					// resource.TestCheckResourceAttr(
					// 	resourceName, "enable_real_user_monitoring", "true"),
				),
			},
			// Test: Import
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
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

func testAccCheckNewRelicApplicationExists2(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no application ID is set")
		}

		key := rs.Primary.ID
		log.Printf("\n\n\n[ZACH] Key %s", key)

		// id := strings.Split(key, ":")
		// category := id[0]
		// name := id[1]
		//
		// client := testAccProvider.Meta().(*ProviderConfig).NewClient
		//
		// app, err := client.APM.GetApplication(key)
		// if err != nil {
		// 	return err
		// }
		//
		// if strings.EqualFold(app.Name, category) && !strings.EqualFold(label.Name, name) {
		// 	return nil
		// }
		//
		// return fmt.Errorf("application label not found: %v", key)
		return nil
	}
}
