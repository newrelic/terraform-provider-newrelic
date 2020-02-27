package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicApplicationLabel(t *testing.T) {
	// Skipping this test because we're awaiting an upstream fix on a deprecated feature.
	t.Skip("Skipping TestAccNewRelicApplicationLabel.")

	resourceName := "newrelic_application_label.foo"
	rCategory := acctest.RandString(10)
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicApplicationLabelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicApplicationLabelConfig(rCategory, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicApplicationLabelExists(resourceName),
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

// nolint:deadcode,unused
func testAccCheckNewRelicApplicationLabelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_application_label" {
			continue
		}
		key := r.Primary.ID

		_, err := client.APM.GetLabel(key)

		if err == nil {
			return fmt.Errorf("application label still exists")
		}
	}
	return nil
}

// nolint:deadcode,unused
func testAccNewRelicApplicationLabelConfig(category string, name string) string {
	return fmt.Sprintf(`
resource "newrelic_application_label" "foo" {
	category = "%s"
	name = "%s"
	links {
		applications = [215037795]
		servers = []
	}
}
`, category, name)
}

// nolint:deadcode,unused
func testAccCheckNewRelicApplicationLabelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no application label ID is set")
		}
		key := rs.Primary.ID
		id := strings.Split(key, ":")
		category := id[0]
		name := id[1]

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		label, err := client.APM.GetLabel(key)
		if err != nil {
			return err
		}

		if strings.EqualFold(label.Category, category) && !strings.EqualFold(label.Name, name) {
			return nil
		}

		return fmt.Errorf("application label not found: %v", key)
	}
}
