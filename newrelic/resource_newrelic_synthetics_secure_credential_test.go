//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

func TestAccNewRelicSyntheticsSecureCredential_Basic(t *testing.T) {
	resourceName := "newrelic_synthetics_secure_credential.foo"
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsSecureCredentialDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsSecureCredentialConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsSecureCredentialExists(resourceName),
				),
			},
			// Test: Update
			{
				Config: testAccNewRelicSyntheticsSecureCredentialConfigUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsSecureCredentialExists(resourceName),
				),
			},
			// Test: Import
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					// not returned from the API
					"value",
					"last_updated",
				},
			},
		},
	})
}

func TestAccNewRelicSyntheticsSecureCredential_Error(t *testing.T) {
	resourceName := "newrelic_synthetics_secure_credential.foo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicSyntheticsSecureCredentialDestroy,
		Steps: []resource.TestStep{
			// Test: Create
			{
				Config: testAccNewRelicSyntheticsSecureCredentialConfig("bad-key"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicSyntheticsSecureCredentialExists(resourceName),
				),
				ExpectError: regexp.MustCompile("Invalid key specified: key is blank, too long, or contains invalid characters."),
			},
		},
	})
}

func testAccCheckNewRelicSyntheticsSecureCredentialExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no synthetics secure credential key is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(60 * time.Second)

		nrqlQuery := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED' AND name = '%s' AND accountId = '%d'", rs.Primary.ID, testAccountID)
		found, err := client.Entities.GetEntitySearchByQuery(entities.EntitySearchOptions{}, nrqlQuery, []entities.EntitySearchSortCriteria{})
		if err != nil {
			return err
		}

		for _, e := range found.Results.Entities {
			if !strings.EqualFold(e.GetName(), rs.Primary.ID) {
				return fmt.Errorf("synthetics secure credential not found: %v - %v", rs.Primary.ID, found)
			}
		}

		return nil
	}
}

func testAccCheckNewRelicSyntheticsSecureCredentialDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*ProviderConfig).NewClient
	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_synthetics_secure_credential" {
			continue
		}

		// Unfortunately we still have to wait due to async delay with entity indexing :(
		time.Sleep(60 * time.Second)

		nrqlQuery := fmt.Sprintf("domain = 'SYNTH' AND type = 'SECURE_CRED' AND name = '%s' AND accountId = '%d'", r.Primary.ID, testAccountID)
		found, err := client.Entities.GetEntitySearchByQuery(entities.EntitySearchOptions{}, nrqlQuery, []entities.EntitySearchSortCriteria{})
		if err != nil {
			return err
		}

		if found.Count != 0 {
			return fmt.Errorf("synthetics secure credential still exists")
		}

	}
	return nil
}

func testAccNewRelicSyntheticsSecureCredentialConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_secure_credential" "foo" {
	key          = "%[1]s"
	value        = "Test Value"
	description  = "Test Description"
}
`, name)
}

func testAccNewRelicSyntheticsSecureCredentialConfigUpdated(name string) string {
	return fmt.Sprintf(`
resource "newrelic_synthetics_secure_credential" "foo" {
	key         = "%[1]s"
	value        = "Test Value Updated"
	description  = "Test Description"
}
`, name)
}
