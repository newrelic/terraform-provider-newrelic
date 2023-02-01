//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicUserGroup(t *testing.T) {
	rName := generateNameForIntegrationTestResource()
	authId := "0cc21d98-8dc2-484a-bb26-258e17ede584"
	resourceName := "newrelic_user_group.bar"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckEnvVars(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNewRelicGroupDestroy,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicGroupConfig(rName, authId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicGroupExists(resourceName)),
			},
		},
	})
}

func testAccCheckNewRelicGroupDestroy(s *terraform.State) error {
	providerConfig := testAccProvider.Meta().(*ProviderConfig)
	client := providerConfig.NewClient

	for _, r := range s.RootModule().Resources {
		if r.Type != "newrelic_user_group" {
			continue
		}
		authenticationDomainId := "0cc21d98-8dc2-484a-bb26-258e17ede584"
		resp, err := getUserGroupID(context.Background(), client, authenticationDomainId, r.Primary.ID)
		if resp != nil {
			fmt.Errorf("groups still exists")
		}

		if err != nil {
			return err
		}

	}
	return nil
}

func testAccCheckNewRelicGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient
		authenticationDomainId := "0cc21d98-8dc2-484a-bb26-258e17ede584"

		_, err := getUserGroupID(context.Background(), client, authenticationDomainId, rs.Primary.ID)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccNewRelicGroupConfig(name string, authId string) string {
	return fmt.Sprintf(`
resource "newrelic_user_group" "bar"{
name = "%[1]s" 
authentication_domain_id = "%[2]s"
}
`, name, authId)
}
