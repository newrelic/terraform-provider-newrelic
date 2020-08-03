// +build integration

package newrelic

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestAccNewRelicProvider_Region(t *testing.T) {
	// This error message will occur when configuring
	// US region with EU API URLs when using the TF test account.
	expectedErrorMsg := "403 response returned"
	rName := acctest.RandString(5)

	resource.ParallelTest(t, resource.TestCase{
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
				ExpectError: regexp.MustCompile(expectedErrorMsg),
			},
			// Test: Override US region URLs with EU region URLs (will result in an auth error)
			{
				Config:      testAccNewRelicProviderConfig("US", `nerdgraph_api_url = "https://api.eu.newrelic.com/graphql"`, rName),
				ExpectError: regexp.MustCompile(expectedErrorMsg),
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
