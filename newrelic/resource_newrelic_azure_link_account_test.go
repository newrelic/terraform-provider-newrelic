//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAzureLinkAccount_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccNewRelicAzureLinkAccountDestroy,
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

func testAccNewRelicAzureLinkAccountDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_azure_link_account" {
			continue
		}
		resourceId, err := strconv.Atoi(r.Primary.ID)
		if err != nil {
			fmt.Errorf("unable to convert string to int")
		}
		_, err = client.Cloud.GetLinkedAccount(testAccountID, resourceId)
		if err != nil {
			return err
		}
	}
	return nil
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
