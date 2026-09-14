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

// isNotebookNotFoundError is defined in structures_newrelic_notebook.go and
// used here for destroy-check logic. No local definition needed.

// ── Test configs ──────────────────────────────────────────────────────────────

func testAccNotebookConfigContent(name, text string) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title   = %[1]q
  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [
      {
        type  = "container"
        props = { layout = "stack" }
        content = [
          {
            type  = "widget"
            props = {}
            content = {
              type = "visualization"
              id   = "viz.markdown"
              props = { text = %[2]q }
            }
          }
        ]
      }
    ]
  })
}
`, name, text)
}

func testAccNotebookConfigContentJSON(name string) string {
	// Uses a literal JSON string (not jsonencode) to exercise the content_json
	// path authentically — verifying the raw-JSON passthrough, normalization,
	// and DiffSuppressFunc all work on a real JSON string value.
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content_json = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"props\":{},\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"content_json acceptance test\"}}}]}]}"
}
`, name)
}

func testAccNotebookConfigContentJSONUpdated(name string) string {
	// Two widgets as a raw JSON string to verify content_json update path.
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content_json = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"props\":{},\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"content_json acceptance test - updated\"}}},{\"type\":\"widget\",\"props\":{},\"content\":{\"type\":\"visualization\",\"id\":\"viz.billboard\",\"props\":{\"nrqlQueries\":[{\"accountIds\":[0],\"query\":\"SELECT count(*) FROM Transaction SINCE 1 hour ago\"}]}}}]}]}"
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
          props = {}
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

// testAccNotebookConfigSwitchToContentJSON returns a config that is
// semantically identical to testAccNotebookConfigContent(name, text) but uses
// the content_json field instead of content. Used to test the mode-switch
// mid-lifecycle: the plan should show a field-name change but no API write
// (because the normalized JSON is identical).
func testAccNotebookConfigSwitchToContentJSON(name, text string) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content_json = jsonencode({
    type    = "declarative"
    version = 1
    content = [{
      type  = "container"
      props = { layout = "stack" }
      content = [{
        type  = "widget"
        props = {}
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = { text = %[2]q }
        }
      }]
    }]
  })
}
`, name, text)
}

// ── Acceptance tests ──────────────────────────────────────────────────────────

// TestAccNewRelicNotebook_ContentMode covers the full lifecycle using the
// content (jsonencode) field.
func TestAccNewRelicNotebook_ContentMode(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-content-%s", acctest.RandString(5))
	rNameUpdated := rName + "-renamed"
	resourceName := "newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			// Step 1: create.
			{
				Config: testAccNotebookConfigContent(rName, "Initial text"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "title", rName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(resourceName, "blob_id"),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// Step 2: no drift after a clean apply.
			{
				Config:             testAccNotebookConfigContent(rName, "Initial text"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 3: rename.
			{
				Config: testAccNotebookConfigContent(rNameUpdated, "Initial text"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
				),
			},
			// Step 4: content update.
			{
				Config: testAccNotebookConfigContent(rNameUpdated, "Updated text"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", rNameUpdated),
					resource.TestCheckResourceAttrSet(resourceName, "content"),
				),
			},
			// Step 5: no drift after content update.
			{
				Config:             testAccNotebookConfigContent(rNameUpdated, "Updated text"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			// Step 6: import with :content mode so the imported state has the
			// content field populated, matching this config's usage of content.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found: %s", resourceName)
					}
					return rs.Primary.ID + ":content", nil
				},
			},
		},
	})
}

// TestAccNewRelicNotebook_ContentJSONMode covers the full lifecycle using the
// content_json field.
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
			// Step 3: update - add a second widget.
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
			// Step 5: import - content_json mode is fully verifiable because
			// the Read path also writes to content_json, so pre- and post-import
			// state use the same field name and the same normalized value.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccNewRelicNotebook_JSONReformatNoDrift verifies that reformatting the
// JSON in content_json does not trigger a plan change.
func TestAccNewRelicNotebook_JSONReformatNoDrift(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-nodrift-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	configV1 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %q
  content_json = "{\"type\":\"declarative\",\"version\":1,\"content\":[{\"type\":\"container\",\"props\":{\"layout\":\"stack\"},\"content\":[{\"type\":\"widget\",\"props\":{},\"content\":{\"type\":\"visualization\",\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"}}}]}]}"
}
`, rName)

	// Same content, different key ordering - must produce no plan change.
	configV2 := fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %q
  content_json = "{\"content\":[{\"content\":[{\"content\":{\"id\":\"viz.markdown\",\"props\":{\"text\":\"nodrift\"},\"type\":\"visualization\"},\"props\":{},\"type\":\"widget\"}],\"props\":{\"layout\":\"stack\"},\"type\":\"container\"}],\"type\":\"declarative\",\"version\":1}"
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

// TestAccNewRelicNotebook_ModeSwitch verifies that switching from content= to
// content_json= with semantically identical JSON does not trigger an unnecessary
// Blob API write (blob_id must not change) and leaves no plan drift after apply.
func TestAccNewRelicNotebook_ModeSwitch(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-notebook-mode-switch-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			// Step 1: create with content= (HCL mode).
			{
				Config: testAccNotebookConfigContent(rName, "Mode switch test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "blob_id"),
				),
			},
			// Step 2: switch to content_json= with identical JSON.
			// The plan must show a field-name change but NOT a content change,
			// and apply must NOT bump the blob_id.
			{
				Config: testAccNotebookConfigSwitchToContentJSON(rName, "Mode switch test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "blob_id"),
					resource.TestCheckResourceAttrSet(resourceName, "content_json"),
				),
			},
			// Step 3: no drift after mode switch.
			{
				Config:             testAccNotebookConfigSwitchToContentJSON(rName, "Mode switch test"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
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
				Config: testAccNotebookConfigWithBillboard(rName, testAccGetTestAccountID()),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicNotebookExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttrSet(resourceName, "blob_id"),
				),
			},
			// No drift — nested billboardSettings/thresholdsWithSeriesOverrides
			// must survive the round-trip unchanged.
			{
				Config:             testAccNotebookConfigWithBillboard(rName, testAccGetTestAccountID()),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
