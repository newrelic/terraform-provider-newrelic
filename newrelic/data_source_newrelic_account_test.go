//go:build integration || AUTH
// +build integration AUTH

package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicAccountDataSource_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataSourceConfigByID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataSourceExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountDataSource_ByName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataSourceExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountDataSource_MissingAttributes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataSourceConfigMissingAttributes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataSourceExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountDataSource_ConflictingAttributes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAccountDataSourceConfigConflictingAttributes(),
				ExpectError: regexp.MustCompile("conflicts with"),
			},
		},
	})
}

func testAccCheckNewRelicAccountDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		id := r.Primary.ID
		a := r.Primary.Attributes

		if id == "" {
			return fmt.Errorf("expected to get an account from New Relic")
		}

		if a["name"] == "" {
			return fmt.Errorf("expected to get an account name from New Relic")
		}

		if a["account_id"] == "" {
			return fmt.Errorf("expected to get an account ID from New Relic")
		}

		if a["region"] == "" {
			return fmt.Errorf("expected to get an region code from New Relic")
		}

		if a["name"] != testAccountName {
			return fmt.Errorf("expected account name to be %s, got %s", testAccountName, a["name"])
		}

		if a["account_id"] != strconv.Itoa(testAccountID) {
			return fmt.Errorf("expected account ID to be %d, got %s", testAccountID, a["account_id"])
		}

		if a["region"] != "us01" {
			return fmt.Errorf("expected region code to be us01, got %s", a["region"])
		}

		return nil
	}
}

func testAccNewRelicAccountDataSourceConfigByID() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	account_id = "%d"
}
`, testAccountID)
}

func testAccNewRelicAccountDataSourceConfigByName() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	name = "%s"
}
`, testAccountName)
}

func testAccNewRelicAccountDataSourceConfigMissingAttributes() string {
	return `data "newrelic_account" "acc" {}`
}

func testAccNewRelicAccountDataSourceConfigConflictingAttributes() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	name = "%s"
	account_id = "%d"
}
`, testAccountName, testAccountID)
}
