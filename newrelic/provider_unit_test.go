//go:build unit
// +build unit

package newrelic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestGenerateNameForIntegrationTestResource(t *testing.T) {
	t.Parallel()

	result := generateNameForIntegrationTestResource()
	require.Contains(t, result, "tf-test-")
}

func TestBuildUserAgentStringWithDefaultServiceName(t *testing.T) {
	t.Parallel()

	tfUA := "HashiCorp Terraform/1.3.5 (+https://www.terraform.io) Terraform Plugin SDK/2.10.1"

	result := buildUserAgentString(tfUA, getUserAgentServiceName(), ProviderVersion)
	require.Contains(t, result, " terraform-provider-newrelic/dev---fixthisbcitsatest")
}

func TestBuildUserAgentStringWithCustomServiceName(t *testing.T) {
	t.Parallel()

	// In a real scenario, this is set via -ldflags at compile time.
	// This would be "pulumi" or some other third-party service in a real scenario.
	UserAgentServiceName = "test"

	tfUA := "HashiCorp Terraform/1.3.5 (+https://www.terraform.io) Terraform Plugin SDK/2.10.1"
	result := buildUserAgentString(tfUA, getUserAgentServiceName(), ProviderVersion)
	require.Contains(t, result, " test/terraform-provider-newrelic/dev")

	// Reset the package variable to default to avoid polluting other tests.
	UserAgentServiceName = ""
}

func TestGetUserAgentServiceNameDefault(t *testing.T) {
	t.Parallel()

	result := getUserAgentServiceName()
	require.Equal(t, "terraform-provider-newrelic", result)
}

func TestGetUserAgentServiceNameCustom(t *testing.T) {
	t.Parallel()

	// In a real scenario, this is set via -ldflags at compile time.
	UserAgentServiceName = "test"

	result := getUserAgentServiceName()
	require.Equal(t, "test/terraform-provider-newrelic", result)

	// Reset the package variable to default to avoid polluting other tests.
	UserAgentServiceName = ""
}
