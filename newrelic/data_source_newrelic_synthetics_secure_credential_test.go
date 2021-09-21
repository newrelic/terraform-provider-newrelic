//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsSecureCredentialDataSource_Basic(t *testing.T) {
	rand := acctest.RandString(5)
	rName := fmt.Sprintf("tf_test_%s", rand)
	resourceName := "data.newrelic_synthetics_secure_credential.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicSyntheticsSecureCredentialDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", strings.ToUpper(rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "Test description"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "last_updated"),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsSecureCredentialDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_secure_credential" "sc" {
	key = "%[1]s"
	value = "%[1]s"
	description = "Test description"
}

data "newrelic_synthetics_secure_credential" "foo" {
	key = newrelic_synthetics_secure_credential.sc.key
}
`, name)
}
