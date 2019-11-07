package version

import (
	"testing"
)

func TestDefaultProviderVersion(t *testing.T) {
	if ProviderVersion != "dev" {
		t.Fatal("incorrect default for ProviderVersion")
	}
}
