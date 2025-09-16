//go:build integration || APIKS
// +build integration APIKS

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNewRelicAPIAccessKey_BasicIngestBrowser(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestBrowser, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", fmt.Sprintf("%d", testSubAccountID)),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestBrowser),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
			{
				// Ensure a subsequent plan has no drift
				Config:   testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestBrowser, keyName, keyNotes),
				PlanOnly: true,
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_BasicIngestLicense(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestLicense, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", fmt.Sprintf("%d", testSubAccountID)),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestLicense),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
			{
				// Ensure a subsequent plan has no drift
				Config:   testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestLicense, keyName, keyNotes),
				PlanOnly: true,
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_BasicUser(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))

	if testUserID == 0 {
		t.Skipf("Skipping this test, as NEW_RELIC_TEST_USER_ID must be set for this test to run.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyUser(testSubAccountID, testUserID, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", fmt.Sprintf("%d", testSubAccountID)),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeUser),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "user_id", fmt.Sprintf("%d", testUserID)),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
			{
				// Ensure a subsequent plan has no drift
				Config:   testAccCheckNewRelicAPIAccessKeyUser(testSubAccountID, testUserID, keyName, keyNotes),
				PlanOnly: true,
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_BasicIngestBrowserNoNotesNames(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngestNoNameNotes(testSubAccountID, keyTypeIngestBrowser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", fmt.Sprintf("%d", testSubAccountID)),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestBrowser),
					// ignoring the following since there is a drift - when no name is assigned, the API gives an autogenerated name
					// resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", ""),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", ""),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
			{
				// Ensure a subsequent plan has no drift
				Config:   testAccCheckNewRelicAPIAccessKeyIngestNoNameNotes(testSubAccountID, keyTypeIngestBrowser),
				PlanOnly: true,
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_ImportBasic(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestLicense, keyName, keyNotes),
			},
			{
				ResourceName:      "newrelic_api_access_key.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNewRelicAPIAccessKeyImportStateIdFunc_Basic("newrelic_api_access_key.foobar"),
			},
			{
				// Planning after import should not show changes
				Config:   testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestLicense, keyName, keyNotes),
				PlanOnly: true,
			},
		},
	})
}

func testAccNewRelicAPIAccessKeyImportStateIdFunc_Basic(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("[ERROR] Not found: %s", resourceName)
		}

		return fmt.Sprintf("%s:%s", rs.Primary.Attributes["id"], rs.Primary.Attributes["key_type"]), nil
	}

}

func testAccCheckNewRelicAPIAccessKeyIngestNoNameNotes(accountID int, ingestType string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id  = %d
	key_type    = "INGEST"
	ingest_type = "%s"
}
`, accountID, ingestType)
}

func testAccCheckNewRelicAPIAccessKeyIngest(accountID int, ingestType, name, notes string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id  = %d
	key_type    = "INGEST"
	ingest_type = "%s"
	name        = "%s"
    notes       = "%s"
}
`, accountID, ingestType, name, notes)
}

func testAccCheckNewRelicAPIAccessKeyUser(accountID, userID int, name, notes string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id  = %d
	key_type    = "USER"
	user_id     = "%d"
	name        = "%s"
    notes       = "%s"
}
`, accountID, userID, name, notes)
}
