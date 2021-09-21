//go:build integration
// +build integration

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAPIAccessKey_BasicIngestBrowser(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))
	accountIDRaw, accountID := retrieveIdsFromEnvOrSkip(t, "NEW_RELIC_TEST_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngest(accountID, keyTypeIngestBrowser, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", accountIDRaw),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestBrowser),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_BasicIngestLicense(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))
	accountIDRaw, accountID := retrieveIdsFromEnvOrSkip(t, "NEW_RELIC_TEST_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngest(accountID, keyTypeIngestLicense, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", accountIDRaw),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestLicense),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_BasicUser(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))
	accountIDRaw, accountID := retrieveIdsFromEnvOrSkip(t, "NEW_RELIC_TEST_ACCOUNT_ID")
	userIDRaw, userID := retrieveIdsFromEnvOrSkip(t, "NEW_RELIC_TEST_USER_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyUser(accountID, userID, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", accountIDRaw),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeUser),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "user_id", userIDRaw),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_BasicIngestBrowserNoNotesNames(t *testing.T) {
	accountIDRaw, accountID := retrieveIdsFromEnvOrSkip(t, "NEW_RELIC_TEST_ACCOUNT_ID")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckNewRelicAPIAccessKeyIngestNoNameNotes(accountID, keyTypeIngestBrowser),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "account_id", accountIDRaw),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestBrowser),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", ""),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "notes", ""),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
		},
	})
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
