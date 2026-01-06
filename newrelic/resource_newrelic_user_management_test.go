//go:build integration || AUTH

package newrelic

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
)

var userNamePrefix string = "terraform-provider-newrelic-integration-test-mock-user"
var userEmailPrefix string = "developer-toolkit+#@newrelic.com"
var authenticationDomainId string = "84cb286a-8eb0-4478-b469-cdf2ccfef553"
var mockAuthenticationDomainId string = "fae55e6b-b1ce-4a0f-83b2-ee774798f2cc"
var userManagementResourceName string = "newrelic_user.foo"

func TestAccNewRelicUser_Basic(t *testing.T) {
	createMap := map[string]string{
		"name":                     fmt.Sprintf("%s-%s", userNamePrefix, testhelpers.RandSeq(10)),
		"email_id":                 strings.ReplaceAll(userEmailPrefix, "#", testhelpers.RandSeq(10)),
		"authentication_domain_id": authenticationDomainId,
		"user_type":                "CORE_USER_TIER",
	}
	updateMap := map[string]string{
		"name":                     fmt.Sprintf("%s-%s-%s", userNamePrefix, testhelpers.RandSeq(10), "updated"),
		"email_id":                 strings.ReplaceAll(userEmailPrefix, "#", testhelpers.RandSeq(5)),
		"authentication_domain_id": authenticationDomainId,
		"user_type":                "BASIC_USER_TIER",
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccNewRelicUserManagementConfiguration(createMap),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckUserExists(userManagementResourceName),
				),
			},
			// Update
			{
				Config: testAccNewRelicUserManagementConfiguration(updateMap),
				Check: resource.ComposeTestCheckFunc(
					testAccNewRelicCheckUserExists(userManagementResourceName),
				),
			},
			// Import
			{
				ResourceName:      userManagementResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNewRelicUser_InvalidAuthenticationDomainError(t *testing.T) {
	createMap := map[string]string{
		"name":                     fmt.Sprintf("%s-%s", userNamePrefix, testhelpers.RandSeq(10)),
		"email_id":                 strings.ReplaceAll(userEmailPrefix, "#", testhelpers.RandSeq(10)),
		"authentication_domain_id": mockAuthenticationDomainId,
		"user_type":                "CORE_USER_TIER",
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicUserManagementConfiguration(createMap),
				ExpectError: regexp.MustCompile(`Could not find the target or you are unauthorized.`),
			},
		},
	})
}

func TestAccNewRelicUser_EmailAlreadyTakenError(t *testing.T) {
	createMap := map[string]string{
		"name": fmt.Sprintf("%s-%s", userNamePrefix, testhelpers.RandSeq(10)),
		// a user with the following email_id already exists, hence, this test is expected to fail
		"email_id":                 strings.ReplaceAll(userEmailPrefix, "#", "integration"),
		"authentication_domain_id": authenticationDomainId,
		"user_type":                "CORE_USER_TIER",
	}

	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccNewRelicUserManagementConfiguration(createMap),
				ExpectError: regexp.MustCompile(`has already been taken within authentication domain`),
			},
		},
	})
}

func testAccNewRelicCheckUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no user ID found")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		resp, err := client.UserManagement.UserManagementGetUsers([]string{authenticationDomainId}, []string{rs.Primary.ID}, "", "")
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		for _, authDomain := range resp.AuthenticationDomains {
			if authDomain.ID == authenticationDomainId {
				for _, u := range authDomain.Users.Users {
					if u.ID != rs.Primary.ID {
						return fmt.Errorf("user not found")
					}
				}
			}

		}

		return nil
	}
}

func testAccNewRelicUserManagementConfiguration(values map[string]string) string {
	return fmt.Sprintf(`
	resource "newrelic_user" "foo" {
	  name  = "%s"
	  email_id = "%s"
	  authentication_domain_id = "%s"
	  user_type = "%s"
	}
`,
		values["name"],
		values["email_id"],
		values["authentication_domain_id"],
		values["user_type"],
	)
}
