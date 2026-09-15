//go:build integration || DASHBOARDS

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNewRelicNotebookDataSource_MetadataOnly verifies that the data source
// returns title, organization_id, and blob_id when fetch_content = false, and
// that content is empty (no Blob Storage API call made).
func TestAccNewRelicNotebookDataSource_MetadataOnly(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-ds-notebook-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"
	dataSourceName := "data.newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotebookDataSourceConfig(rName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttr(resourceName, "title", rName),

					resource.TestCheckResourceAttrPair(dataSourceName, "guid", resourceName, "guid"),
					resource.TestCheckResourceAttr(dataSourceName, "title", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_id"),

					// content must be empty when fetch_content = false.
					resource.TestCheckResourceAttr(dataSourceName, "content", ""),
				),
			},
		},
	})
}

// TestAccNewRelicNotebookDataSource_WithContent verifies that fetch_content = true
// populates content with normalized JSON, and that blob_id on the data source
// matches the resource's blob_id.
func TestAccNewRelicNotebookDataSource_WithContent(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-ds-notebook-content-%s", acctest.RandString(5))
	resourceName := "newrelic_notebook.test"
	dataSourceName := "data.newrelic_notebook.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicNotebookDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNotebookDataSourceConfig(rName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "guid"),

					resource.TestCheckResourceAttrPair(dataSourceName, "guid", resourceName, "guid"),
					resource.TestCheckResourceAttr(dataSourceName, "title", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_id"),

					// content must be populated when fetch_content = true.
					resource.TestCheckResourceAttrSet(dataSourceName, "content"),

					// blob_id from data source must match the managed resource.
					resource.TestCheckResourceAttrPair(dataSourceName, "blob_id", resourceName, "blob_id"),
				),
			},
		},
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func testAccNotebookDataSourceConfig(name string, fetchContent bool) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title        = %[1]q
  content = jsonencode({
    type    = "declarative"
    version = 1
    content = [{
      type  = "container"
      props = { layout = "stack" }
      content = [{
        type  = "widget"
        content = {
          type = "visualization"
          id   = "viz.markdown"
          props = { text = "Data source acceptance test notebook." }
        }
      }]
    }]
  })
}

data "newrelic_notebook" "test" {
  guid          = newrelic_notebook.test.guid
  fetch_content = %[2]t
}
`, name, fetchContent)
}
