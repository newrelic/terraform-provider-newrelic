//go:build integration || AUTH

package newrelic

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var authenticationDomainName = "Test-Auth-Domain DO NOT DELETE"
var existingUserEmail = strings.ReplaceAll(userEmailPrefix, "#", "integration")
var existingUserName = "Integration Test User 1 DO NOT DELETE"

func TestAccNewRelicUserDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicUserDataSourceConfiguration(authenticationDomainName, existingUserEmail, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckUserDataSourceExists(t, "data.newrelic_user.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicUserDataSource_EmailAndName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicUserDataSourceConfiguration(authenticationDomainName, existingUserEmail, existingUserName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckUserDataSourceExists(t, "data.newrelic_user.foo"),
				),
			},
		},
	})
}

func TestAccNewRelicUserDataSource_MissingError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicUserDataSourceConfiguration(authenticationDomainName, "", fmt.Sprintf("%s-Invalid", existingUserName)),
				ExpectError: regexp.MustCompile(`no user found with the specified parameters`),
			},
		},
	})
}

func testAccNewRelicCheckUserDataSourceExists(t *testing.T, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("expected to get an ID of the matching user")
		}

		return nil
	}
}

func testAccNewRelicUserDataSourceConfiguration(authenticationDomainName string, userEmailID string, userName string) string {
	return fmt.Sprintf(`
data "newrelic_authentication_domain" "foo" {
  name = "%s"
}

data "newrelic_user" "foo" {
	authentication_domain_id = data.newrelic_authentication_domain.foo.id
	email_id = "%s"
	name 	 = "%s"
}
`, authenticationDomainName, userEmailID, userName)
}
