//go:build integration || AUTH
// +build integration AUTH

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var existingGroupName = "Integration Test Group 1 DO NOT DELETE"

func TestAccNewRelicGroupDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicGroupDataSourceConfiguration(authenticationDomainName, existingGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckGroupDataSourceAttributesExported(t, "data.newrelic_group.foo", existingGroupName),
				),
			},
		},
	})
}

func TestAccNewRelicGroupDataSource_MissingError(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicGroupDataSourceConfiguration(authenticationDomainName, fmt.Sprintf("%s-Invalid", existingGroupName)),
				ExpectError: regexp.MustCompile(`no group found with the specified parameters`),
			},
		},
	})
}

func testAccNewRelicCheckGroupDataSourceAttributesExported(t *testing.T, n string, groupName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("no group found")
		}

		return nil
	}
}

func testAccNewRelicGroupDataSourceConfiguration(authenticationDomainName string, groupName string) string {
	return fmt.Sprintf(`
	data "newrelic_authentication_domain" "foo" {
	  name = "%s"
	}
	
	data "newrelic_group" "foo" {
		authentication_domain_id = data.newrelic_authentication_domain.foo.id
		name 	 				 = "%s"
	}
`, authenticationDomainName, groupName)
}
