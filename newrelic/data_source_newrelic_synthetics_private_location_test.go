//go:build integration || SYNTHETICS
// +build integration SYNTHETICS

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

func TestAccNewRelicSyntheticsPrivateLocationDataSource_Basic(t *testing.T) {
	t.Parallel()

	privateLocationName := fmt.Sprintf("tf-test-%s", acctest.RandString(5))
	client, err := newIntegrationTestClient()
	if err != nil {
		t.Skipf("Skipping test due to error instantiating an integration test client: %s", err)
	}

	var privateLocationGUID synthetics.EntityGUID
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

			result, err := client.Synthetics.SyntheticsCreatePrivateLocation(
				testAccountID,
				"created via TF integration tests",
				privateLocationName,
				false,
			)

			if err != nil {
				t.Skipf("Skipping test due to error creating test sythentics private location: %s", err)
			}

			privateLocationGUID = result.GUID

			// Workaround for async entity creation so we can test the data source below
			time.Sleep(60 * time.Second)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigDataSourceSyntheticsLocation(privateLocationName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_private_location.bar", "name"),
				),
			},
		},
	})

	// Clean up extra resource(s) needed for testing.
	defer func(guid synthetics.EntityGUID) {
		client.Synthetics.SyntheticsDeletePrivateLocation(synthetics.EntityGUID(guid))
	}(privateLocationGUID)
}

func testAccNewRelicSyntheticsLocationDataSource(n string, attr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if _, ok := a[attr]; !ok {
			return fmt.Errorf("expected to read synthetics monitor location data from New Relic using attribute `%s`", attr)
		}
		return nil
	}
}

func testConfigDataSourceSyntheticsLocation(value string) string {
	return fmt.Sprintf(`
data "newrelic_synthetics_private_location" "bar" {
	name = "%s"
}
`, value)
}
