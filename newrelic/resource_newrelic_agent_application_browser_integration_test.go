//go:build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/agentapplications"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

func TestAccNewRelicAgentApplicationBrowser(t *testing.T) {
	resourceName := "newrelic_browser_application.foo"
	rName := generateNameForIntegrationTestResource()

	accountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skip(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicAgentApplicationBrowserConfig(
					accountID,
					rName,
					string(agentapplications.AgentApplicationBrowserLoaderTypes.SPA),
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAgentApplicationBrowserExists(resourceName),
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

func TestAccNewRelicAgentApplicationBrowser_InvalidLoaderType(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	accountID, err := testhelpers.GetTestAccountID()
	if err != nil {
		t.Skip(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsMonitorResourceDestroy,
		Steps: []resource.TestStep{
			// Create with invalid loader type. Expect an error.
			{
				Config: testAccNewRelicAgentApplicationBrowserConfig(
					accountID,
					rName,
					"INVALID_LOADER_TYPE",
				),
				ExpectError: regexp.MustCompile(`Expected type "AgentApplicationBrowserLoader"`),
			},
		},
	})
}

func testAccCheckNewRelicAgentApplicationBrowserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no browser agent application ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// Provide a minimal delay to allow for the entity to be indexed.
		time.Sleep(2 * time.Second)
		result, err := client.Entities.GetEntity(common.EntityGUID(rs.Primary.ID))
		if err != nil {
			return err
		}
		if result != nil {
			if string((*result).GetGUID()) != rs.Primary.ID {
				return fmt.Errorf("the browser agent application was not found %v - %v", (*result).GetGUID(), rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccNewRelicAgentApplicationBrowserConfig(accountID int, name string, loaderType string) string {
	return fmt.Sprintf(`
resource newrelic_browser_application foo {
 	account_id = %[1]d
	name = "%[2]s"
	cookies_enabled = true
	distributed_tracing_enabled = true
	loader_type = "%[3]s"
}
`, accountID, name, loaderType)
}
