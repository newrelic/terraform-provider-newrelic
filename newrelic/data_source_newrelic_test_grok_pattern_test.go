//go:build integration || LOGGING_INTEGRATIONS
// +build integration LOGGING_INTEGRATIONS

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicTestGrokDataSource_Basic(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//read
			{
				Config: testAccNewRelicTestGrokDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicTestGrokDataSourceExists("data.newrelic_test_grok_pattern.grok")),
			},
		},
	})
}

func TestAccNewRelicTestGrokDataSource_InvalidGrokPattern(t *testing.T) {

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//read
			{
				Config:      testAccNewRelicTestGrokInvalidPatternDataSourceConfig(),
				ExpectError: regexp.MustCompile("Invalid Grok pattern"),
			},
		},
	})
}

func testAccNewRelicTestGrokInvalidPatternDataSourceConfig() string {
	return fmt.Sprintf(`
data "newrelic_test_grok_pattern" "grok"{
	account_id = %[1]d
grok = "{IP:host_ip}"
	log_lines = ["host_ip: 43.3.120.2","bytes_received: 2048"]
}
`, testAccountID)
}

func testAccNewRelicTestGrokDataSourceConfig() string {
	return fmt.Sprintf(`
data "newrelic_test_grok_pattern" "grok"{
	account_id = %[1]d
	grok = "%%%%{IP:host_ip}"
	log_lines = ["host_ip: 43.3.120.2","bytes_received: 2048"]
}
`, testAccountID)
}

func testAccCheckNewRelicTestGrokDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["grok"] == "" {
			return fmt.Errorf("expected to get a grok from New Relic")
		}
		return nil
	}
}
