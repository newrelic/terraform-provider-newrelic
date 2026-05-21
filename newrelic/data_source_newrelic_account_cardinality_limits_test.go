//go:build integration || INGEST

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAccountCardinalityLimits_DataSource(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "newrelic_account_cardinality_limits" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.newrelic_account_cardinality_limits.test", "account_id"),
					resource.TestCheckResourceAttrSet("data.newrelic_account_cardinality_limits.test", "limits.#"),
				),
			},
		},
	})
}
