package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	newrelic "github.com/paultyng/go-newrelic/api"
)

func TestAccNewRelicDashboard_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckNewRelicDashboardConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicDashboardExists("newrelic_dashboard.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_dashboard.foo", "title", fmt.Sprintf("tf-test-%s", rName)),
				),
			},
		},
	})
}

func testAccCheckNewRelicDashboardExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No dashboard ID is set")
		}

		client := testAccProvider.Meta().(*newrelic.Client)

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.GetDashboard(int(id))
		if err != nil {
			return err
		}

		if strconv.Itoa(found.ID) != rs.Primary.ID {
			return fmt.Errorf("Policy not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicDashboardConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_dashboard" "foo" {
  title = "tf-test-%s"
}
`, rName)
}
