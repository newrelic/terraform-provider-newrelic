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

func testAccNotebookConfigContent(name string) string {
	// Raw JSON string — exercises the string-passthrough, normalization,
	// and DiffSuppressFunc code paths directly.
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"content acceptance test\"}}}]}]}"
}
`, name)
}

func testAccNotebookConfigContentUpdated(name string) string {
	// Two widgets — verifies the content update path.
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"content acceptance test - updated\"}}},{\"type\":\"widget\",\"props\":{\"title\":\"\"},\"content\":{\"type\":\"visualization\",\"id\":\"viz.billboard\",\"props\":{\"nrqlQueries\":[{\"accountIds\":[0],\"query\":\"SELECT count(*) FROM Transaction SINCE 1 hour ago\"}]}}}]}]}"
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
  content = jsonencode({
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
                  { to = 1,           severity = "success"  },
                  { from = 1, to = 5, severity = "warning"  },
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

// TestAccNewRelicNotebook_ContentMode covers the full lifecycle using
// content with a raw JSON string.
func TestAccNewRelicNotebook_ContentMode(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-json-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			// Step 1: create.
			{
				Config: testAccNotebookConfigContent(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// Step 2: no drift.
			{
				Config:             testAccNotebookConfigContent(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: update — add a second widget.
			{
				Config: testAccNotebookConfigContentUpdated(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// Step 4: no drift after update.
			{
				Config:             testAccNotebookConfigContentUpdated(rName),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 5: import — Read populates content so pre- and
			// post-import state use the same field and the same normalized value.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// testAccNotebookConfigHeredocJSON returns a notebook whose content is written
// as a multi-line heredoc raw JSON string — the authoring style most users
// adopt when they don't want to use jsonencode. Exercises the heredoc JSON
// parse, normalize, and DiffSuppress paths.
func testAccNotebookConfigHeredocJSON(name string, accountID int) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title = %[1]q
  content = <<-JSON
    {
      "type": "declarative",
      "version": 1,
      "content": [
        {
          "type": "container",
          "props": { "layout": "stack" },
          "content": [
            {
              "type": "widget",
              "content": {
                "type": "visualization",
                "id": "viz.markdown",
                "props": { "text": "## Heredoc JSON test" }
              }
            },
            {
              "type": "widget",
              "props": { "title": "Request count" },
              "content": {
                "type": "visualization",
                "id": "viz.billboard",
                "props": {
                  "nrqlQueries": [
                    { "accountIds": [%[2]d], "query": "SELECT count(*) FROM Transaction SINCE 1 hour ago" }
                  ]
                }
              }
            }
          ]
        }
      ]
    }
  JSON
}
`, name, accountID)
}

// TestAccNewRelicNotebook_HeredocJSON verifies the full lifecycle using a
// multi-line heredoc raw JSON string — the most common real-world authoring
// style. Confirms that create, no-drift, and import all work correctly.
func TestAccNewRelicNotebook_HeredocJSON(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-heredoc-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotebookConfigHeredocJSON(rName, testAccountID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// No drift — normalized heredoc JSON must round-trip cleanly.
			{
				Config:             testAccNotebookConfigHeredocJSON(rName, testAccountID),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicNotebook_ReformatNoDrift verifies that reformatting
// content (key reordering, whitespace) does not trigger a plan change.
func TestAccNewRelicNotebook_ReformatNoDrift(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-nodrift-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	configV1 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %q
  content = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"}}}]}]}"
}
`, rName)

	// Same content, different key ordering — must produce no plan change.
	configV2 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %q
  content = "{\"content\":[{\"content\":[{\"content\":{\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"},\"type\":\"visualization\"},\"type\":\"widget\"}],\"props\":{\"layout\":\"stack\"},\"type\":\"container\"}],\"type\":\"declarative\",\"version\":1}"
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
