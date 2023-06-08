//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicEntityData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig(testAccExpectedApplicationName, testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", testAccExpectedApplicationName, testAccountID),
				),
			},
		},
	})
}

// This test case checks if an entity with a single quote "'" in its name
// is created, and is subsequently, fetched by the data source.
func TestAccNewRelicSingleQuotedEntityData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccSingleQuotedPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig(testAccExpectedSingleQuotedApplicationName, testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", testAccExpectedSingleQuotedApplicationName, testAccountID),
				),
			},
		},
	})
}

func TestAccNewRelicEntityData_Missing(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicEntityDataConfig(strings.ToUpper(testAccExpectedApplicationName), testAccountID),
				ExpectError: regexp.MustCompile(`no entities found`),
			},
		},
	})
}

func TestAccNewRelicEntityData_IgnoreCase(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig_IgnoreCase(strings.ToUpper(testAccExpectedApplicationName), testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", testAccExpectedApplicationName, testAccountID),
				),
			},
		},
	})
}

func TestAccNewRelicEntityData_EntityInSubAccount(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicEntityDataConfig_EntityInSubAccount("Dummy App Two", 3957524),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicEntityDataExists(t, "data.newrelic_entity.entity", "Dummy App Two", 3957524),
				),
			},
		},
	})
}

func TestAccNewRelicEntityData_EntityAbsentInSubAccount(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicEntityDataConfig_EntityInSubAccount("Dummy App Two", 3814156),
				ExpectError: regexp.MustCompile(`no entities found`),
			},
		},
	})
}

func testAccCheckNewRelicEntityDataExists(t *testing.T, n string, appName string, accountID int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["account_id"] != strconv.Itoa(accountID) {
			return fmt.Errorf("account ID mismatch: the account ID specified in the configuration used to create this entity, and the account ID linked to the retrieved entity are not identical")
		}
		if a["guid"] == "" {
			return fmt.Errorf("expected to get an entity GUID")
		}

		if a["application_id"] == "" {
			return fmt.Errorf("expected to get an application ID")
		}

		if a["name"] != appName {
			return fmt.Errorf("expected the entity name to be: %s, but got: %s", appName, a["name"])
		}

		return nil
	}
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "application"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
	tag {
		key = "account"
		value = "New Relic Terraform Provider Acceptance Testing"
	}
}
`, name, accountId)
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig_IgnoreCase(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	ignore_case = true
	type = "application"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
}
`, name, accountId)
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig_InvalidType(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "app"
	domain = "apm"
	tag {
		key = "accountId"
		value = "%d"
	}
	tag {
		key = "account"
		value = "New Relic Terraform Provider Acceptance Testing"
	}
}
`, name, accountId)
}

// The test entity for this data source is created in provider_test.go
func testAccNewRelicEntityDataConfig_InvalidDomain(name string, accountId int) string {
	return fmt.Sprintf(`
data "newrelic_entity" "entity" {
	name = "%s"
	type = "application"
	domain = "VIZ"
	tag {
		key = "accountId"
		value = "%d"
	}
	tag {
		key = "account"
		value = "New Relic Terraform Provider Acceptance Testing"
	}
}
`, name, accountId)
}

// testAccNewRelicEntityDataConfig_EntityInSubAccount checks if the entity retrieved by applying
// the following configuration corresponds to the entity in the sub-account, and not the entity
// with an identical name in the main account (NEW_RELIC_ACCOUNT_ID).
func testAccNewRelicEntityDataConfig_EntityInSubAccount(name string, subAccountID int) string {
	return fmt.Sprintf(`
			provider "newrelic" {
  				account_id = %d
  				alias      = "entity-data-source-test-provider"
			}

			data "newrelic_entity" "entity" {
				provider = newrelic.entity-data-source-test-provider
				name = "%s"
				type = "APPLICATION"
				domain = "APM"
			}
`, subAccountID, name)
}
