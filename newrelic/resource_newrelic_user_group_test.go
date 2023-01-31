//go:build integration
// +build integration

package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccNewRelicUserGroup(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEnvVars(t)
		},
		ProviderFactories: testAccProviders,
	})
}
