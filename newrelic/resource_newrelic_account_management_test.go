//go:build integration
// +build integration

package newrelic

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicAccountManagement(t *testing.T) {
	resourceName := "newrelic_account_management.foo"
	rName := acctest.RandString(7)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: testAccNewRelicAccountCreateConfig("Test " + rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountExists(resourceName)),
			},
			//update
			{
				Config: testAccNewRelicAccountUpdateConfig("Update " + rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNewRelicAccountExists(resourceName)),
			},
			//import
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      resourceName,
			},
		},
	})
}
func TestAccNewRelicAccountManagementInvalidRegion(t *testing.T) {
	rName := acctest.RandString(7)
	expectedErrorMsg := regexp.MustCompile(`expected region to be one of \[us01 eu01\], got abcd01`)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicAccountCreateInvalidRegionConfig("Test " + rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}
func TestAccNewRelicAccountManagementInCorrectRegion(t *testing.T) {
	rName := acctest.RandString(7)
	expectedErrorMsg := regexp.MustCompile(`An error occurred resolving this field`)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config:      testAccNewRelicAccountCreateConfigInCorrectRegion("Test " + rName),
				ExpectError: expectedErrorMsg,
			},
		},
	})
}

func testAccNewRelicAccountCreateConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[2]s"
	region= "us01"
}
`, testAccountID, name)
}

func testAccNewRelicAccountUpdateConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[2]s"
	region= "us01"
}
`, testAccountID, name)
}

func testAccCheckNewRelicAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		client := testAccProvider.Meta().(*ProviderConfig).NewClient

		account, err := getCreatedAccountByID(context.Background(), client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if account == nil {
			return fmt.Errorf("account not found")
		}

		return nil
	}
}

func testAccNewRelicAccountCreateInvalidRegionConfig(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[2]s"
	region= "abcd01"
}
`, testAccountID, name)
}

func testAccNewRelicAccountCreateConfigInCorrectRegion(name string) string {
	return fmt.Sprintf(`
resource "newrelic_account_management" "foo"{
	name=  "%[2]s"
	region= "eu01"
}
`, testAccountID, name)
}
