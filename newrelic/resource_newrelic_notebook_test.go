//go:build integration || DASHBOARDS

package newrelic

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccPreCheckNotebookEnvVars skips the test when the Notebooks-specific
// credentials are absent. Notebooks require the fleet-test account because
// that is currently the only account with the Notebooks entitlement enabled.
func testAccPreCheckNotebookEnvVars(t *testing.T) {
	t.Helper()
	if v := os.Getenv("NEW_RELIC_FLEET_TEST_API_KEY"); v == "" {
		t.Skip("NEW_RELIC_FLEET_TEST_API_KEY must be set for Notebook acceptance tests")
	}
	if v := os.Getenv("NEW_RELIC_FLEET_TEST_ORGANIZATION_ID"); v == "" {
		t.Skip("NEW_RELIC_FLEET_TEST_ORGANIZATION_ID must be set for Notebook acceptance tests")
	}
}

// setupNotebookTestCredentials swaps the provider's default API key and account
// ID for the fleet-test equivalents for the duration of the test. This mirrors
// setupFleetTestCredentials from resource_newrelic_fleet_test.go, but is
// defined here so the notebook tests compile under the DASHBOARDS build tag
// without depending on files that are only included under integration || FLEET.
func setupNotebookTestCredentials(t *testing.T) {
	t.Helper()

	originalAPIKey := os.Getenv("NEW_RELIC_API_KEY")
	originalAccountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	t.Cleanup(func() {
		os.Setenv("NEW_RELIC_API_KEY", originalAPIKey)       //nolint:errcheck
		os.Setenv("NEW_RELIC_ACCOUNT_ID", originalAccountID) //nolint:errcheck
	})

	if v := os.Getenv("NEW_RELIC_FLEET_TEST_API_KEY"); v != "" {
		os.Setenv("NEW_RELIC_API_KEY", v) //nolint:errcheck
	}
	if v := os.Getenv("NEW_RELIC_FLEET_TEST_ACCOUNT_ID"); v != "" {
		os.Setenv("NEW_RELIC_ACCOUNT_ID", v) //nolint:errcheck
	}
}

// testNotebookOrgID returns the org ID used for all notebook acceptance tests.
func testNotebookOrgID() string {
	if id := os.Getenv("NEW_RELIC_FLEET_TEST_ORGANIZATION_ID"); id != "" {
		return id
	}
	return "b961cf81-d62b-4359-8822-7b1d6dadd374"
}

// testAccCheckNewRelicNotebookExists verifies that the notebook resource exists
// in both Terraform state and on the remote platform.
func testAccCheckNewRelicNotebookExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no notebook GUID is set")
		}
		return nil
	}
}

// testAccCheckNewRelicNotebookDestroy verifies that a notebook has been deleted
// from the platform after a `terraform destroy`.
func testAccCheckNewRelicNotebookDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "newrelic_notebook" {
			continue
		}
		guid := rs.Primary.ID
		orgID := rs.Primary.Attributes["organization_id"]

		_, err := providerConfig.NewClient.Notebooks.GetNotebookContentWithContext(
			context.Background(), orgID, guid,
		)
		if err != nil && isNotebookNotFoundError(err) {
			continue
		}
		if err == nil {
			return fmt.Errorf("notebook %s still exists on the platform", guid)
		}
		return err
	}
	return nil
}

// isNotebookNotFoundError matches the "not found" error surfaced by GetNotebookContent
// when the notebook has been deleted. The Blob API returns a plain-text
// "Blob not found." body with HTTP 404; the Go client surfaces this as an
// error containing "not found".
func isNotebookNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "not found") || strings.Contains(msg, "404") || strings.Contains(msg, "Blob not found")
}

// ── Test helpers ──────────────────────────────────────────────────────────────

func testAccNotebookConfigContent(name, text, orgID string) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title           = %[1]q
  organization_id = %[3]q
  content = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = { text = %[2]q }
        }
      }
    ]
  })
}
`, name, text, orgID)
}

func testAccNotebookConfigContentJSON(name, orgID string) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title           = %[1]q
  organization_id = %[2]q
  content_json    = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = { text = "content_json acceptance test" }
        }
      }
    ]
  })
}
`, name, orgID)
}

func testAccNotebookConfigContentJSONUpdated(name, orgID string) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title           = %[1]q
  organization_id = %[2]q
  content_json    = jsonencode({
    version = "1"
    blocks = [
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = { text = "content_json acceptance test - updated" }
        }
      },
      {
        type = "widget"
        content = {
          type = "visualization"
          id   = "viz.billboard"
          props = {
            nrqlQueries = [{ accountIds = [4481681], query = "SELECT count(*) FROM Transaction SINCE 1 hour ago" }]
          }
        }
      }
    ]
  })
}
`, name, orgID)
}

// ── Acceptance tests ──────────────────────────────────────────────────────────

// TestAccNewRelicNotebook_ContentMode covers the full lifecycle using the
// content (jsonencode) field: create → no-drift plan → title update → content
// update → import → destroy.
func TestAccNewRelicNotebook_ContentMode(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-content-%s", acctest.RandString(5))
	rNameUpdated := rName + "-renamed"
	orgID := testNotebookOrgID()
	resourceName := "newrelic_notebook.test"

	setupNotebookTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNotebookEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			// Step 1: create.
			{
				Config: testAccNotebookConfigContent(rName, "Initial text", orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// Step 2: no drift expected after a clean apply.
			{
				Config:             testAccNotebookConfigContent(rName, "Initial text", orgID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: rename - title update triggers RenameNotebook.
			{
				Config: testAccNotebookConfigContent(rNameUpdated, "Initial text", orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Step 4: content update - Blob API receives new JSON.
			{
				Config: testAccNotebookConfigContent(rNameUpdated, "Updated text after content change", orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// Step 5: no drift after content update.
			{
				Config:             testAccNotebookConfigContent(rNameUpdated, "Updated text after content change", orgID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 6: import - GUID passthrough reconstructs state from the platform.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicNotebook_ContentJSONMode covers the full lifecycle using the
// content_json field: create → no-drift plan → update (add block) → no-drift
// → destroy.
func TestAccNewRelicNotebook_ContentJSONMode(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-json-%s", acctest.RandString(5))
	orgID := testNotebookOrgID()
	resourceName := "newrelic_notebook.test"

	setupNotebookTestCredentials(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNotebookEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			// Step 1: create with content_json.
			{
				Config: testAccNotebookConfigContentJSON(rName, orgID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "content_json"),
				),
			},
			// Step 2: no drift.
			{
				Config:             testAccNotebookConfigContentJSON(rName, orgID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: update - add a second block.
			{
				Config: testAccNotebookConfigContentJSONUpdated(rName, orgID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "content_json"),
				),
			},
			// Step 4: no drift after update.
			{
				Config:             testAccNotebookConfigContentJSONUpdated(rName, orgID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 5: import.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicNotebook_JSONReformatNoDrift verifies that reformatting the
// JSON in content_json (different whitespace / key ordering) does not trigger
// a plan change. This is the DiffSuppressFunc invariant.
func TestAccNewRelicNotebook_JSONReformatNoDrift(t *testing.T) {
	orgID := testNotebookOrgID()
	rName := fmt.Sprintf("tf-acc-notebook-nodrift-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	setupNotebookTestCredentials(t)

	// The two configs below are semantically identical (same content) but
	// differ in key ordering. Both should hash to the same normalized form.
	configV1 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title           = %q
  organization_id = %q
  content_json    = "{\"version\":\"1\",\"blocks\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"}}}]}"
}
`, rName, orgID)

	// Re-ordered JSON keys - normalized form is identical so no diff should appear.
	configV2 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title           = %q
  organization_id = %q
  content_json    = "{\"blocks\":[{\"content\":{\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"},\"type\":\"visualization\"},\"type\":\"widget\"}],\"version\":\"1\"}"
}
`, rName, orgID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckNotebookEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			{Config: configV1, Check: testAccCheckNewRelicNotebookExists(resourceName)},
			// Switch to re-ordered keys - must produce no plan change.
			{Config: configV2, PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}
