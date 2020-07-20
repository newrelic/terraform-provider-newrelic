package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
				Config:      testAccNewRelicAccountDataSourceConfigMissingAttributes(),
				ExpectError: regexp.MustCompile("one of"),
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
				ExpectError: regexp.MustCompile("exactly one of"),
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
