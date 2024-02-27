//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var groupNamePrefix string = "terraform-provider-newrelic-integration-test-mock-group"
var groupManagementResourceName string = "newrelic_group.foo"

func TestAccNewRelicGroup_Basic(t *testing.T) {
	createMap := map[string]string{
		"name":                     fmt.Sprintf("%s-%s", groupNamePrefix, testhelpers.RandSeq(10)),
		"authentication_domain_id": authenticationDomainId,
		"user_ids":                 "",
	}
	updateMap := map[string]string{
		"name":                     fmt.Sprintf("%s-%s-%s", groupNamePrefix, testhelpers.RandSeq(10), "updated"),
		"authentication_domain_id": authenticationDomainId,
		"user_ids":                 "\"1005772336\", \"1005772340\"",
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicGroupManagementConfiguration(createMap),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckGroupExists(groupManagementResourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicGroupManagementConfiguration(updateMap),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckGroupExists(groupManagementResourceName),
				),
			},
			// Import
			{
				ResourceName:      groupManagementResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicGroup_InvalidAuthenticationDomainError(t *testing.T) {
	createMap := map[string]string{
		"name":                     fmt.Sprintf("%s-%s", groupNamePrefix, testhelpers.RandSeq(10)),
		"authentication_domain_id": mockAuthenticationDomainId,
		"user_ids":                 "",
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicGroupManagementConfiguration(createMap),
				ExpectError: regexp.MustCompile(`Could not find the target or you are unauthorized.`),
			},
		},
	})
}

func TestAccNewRelicGroup_DuplicateNameError(t *testing.T) {
	createMap := map[string]string{
		"name":                     "Integration Test Group 1 DO NOT DELETE",
		"authentication_domain_id": authenticationDomainId,
		"user_ids":                 "",
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicGroupManagementConfiguration(createMap),
				ExpectError: regexp.MustCompile(`Display name has already been taken`),
			},
		},
	})
}

func testAccNewRelicCheckGroupExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no group ID found")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resp, err := client.UserManagement.UserManagementGetGroupsWithUsers([]string{authenticationDomainId}, []string{rs.Primary.ID}, "")
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		for _, authDomain := range resp.AuthenticationDomains {
			if authDomain.ID == authenticationDomainId {
				for _, u := range authDomain.Groups.Groups {
					if u.ID != rs.Primary.ID {
						return fmt.Errorf("group not found")
					}
				}
			}

		}

		return nil
	}
}

func testAccNewRelicGroupManagementConfiguration(values map[string]string) string {
	return fmt.Sprintf(`
	resource "newrelic_group" "foo" {
  		name                     = "%s"
  		authentication_domain_id = "%s"
		user_ids 			 	 = [%s] 
	}
`,
		values["name"],
		values["authentication_domain_id"],
		values["user_ids"],
	)
}
