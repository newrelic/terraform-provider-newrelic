package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccExpectedKeyTransactionName = "get /"
)

func TestAccNewRelicKeyTransactionDataSource_Basic(t *testing.T) {
	if !nrInternalAccount {
		t.Skipf("New Relic internal testing account required")
	}

	resourceName := "data.newrelic_key_transaction.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicKeyTransactionDataSourceConfig(testAccExpectedKeyTransactionName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicKeyTransactionDataSourceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", testAccExpectedKeyTransactionName),
				),
			},
		},
	})
}

func testAccNewRelicKeyTransactionDataSourceConfig(rName string) string {
	return fmt.Sprintf(`
data "newrelic_key_transaction" "foo" {
	name = "%s"
}
`, testAccExpectedKeyTransactionName)
}

func testAccCheckNewRelicKeyTransactionDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get a key transaction from New Relic")
		}

		if !strings.EqualFold(testAccExpectedKeyTransactionName, a["name"]) {
			return fmt.Errorf("expected the key transaction name to be: %s, but got: %s", testAccExpectedKeyTransactionName, a["name"])
		}

		return nil
	}
}
