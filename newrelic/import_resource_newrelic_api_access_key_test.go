// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicAPIAccessKey_importBasic(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))
	_, accountID := retrieveIdsFromEnvOrSkip(t, "NEW_RELIC_TEST_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngest(accountID, keyTypeIngestLicense, keyName, keyNotes),
			},
			{
				ResourceName:      "newrelic_api_access_key.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNewRelicAPIAccessKeyImportStateIdFunc_Basic("newrelic_api_access_key.foobar"),
			},
		},
	})
}

func testAccNewRelicAPIAccessKeyImportStateIdFunc_Basic(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["id"], rs.Primary.Attributes["key_type"]), nil
	}
}
