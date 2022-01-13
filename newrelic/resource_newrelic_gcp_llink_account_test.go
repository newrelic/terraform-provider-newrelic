package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicGcpLinkAccount_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicGcpLinkAccountConfig(rName),
			},
			//Test: Update
			//TODO
			{
				Config: testAccNewRelicGcpLinkAccountConfigUpdated(rName),
			},
		},
	})
}

func testAccNewRelicGcpLinkAccountConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_gcp_link_account" "gcp_account"{
		name = "%[1]s"
		project_id = ""
	}
	`, name)
}

func testAccNewRelicGcpLinkAccountConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_gcp_link_account" "gcp_account"{
		name = "%[1]s"
		project_id = ""
	}
	`, name)
}
