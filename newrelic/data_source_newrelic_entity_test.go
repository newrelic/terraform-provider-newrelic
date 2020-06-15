package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicEntityData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			// We need to give the entity search engine time to index the app
			time.Sleep(5 * time.Second)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity"),
				),
			},
		},
	})
}

func testAccCheckNewRelicEntityDataExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["guid"] == "" {
			return fmt.Errorf("expected to get an entity GUID")
		}

		if a["domain_id"] == "" {
			return fmt.Errorf("expected to get a domain ID")
		}

		if a["name"] != testAccExpectedApplicationName {
			return fmt.Errorf("expected the entity name to be: %s, but got: %s", testAccExpectedApplicationName, a["name"])
		}

		return nil
	}
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig() string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "APPLICATION"
	domain = "APM"
}
`, testAccExpectedApplicationName)
}
