//go:build integration || DASHBOARDS

package newrelic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccCheckNewRelicNotebookExists verifies the notebook is in Terraform state.
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

// testAccCheckNewRelicNotebookDestroy verifies the notebook has been removed
// from the platform after a terraform destroy.
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

// ── Test configs ──────────────────────────────────────────────────────────────

func testAccNotebookConfigContentJSON(name string) string {
	// Raw JSON string — exercises the string-passthrough, normalization,
	// and DiffSuppressFunc code paths directly.
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content_json = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"content_json acceptance test\"}}}]}]}"
}
`, name)
}

func testAccNotebookConfigContentJSONUpdated(name string) string {
	// Two widgets — verifies the content_json update path.
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content_json = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"content_json acceptance test - updated\"}}},{\"type\":\"widget\",\"props\":{\"title\":\"\"},\"content\":{\"type\":\"visualization\",\"id\":\"viz.billboard\",\"props\":{\"nrqlQueries\":[{\"accountIds\":[0],\"query\":\"SELECT count(*) FROM Transaction SINCE 1 hour ago\"}]}}}]}]}"
}
`, name)
}

// testAccNotebookConfigWithBillboard produces a two-widget notebook containing
// a markdown header and a billboard with threshold configuration. Used to verify
// that nested billboardSettings and thresholdsWithSeriesOverrides props survive
// the full create → read → no-drift cycle.
func testAccNotebookConfigWithBillboard(name string, accountID int) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title = %[1]q
  content_json = jsonencode({
    type    = "declarative"
    version = 1
    content = [{
      type  = "container"
      props = { layout = "stack" }
      content = [
        {
          type  = "widget"
          content = {
            type = "visualization"
            id   = "viz.markdown"
            props = { text = "## Notebook with billboard" }
          }
        },
        {
          type  = "widget"
          props = { title = "Error rate" }
          content = {
            type = "visualization"
            id   = "viz.billboard"
            props = {
              nrqlQueries = [{
                accountIds = [%[2]d]
                query      = "SELECT percentage(count(*), WHERE error IS true) AS 'Error %%' FROM Transaction SINCE 1 hour ago"
              }]
              thresholdsWithSeriesOverrides = {
                thresholds = [
                  { to = 1,           severity = "success"  }
                  { from = 1, to = 5, severity = "warning"  }
                  { from = 5,         severity = "critical" }
                ]
              }
              billboardSettings = {
                visual = { alignment = "inline", display = "auto" }
              }
            }
          }
        }
      ]
    }]
  })
}
`, name, accountID)
}

// ── Acceptance tests ──────────────────────────────────────────────────────────

// TestAccNewRelicNotebook_ContentJSONMode covers the full lifecycle using
// content_json with a raw JSON string.
func TestAccNewRelicNotebook_ContentJSONMode(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-json-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			// Step 1: create.
			{
				Config: testAccNotebookConfigContentJSON(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "content_json"),
				),
			},
			// Step 2: no drift.
			{
				Config:             testAccNotebookConfigContentJSON(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: update — add a second widget.
			{
				Config: testAccNotebookConfigContentJSONUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "content_json"),
				),
			},
			// Step 4: no drift after update.
			{
				Config:             testAccNotebookConfigContentJSONUpdated(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 5: import — Read populates content_json so pre- and
			// post-import state use the same field and the same normalized value.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicNotebook_JSONReformatNoDrift verifies that reformatting
// content_json (key reordering, whitespace) does not trigger a plan change.
func TestAccNewRelicNotebook_JSONReformatNoDrift(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-nodrift-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	configV1 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %q
  content_json = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"}}}]}]}"
}
`, rName)

	// Same content, different key ordering — must produce no plan change.
	configV2 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %q
  content_json = "{\"content\":[{\"content\":[{\"content\":{\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"},\"type\":\"visualization\"},\"type\":\"widget\"}],\"props\":{\"layout\":\"stack\"},\"type\":\"container\"}],\"type\":\"declarative\",\"version\":1}"
}
`, rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			{Config: configV1, Check: testAccCheckNewRelicNotebookExists(resourceName)},
			{Config: configV2, PlanOnly: true, ExpectNonEmptyPlan: false},
		},
	})
}

// TestAccNewRelicNotebook_WithBillboard verifies that a notebook containing
// billboardSettings, thresholdsWithSeriesOverrides, and nested viz props
// round-trips through create → read → no-drift without any spurious diff.
func TestAccNewRelicNotebook_WithBillboard(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-billboard-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotebookConfigWithBillboard(rName, testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "blob_id"),
				),
			},
			// No drift — nested billboardSettings/thresholdsWithSeriesOverrides
			// must survive the round-trip unchanged.
			{
				Config:             testAccNotebookConfigWithBillboard(rName, testAccountID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
