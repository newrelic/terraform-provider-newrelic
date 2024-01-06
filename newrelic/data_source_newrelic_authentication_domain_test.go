//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicAuthenticationDomain_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAuthenticationDomainDataSourceConfiguration("Test-Auth-Domain"),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckAuthenticationDomainExists(t, "data.newrelic_authentication_domain.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicAuthenticationDomain_MissingError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAuthenticationDomainDataSourceConfiguration("Invalid-Auth-Domain"),
				ExpectError: regexp.MustCompile(`no authentication domain found`),
			},
		},
	})
}

func testAccNewRelicCheckAuthenticationDomainExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get an ID of the matching authentication domain")
		}

		return nil
	}
}

func testAccNewRelicAuthenticationDomainDataSourceConfiguration(name string) string {
	return fmt.Sprintf(`
data "newrelic_authentication_domain" "foo" {
  name = "%s"
}
`, name)
}
