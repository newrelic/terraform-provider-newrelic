// +build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccExpectedAlertChannelName string
	testAccExpectedApplicationName  string
	testAccExpectedAlertPolicyName  string
	testAccAPIKey                   string
	testAccProviders                map[string]terraform.ResourceProvider
	testAccProvider                 *schema.Provider
	testAccountID                   int
	testAccountName                 string
	//testAccCleanupComplete          = false
)

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfig(t *testing.T) {
	c := ProviderConfig{
		PersonalAPIKey: "abc123",
		AccountID:      123,
	}

	hasNerdGraphCreds := c.hasNerdGraphCredentials()

	if !hasNerdGraphCreds {
		t.Error("hasNerdGraphCreds should be true")
	}
}
