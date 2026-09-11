//go:build integration || DASHBOARDS

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccNewRelicNotebookDataSource_MetadataOnly verifies that the data source
// returns title, organization_id and blob_id when fetch_content = false and
// does not make a Blob Storage API call (content is empty).
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
					// Resource fields populated.
					resource.TestCheckResourceAttrSet(resourceName, "guid"),
					resource.TestCheckResourceAttr(resourceName, "title", rName),

					// Data source matches the resource.
					resource.TestCheckResourceAttrPair(dataSourceName, "guid", resourceName, "guid"),
					resource.TestCheckResourceAttr(dataSourceName, "title", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_id"),

					// Content NOT fetched - should be empty.
					resource.TestCheckResourceAttr(dataSourceName, "content", ""),
				),
			},
		},
	})
}

// TestAccNewRelicNotebookDataSource_WithContent verifies that when
// fetch_content = true the data source populates the content attribute with
// the normalized notebook JSON returned by the Blob Storage API.
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

					// Data source attributes.
					resource.TestCheckResourceAttrPair(dataSourceName, "guid", resourceName, "guid"),
					resource.TestCheckResourceAttr(dataSourceName, "title", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "organization_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "blob_id"),

					// Content IS fetched - must be non-empty JSON.
					resource.TestCheckResourceAttrSet(dataSourceName, "content"),
				),
			},
		},
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func testAccNotebookDataSourceConfig(name string, fetchContent bool) string {
	return fmt.Sprintf(`
resource "newrelic_notebook" "test" {
  title = %[1]q
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
              props = { text = "Data source acceptance test notebook." }
            }
          }
        ]
      }
    ]
  })
}

data "newrelic_notebook" "test" {
  guid          = newrelic_notebook.test.guid
  fetch_content = %[2]t
}
`, name, fetchContent)
}
