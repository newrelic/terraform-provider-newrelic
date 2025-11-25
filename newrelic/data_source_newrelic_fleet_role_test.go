//go:build integration || FLEET
// +build integration FLEET

package newrelic

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicFleetRole(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		resourceName string
		config       string
		checks       []resource.TestCheckFunc
	}{
		{
			name:         "Default",
			resourceName: "data.newrelic_fleet_role.default",
			config:       testAccDataSourceNewRelicFleetRoleDefault(),
			checks:       testAccCheckFleetRoleAttributes("data.newrelic_fleet_role.default", "Fleet Manager", "STANDARD"),
		},
		{
			name:         "ByName",
			resourceName: "data.newrelic_fleet_role.by_name",
			config:       testAccDataSourceNewRelicFleetRoleByName(),
			checks:       testAccCheckFleetRoleAttributes("data.newrelic_fleet_role.by_name", "Fleet Manager", "STANDARD"),
		},
		{
			name:         "ByType",
			resourceName: "data.newrelic_fleet_role.by_type",
			config:       testAccDataSourceNewRelicFleetRoleByType(),
			checks:       testAccCheckFleetRoleAttributes("data.newrelic_fleet_role.by_type", "Fleet Role x2 Custom", "CUSTOM"),
		},
		{
			name:         "ByNameAndType",
			resourceName: "data.newrelic_fleet_role.by_name_and_type",
			config:       testAccDataSourceNewRelicFleetRoleByNameAndType(),
			checks:       testAccCheckFleetRoleAttributes("data.newrelic_fleet_role.by_name_and_type", "Fleet Role x2 Custom", "CUSTOM"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:  func() { testAccPreCheckEnvVars(t) },
				Providers: testAccProviders,
				Steps: []resource.TestStep{
					{
						Config: tc.config,
						Check:  resource.ComposeTestCheckFunc(tc.checks...),
					},
				},
			})
		})
	}

	t.Run("ErrorCase", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			PreCheck:  func() { testAccPreCheckEnvVars(t) },
			Providers: testAccProviders,
			Steps: []resource.TestStep{
				{
					Config:      testAccDataSourceNewRelicFleetRoleError(),
					ExpectError: regexp.MustCompile(`no fleet role found with the given criteria`),
				},
			},
		})
	})
}

func testAccCheckFleetRoleAttributes(resourceName, name, roleType string) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		resource.TestCheckResourceAttr(resourceName, "name", name),
		resource.TestCheckResourceAttr(resourceName, "type", roleType),
		resource.TestCheckResourceAttrSet(resourceName, "id"),
	}
}

func testAccDataSourceNewRelicFleetRoleDefault() string {
	return `
		data "newrelic_fleet_role" "default" {}
`
}

func testAccDataSourceNewRelicFleetRoleByName() string {
	return `
		data "newrelic_fleet_role" "by_name" {
		  name = "Fleet Manager"
		}
`
}

func testAccDataSourceNewRelicFleetRoleByType() string {
	return `
		data "newrelic_fleet_role" "by_type" {
		  type = "CUSTOM"
		}
`
}

func testAccDataSourceNewRelicFleetRoleByNameAndType() string {
	return `
		data "newrelic_fleet_role" "by_name_and_type" {
		  	name = "Fleet Role x2 Custom"
			type = "CUSTOM"
		}
`
}

func testAccDataSourceNewRelicFleetRoleError() string {
	return `
		data "newrelic_fleet_role" "by_name_and_type" {
		  	name = "Fleet Role - Does Not Exist"
			type = "CUSTOM"
		}
`
}
