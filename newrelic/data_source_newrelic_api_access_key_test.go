//go:build integration || APIKS

package newrelic

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNewRelicAPIAccessKeyDataSource_ByKeyID(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAPIAccessKeyDataSourceConfigByKeyID(testSubAccountID, keyTypeIngestLicense, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "account_id", fmt.Sprintf("%d", testSubAccountID)),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestLicense),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "notes", keyNotes),
					resource.TestCheckResourceAttrSet("data.newrelic_api_access_key.foobar", "key_id"),
					resource.TestCheckResourceAttrSet("data.newrelic_api_access_key.foobar", "key"),
				),
			},
		},
	})
}

func TestAccNewRelicAPIAccessKeyDataSource_ByName(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNotes := fmt.Sprintf("tftest-keynotes-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNewRelicAPIAccessKeyDataSourceConfigByName(testSubAccountID, keyTypeIngestLicense, keyName, keyNotes),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "account_id", fmt.Sprintf("%d", testSubAccountID)),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "key_type", keyTypeIngest),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "ingest_type", keyTypeIngestLicense),
					resource.TestCheckResourceAttr("data.newrelic_api_access_key.foobar", "name", keyName),
					resource.TestCheckResourceAttrSet("data.newrelic_api_access_key.foobar", "key_id"),
					resource.TestCheckResourceAttrSet("data.newrelic_api_access_key.foobar", "key"),
				),
			},
		},
	})
}

func testAccNewRelicAPIAccessKeyDataSourceConfigByKeyID(accountID int, ingestType, name, notes string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id  = %d
	key_type    = "INGEST"
	ingest_type = "%s"
	name        = "%s"
	notes       = "%s"
}

data "newrelic_api_access_key" "foobar" {
	key_id   = newrelic_api_access_key.foobar.id
	key_type = "INGEST"
}
`, accountID, ingestType, name, notes)
}

func testAccNewRelicAPIAccessKeyDataSourceConfigByName(accountID int, ingestType, name, notes string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id  = %d
	key_type    = "INGEST"
	ingest_type = "%s"
	name        = "%s"
	notes       = "%s"
}

data "newrelic_api_access_key" "foobar" {
	account_id  = newrelic_api_access_key.foobar.account_id
	key_type    = "INGEST"
	ingest_type = "%s"
	name        = newrelic_api_access_key.foobar.name
}
`, accountID, ingestType, name, notes, ingestType)
}
