package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicAccountData_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataConfig_ByID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountData_ByName(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAccountDataConfig_ByName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountDataExists("data.newrelic_account.acc"),
				),
			},
		},
	})
}

func TestAccNewRelicAccountData_MissingAttributes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAccountDataConfig_MissingAttributes(),
				ExpectError: regexp.MustCompile("one of"),
			},
		},
	})
}

func TestAccNewRelicAccountData_ConflictingAttributes(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { testAccPreCheck(t) },
		Providers:  testAccProviders,
		IsUnitTest: true,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicAccountDataConfig_ConflictingAttributes(),
				ExpectError: regexp.MustCompile("exactly one of"),
			},
		},
	})
}

func testAccCheckNewRelicAccountDataExists(n string) resource.TestCheckFunc {
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

func testAccNewRelicAccountDataConfig_ByID() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	account_id = "%d"
}
`, testAccountID)
}

func testAccNewRelicAccountDataConfig_ByName() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	name = "%s"
}
`, testAccountName)
}

func testAccNewRelicAccountDataConfig_MissingAttributes() string {
	return `data "newrelic_account" "acc" {}`
}

func testAccNewRelicAccountDataConfig_ConflictingAttributes() string {
	return fmt.Sprintf(`
data "newrelic_account" "acc" {
	name = "%s"
	account_id = "%d"
}
`, testAccountName, testAccountID)
}
