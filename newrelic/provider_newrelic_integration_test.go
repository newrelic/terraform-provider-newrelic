//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func TestAccNewRelicProvider_Region(t *testing.T) {
	// This error message will occur when configuring
	// US region with EU API URLs when using the TF test account.
	expectedErrorMsg := "403 response returned"
	expectedErrorMsgTwo := "Access denied."
	expectedErrorMsgRegex := fmt.Sprintf("%s|%s", expectedErrorMsg, expectedErrorMsgTwo)
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			// Test: Region "US"
			{
				Config: testAccNewRelicProviderConfig("US", "", rName),
			},
			// Test: Region "EU"
			{
				Config:      testAccNewRelicProviderConfig("EU", "", rName),
				ExpectError: regexp.MustCompile(expectedErrorMsgRegex),
			},
			// Test: Override US region URLs with EU region URLs (will result in an auth error)
			{
				Config:      testAccNewRelicProviderConfig("US", `nerdgraph_api_url = "https://api.eu.newrelic.com/graphql"`, rName),
				ExpectError: regexp.MustCompile(expectedErrorMsgRegex),
			},
			// Test: Override EU region URLs with US region URLs (should work since the TF acct is US-based)
			{
				Config: testAccNewRelicProviderConfig("EU", `nerdgraph_api_url = "https://api.newrelic.com/graphql"`, rName),
			},
			// Test: Case insensitivity
			{
				Config: testAccNewRelicProviderConfig("us", "", rName),
			},
		},
	})
}

func TestAccNewRelicProvider_UserAgent(t *testing.T) {
	rName := generateNameForIntegrationTestResource()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEnvVars(t) },
		Providers:    testAccProviders,
		CheckDestroy: func(*terraform.State) error { return nil },
		Steps: []resource.TestStep{
			// Test default user agent
			{
				Config: testAccNewRelicProviderConfig("US", "", rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicProviderConfigurationHasDefaultUserAgent(),
				),
			},
		},
	})
}

func testAccCheckNewRelicProviderConfigurationHasDefaultUserAgent() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerConfig := testAccProvider.Meta().(*ProviderConfig)
		ua := providerConfig.GetUserAgent()
		uaSplit := strings.Split(ua, " ")
		uaServiceName := uaSplit[len(uaSplit)-1]

		if strings.Contains(ua, "/terraform-provider/dev") {
			return fmt.Errorf(`user agent service name "%s" does not match expected result`, uaServiceName)
		}

		if !strings.EqualFold(uaServiceName, "terraform-provider-newrelic/dev") {
			return fmt.Errorf(`user agent service name "%s" does not match expected result`, uaServiceName)
		}

		return nil
	}
}
