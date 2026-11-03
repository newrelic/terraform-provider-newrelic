//go:build integration || APIKS

package newrelic

import (
	"fmt"
	"regexp"
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

func TestAccNewRelicAPIAccessKey_IngestEnabled(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	keyNameUpdated := fmt.Sprintf("%s-updated", keyName)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Create a disabled key
				Config: testAccCheckNewRelicAPIAccessKeyIngestEnabled(testSubAccountID, keyName, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "false"),
					resource.TestCheckResourceAttrSet("newrelic_api_access_key.foobar", "key"),
				),
			},
			{
				// Ensure a disabled key is not treated as deleted and recreated
				Config:   testAccCheckNewRelicAPIAccessKeyIngestEnabled(testSubAccountID, keyName, "false"),
				PlanOnly: true,
			},
			{
				// Import a disabled key
				ResourceName:      "newrelic_api_access_key.foobar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNewRelicAPIAccessKeyImportStateIdFunc_Basic("newrelic_api_access_key.foobar"),
			},
			{
				// Update a disabled key
				Config: testAccCheckNewRelicAPIAccessKeyIngestEnabled(testSubAccountID, keyNameUpdated, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "name", keyNameUpdated),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "false"),
				),
			},
			{
				// Re-enable the key
				Config: testAccCheckNewRelicAPIAccessKeyIngestEnabled(testSubAccountID, keyNameUpdated, "true"),
				Check:  resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "true"),
			},
			{
				// Omitting `enabled` leaves the key's status unchanged
				Config: testAccCheckNewRelicAPIAccessKeyIngest(testSubAccountID, keyTypeIngestBrowser, keyNameUpdated, ""),
				Check:  resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "true"),
			},
			{
				// Disable the key again, so that a disabled key is deleted when the test is torn down
				Config: testAccCheckNewRelicAPIAccessKeyIngestEnabled(testSubAccountID, keyNameUpdated, "false"),
				Check:  resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "false"),
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_IngestOwningTeamIDs(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))
	teamName := fmt.Sprintf("tftest-team-%s", acctest.RandString(10))
	bothTeams := "[newrelic_team.first.id, newrelic_team.second.id]"
	firstTeam := "[newrelic_team.first.id]"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Create a key with owning teams
				Config: testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, bothTeams, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "owning_team_ids.#", "2"),
					resource.TestCheckTypeSetElemAttrPair("newrelic_api_access_key.foobar", "owning_team_ids.*", "newrelic_team.first", "id"),
					resource.TestCheckTypeSetElemAttrPair("newrelic_api_access_key.foobar", "owning_team_ids.*", "newrelic_team.second", "id"),
				),
			},
			{
				// Ensure a subsequent plan has no drift
				Config:   testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, bothTeams, ""),
				PlanOnly: true,
			},
			{
				// Replace the owning teams
				Config: testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, firstTeam, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "owning_team_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("newrelic_api_access_key.foobar", "owning_team_ids.*", "newrelic_team.first", "id"),
				),
			},
			{
				// Omitting `owning_team_ids` leaves the key's owning teams unchanged
				Config: testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "owning_team_ids.#", "1"),
					resource.TestCheckTypeSetElemAttrPair("newrelic_api_access_key.foobar", "owning_team_ids.*", "newrelic_team.first", "id"),
				),
			},
			{
				// Clear all owning teams
				Config: testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, "[]", ""),
				Check:  resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "owning_team_ids.#", "0"),
			},
			{
				// Ensure a subsequent plan has no drift
				Config:   testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, "[]", ""),
				PlanOnly: true,
			},
			{
				// Change owning teams and disable the key in the same apply
				Config: testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(testSubAccountID, keyName, teamName, bothTeams, "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "owning_team_ids.#", "2"),
					resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccNewRelicAPIAccessKey_UserEnabled(t *testing.T) {
	keyName := fmt.Sprintf("tftest-keyname-%s", acctest.RandString(10))

	if testUserID == 0 {
		t.Skipf("Skipping this test, as NEW_RELIC_TEST_USER_ID must be set for this test to run.")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Owning teams are not supported on user keys
				Config:      testAccCheckNewRelicAPIAccessKeyUserOwningTeamIDs(testAccountID, testUserID, keyName),
				PlanOnly:    true,
				ExpectError: regexp.MustCompile("the `owning_team_ids` attribute cannot be set when `key_type` is set to USER"),
			},
			{
				// Create a disabled key
				Config: testAccCheckNewRelicAPIAccessKeyUserEnabled(testAccountID, testUserID, keyName, "false"),
				Check:  resource.TestCheckResourceAttr("newrelic_api_access_key.foobar", "enabled", "false"),
			},
			{
				// Ensure a disabled key is not treated as deleted and recreated
				Config:   testAccCheckNewRelicAPIAccessKeyUserEnabled(testAccountID, testUserID, keyName, "false"),
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

func testAccCheckNewRelicAPIAccessKeyIngestEnabled(accountID int, name, enabled string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id  = %d
	key_type    = "INGEST"
	ingest_type = "BROWSER"
	name        = "%s"
	enabled     = %s
}
`, accountID, name, enabled)
}

// testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs returns a configuration with two temporary teams and an ingest key.
// owningTeamIDs and enabled are HCL expressions; an empty value omits the attribute from the key.
func testAccCheckNewRelicAPIAccessKeyIngestOwningTeamIDs(accountID int, name, teamName, owningTeamIDs, enabled string) string {
	optionalAttributes := ""
	if owningTeamIDs != "" {
		optionalAttributes += fmt.Sprintf("\n\towning_team_ids = %s", owningTeamIDs)
	}
	if enabled != "" {
		optionalAttributes += fmt.Sprintf("\n\tenabled         = %s", enabled)
	}

	return fmt.Sprintf(`
resource "newrelic_team" "first" {
	name = "%[3]s-1"
}

resource "newrelic_team" "second" {
	name = "%[3]s-2"
}

resource "newrelic_api_access_key" "foobar" {
	account_id      = %[1]d
	key_type        = "INGEST"
	ingest_type     = "LICENSE"
	name            = "%[2]s"%[4]s
}
`, accountID, name, teamName, optionalAttributes)
}

func testAccCheckNewRelicAPIAccessKeyUserEnabled(accountID, userID int, name, enabled string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id = %d
	key_type   = "USER"
	user_id    = "%d"
	name       = "%s"
	enabled    = %s
}
`, accountID, userID, name, enabled)
}

func testAccCheckNewRelicAPIAccessKeyUserOwningTeamIDs(accountID, userID int, name string) string {
	return fmt.Sprintf(`
resource "newrelic_api_access_key" "foobar" {
	account_id      = %d
	key_type        = "USER"
	user_id         = "%d"
	name            = "%s"
	owning_team_ids = ["team-guid"]
}
`, accountID, userID, name)
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
