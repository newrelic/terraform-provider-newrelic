//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func TestAccNewRelicKeyTransaction_Basic(t *testing.T) {
	randomName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	resourceName := "newrelic_key_transaction.foo"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicKeyTransactionBasicConfiguration(fmt.Sprintf("%s", randomName)),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckKeyTransactionExists(resourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicKeyTransactionBasicConfiguration(fmt.Sprintf("%s-updated", randomName)),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckKeyTransactionExists(resourceName),
				),
			},
			// Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"domain",
					"type",
				},
			},
		},
	})
}

func TestAccNewRelicKeyTransaction_DuplicateNameError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// create a key transaction with a name that already exists in the UI
			// this is expected to throw an error, as only one key transaction may be created per metric name
			{
				Config:      testAccNewRelicKeyTransactionBasicConfiguration(fmt.Sprintf("%s", "terraform_acceptance_test_key_transaction_donot_delete")),
				ExpectError: regexp.MustCompile("\\s*"),
			},
		},
	})
}

func testAccNewRelicKeyTransactionBasicConfiguration(name string) string {
	return fmt.Sprintf(`
    resource "newrelic_key_transaction" "foo" {
        apdex_index 	     = 0.5
        application_guid     = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTUzNDQ4MjAy"
        browser_apdex_target = 0.5
        metric_name          = "test"
        name                 = "%[1]s"
    }
    `, name)
}

func testAccNewRelicCheckKeyTransactionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no key transaction ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		time.Sleep(5 * time.Second)
		found, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		x, foundOk := (*found).(*entities.KeyTransactionEntity)
		if !foundOk {
			return fmt.Errorf("no key transaction found")
		}
		if x.GUID != common.EntityGUID(rs.Primary.ID) {
			return fmt.Errorf("key transaction not found")
		}

		return nil
	}
}
