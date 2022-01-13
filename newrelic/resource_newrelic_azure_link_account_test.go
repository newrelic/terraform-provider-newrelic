package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAzureLinkAccount_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		Steps: []resource.TestStep{
			//Test: Create
			{
				Config: testAccNewRelicAzureLinkAccountConfig(rName),
			},
			//Test: Update
			//TODO
			{
				Config: testAccNewRelicAzureLinkAccountConfigUpdated(rName),
			},
		},
	})
}

func testAccNewRelicAzureLinkAccountConfig(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_azure_link_account" "azure_account" {
		application_id = ""
		client_secret = ""
		name = "%[1]s"
		subscription_id = ""
		tenant_id = ""
	}
	`, name)
}

func testAccNewRelicAzureLinkAccountConfigUpdated(name string) string {
	return fmt.Sprintf(`
	resource "newrelic_azure_link_account" "azure_account" {
		application_id = ""
		client_secret = ""
		name = "%[1]s"
		subscription_id = ""
		tenant_id = ""
	}
	`, name)
}
