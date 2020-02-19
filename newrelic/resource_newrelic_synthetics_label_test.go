package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicSyntheticsLabel(t *testing.T) {
	resourceName := "newrelic_synthetics_label.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsLabelDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsLabelConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsLabelExists(resourceName),
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

func testAccCheckNewRelicSyntheticsLabelDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_label" {
			continue
		}

		ids := strings.Split(r.Primary.ID, ":")
		monitorID := ids[0]
		labelType := ids[1]
		value := ids[2]

		labels, err := client.Synthetics.GetMonitorLabels(monitorID)

		if err != nil {
			return err
		}

		for _, l := range labels {
			if l.Type == labelType && l.Value == value {
				return fmt.Errorf("synthetics label still exists")
			}
		}
	}
	return nil
}

func testAccNewRelicSyntheticsLabelConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_monitor" "foo" {
	name = "%[1]s"
	type = "SIMPLE"
	frequency = 5
	status = "ENABLED"
	locations = ["AWS_US_EAST_1"]
	uri = "https://example.com"
}

resource "newrelic_synthetics_label" "foo" {
	monitor_id = newrelic_synthetics_monitor.foo.id
	type  = "testType"
	value = "testValue"
}
`, name)
}

func testAccCheckNewRelicSyntheticsLabelExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics label ID is set")
		}

		ids := strings.Split(rs.Primary.ID, ":")
		monitorID := ids[0]
		labelType := ids[1]
		value := ids[2]

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		labels, err := client.Synthetics.GetMonitorLabels(monitorID)
		if err != nil {
			return err
		}

		for _, l := range labels {
			if !strings.EqualFold(l.Type, labelType) || !strings.EqualFold(l.Value, value) {
				continue
			}

			return nil
		}

		return fmt.Errorf("synthetics label not found: %v", rs.Primary.ID)
	}
}
