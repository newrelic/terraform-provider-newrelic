//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/pkg/synthetics"
)

func TestAccNewRelicSyntheticsMonitorLocationDataSource(t *testing.T) {
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
			time.Sleep(20 * time.Second)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigDataSourceSyntheticsLocation("label", privateLocationName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_monitor_location.bar", "label"),
				),
			},
			{
				Config: testConfigDataSourceSyntheticsLocation("name", privateLocationName),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicSyntheticsLocationDataSource("data.newrelic_synthetics_monitor_location.bar", "name"),
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

func testConfigDataSourceSyntheticsLocation(attr string, value string) string {
	return fmt.Sprintf(`
data "newrelic_synthetics_monitor_location" "bar" {
	%s = "%s"
}
`, attr, value)
}
