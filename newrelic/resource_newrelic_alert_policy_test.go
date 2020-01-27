package newrelic

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccNewRelicAlertPolicy_Basic(t *testing.T) {
	resourceName := "newrelic_alert_policy.foo"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "incident_preference", "PER_POLICY"),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicAlertPolicyConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAlertPolicyExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("tf-test-updated-%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "incident_preference", "PER_CONDITION"),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_NoDiffOnReapply(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
			},
			{
				Config:             testAccNewRelicAlertPolicyConfig(rName),
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_ResourceNotFound(t *testing.T) {
	rName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAlertPolicyConfig(rName),
			},
			{
				PreConfig: testAccDeleteNewRelicAlertPolicy(rName),
				Config:    testAccNewRelicAlertPolicyConfig(rName),
			},
		},
	})
}

func TestAccNewRelicAlertPolicy_ErrorThrownWhenNameEmpty(t *testing.T) {
	expectedErrorMsg, _ := regexp.Compile(`name must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		IsUnitTest:   true,
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicAlertPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAlertPolicyConfigNameEmpty(),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccNewRelicAlertPolicyConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name = "tf-test-%s"
}
`, name)
}

func testAccNewRelicAlertPolicyConfigUpdated(rName string) string {
	return fmt.Sprintf(`
resource "newrelic_alert_policy" "foo" {
  name                = "tf-test-updated-%s"
  incident_preference = "PER_CONDITION"
}
`, rName)
}

func testAccCheckNewRelicAlertPolicyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_alert_policy" {
			continue
		}

		id, err := strconv.ParseInt(r.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		_, err = client.Alerts.GetPolicy(int(id))

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

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 32)
		if err != nil {
			return err
		}

		found, err := client.Alerts.GetPolicy(int(id))
		if err != nil {
			return err
		}

		if strconv.Itoa(found.ID) != rs.Primary.ID {
			return fmt.Errorf("policy not found: %v - %v", rs.Primary.ID, found)
		}

		return nil
	}
}

func testAlertPolicyConfigNameEmpty() string {
	return `
provider "newrelic" {
	api_key = "foo"
}

resource "newrelic_alert_policy" "foo" {
  name = ""
}
`
}
