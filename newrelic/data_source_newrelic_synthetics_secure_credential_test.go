//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicSyntheticsSecureCredentialDataSource_Basic(t *testing.T) {
	rName := "INTEGRATION_TEST_SECURE_CREDENTIAL"
	resourceName := "data.newrelic_synthetics_secure_credential.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicSyntheticsSecureCredentialDataSourceConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "key", strings.ToUpper(rName)),
					resource.TestCheckResourceAttr(resourceName, "description", "TF Provider Acceptance Test Secure Cred"), /**/
					resource.TestCheckResourceAttrSet(resourceName, "last_updated"),
				),
			},
		},
	})
}

func testAccNewRelicSyntheticsSecureCredentialDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "newrelic_synthetics_secure_credential" "foo" {
	key = "%[1]s"
}
`, name)
}
