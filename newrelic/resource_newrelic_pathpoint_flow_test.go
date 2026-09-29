//go:build integration || PATHPOINT

package newrelic

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pathpoint"
)

func TestAccNewRelicPathpointFlow_Basic(t *testing.T) {
	resourceName := "newrelic_pathpoint_flow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPathpointFlowDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicPathpointFlowConfigBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPathpointFlowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "account_id", strconv.Itoa(testAccountID)),
				),
			},
			// Test: Update name, description, category, and refresh interval
			{
				Config: testAccNewRelicPathpointFlowConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPathpointFlowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
					resource.TestCheckResourceAttr(resourceName, "category", "Testing"),
					resource.TestCheckResourceAttr(resourceName, "refresh_interval", "FIVE_MINUTES"),
					resource.TestCheckResourceAttr(resourceName, "health_rollup", "AUTOMATIC_ROLL_UP"),
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

func TestAccNewRelicPathpointFlow_WithStages(t *testing.T) {
	resourceName := "newrelic_pathpoint_flow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPathpointFlowDestroy,
		Steps: []resource.TestStep{
			// Test: Create with stages/levels/steps
			{
				Config: testAccNewRelicPathpointFlowConfigWithStages(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPathpointFlowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "stages.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "stages.0.name", "Stage One"),
					resource.TestCheckResourceAttr(resourceName, "stages.0.levels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "stages.0.levels.0.steps.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "stages.0.levels.0.steps.0.name", "Step One"),
					resource.TestCheckResourceAttr(resourceName, "stages.1.name", "Stage Two"),
				),
			},
			// Test: Update — add a step to the first stage's first level
			{
				Config: testAccNewRelicPathpointFlowConfigWithStagesUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPathpointFlowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "stages.0.levels.0.steps.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "stages.0.levels.0.steps.1.name", "Step Two"),
				),
			},
		},
	})
}

func TestAccNewRelicPathpointFlow_WithKPIs(t *testing.T) {
	resourceName := "newrelic_pathpoint_flow.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicPathpointFlowDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicPathpointFlowConfigWithKPIs(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicPathpointFlowExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "kpis.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "kpis.0.name", "Request Rate"),
					resource.TestCheckResourceAttr(resourceName, "kpis.0.query.0.from", "Transaction"),
					resource.TestCheckResourceAttr(resourceName, "kpis.0.query.0.select.0.aggregation_type", "COUNT"),
				),
			},
		},
	})
}

func testAccCheckNewRelicPathpointFlowExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no pathpoint flow ID is set")
		}

		accountIDStr := rs.Primary.Attributes["account_id"]
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account_id %q: %w", accountIDStr, err)
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		guid := pathpoint.EntityGUID(rs.Primary.ID)

		found, err := client.PathPoint.GetFlow(accountID, guid)
		if err != nil {
			return err
		}
		if found == nil || string(found.GUID) == "" {
			return fmt.Errorf("pathpoint flow not found: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckNewRelicPathpointFlowDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_pathpoint_flow" {
			continue
		}

		accountIDStr := r.Primary.Attributes["account_id"]
		accountID, err := strconv.Atoi(accountIDStr)
		if err != nil {
			return fmt.Errorf("invalid account_id %q: %w", accountIDStr, err)
		}

		guid := pathpoint.EntityGUID(r.Primary.ID)
		found, err := client.PathPoint.GetFlow(accountID, guid)
		if err == nil && found != nil && string(found.GUID) != "" {
			return fmt.Errorf("pathpoint flow still exists: %s", r.Primary.ID)
		}
	}

	return nil
}

func testAccNewRelicPathpointFlowConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "newrelic_pathpoint_flow" "foo" {
  account_id = %[1]d
  name       = %[2]q
}
`, testAccountID, name)
}

func testAccNewRelicPathpointFlowConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_pathpoint_flow" "foo" {
  account_id       = %[1]d
  name             = %[2]q
  description      = "updated description"
  category         = "Testing"
  refresh_interval = "FIVE_MINUTES"
  health_rollup    = "AUTOMATIC_ROLL_UP"
}
`, testAccountID, name+"-updated")
}

func testAccNewRelicPathpointFlowConfigWithStages(name string) string {
	return fmt.Sprintf(`
resource "newrelic_pathpoint_flow" "foo" {
  account_id = %[1]d
  name       = %[2]q

  stages {
    name = "Stage One"

    levels {
      steps {
        name = "Step One"
      }
    }
  }

  stages {
    name = "Stage Two"
  }
}
`, testAccountID, name)
}

func testAccNewRelicPathpointFlowConfigWithStagesUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_pathpoint_flow" "foo" {
  account_id = %[1]d
  name       = %[2]q

  stages {
    name = "Stage One"

    levels {
      steps {
        name = "Step One"
      }
      steps {
        name = "Step Two"
      }
    }
  }

  stages {
    name = "Stage Two"
  }
}
`, testAccountID, name)
}

func testAccNewRelicPathpointFlowConfigWithKPIs(name string) string {
	return fmt.Sprintf(`
resource "newrelic_pathpoint_flow" "foo" {
  account_id = %[1]d
  name       = %[2]q

  kpis {
    name = "Request Rate"

    query {
      from = "Transaction"

      select {
        aggregation_type = "COUNT"
      }
    }
  }
}
`, testAccountID, name)
}
