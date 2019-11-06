package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"regexp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNewRelicAlertPolicy_Basic(t *testing.T) {
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists("newrelic_alert_policy.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_policy.foo", "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_policy.foo", "incident_preference", "PER_POLICY"),
				),
			},
			{
				Config: testAccCheckNewRelicAlertPolicyConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists("newrelic_alert_policy.foo"),
					resource.TestCheckResourceAttr(
						"newrelic_alert_policy.foo", "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(
						"newrelic_alert_policy.foo", "incident_preference", "PER_CONDITION"),
				),
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_import(t *testing.T) {
	resourceName := "newrelic_alert_policy.foo"
	rName := acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAlertPolicyConfig(rName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNewRelicAlertPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).Client
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_policy" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = client.GetAlertPolicy(int(id))

		if err == nil {
			return fmt.Errorf("policy still exists")
		}

	}
	return nil
}

func testAccCheckNewRelicAlertPolicyExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no policy ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).Client

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.GetAlertPolicy(int(id))
		if err != nil {
			return err
		}

		if strconv.Itoa(found.ID) != rs.Primary.ID {
			return fmt.Errorf("policy not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAccCheckNewRelicAlertPolicyConfig(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%s"
}
`, rName)
}

func testAccCheckNewRelicAlertPolicyConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name                = "tf-test-updated-%s"
  incident_preference = "PER_CONDITION"
}
`, rName)
}

func TestErrorThrownUponPolicyNameGreaterThan64Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	rName := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponPolicyNameGreaterThan64Char(rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testErrorThrownUponPolicyNameGreaterThan64Char(resourceName string) string {
	return fmt.Sprintf(`
provider "newrelic" {
  api_key = "foo"
}
resource "newrelic_alert_policy" "foo" {
  name = "really-long-name-that-is-more-than-sixtyfour-characters-long-tf-test-%[1]s"
}
`, resourceName, testAccExpectedApplicationName)
}

func TestErrorThrownUponPolicyNameLessThan1Char(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`expected length of name to be in the range \(1 \- 64\)`)
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testErrorThrownUponPolicyNameLessThan1Char(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testErrorThrownUponPolicyNameLessThan1Char() string {
	return `
provider "newrelic" {
  api_key = "foo"
}
resource "newrelic_alert_policy" "foo" {
  name = ""
}
`
}
