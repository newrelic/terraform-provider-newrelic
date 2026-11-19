//go:build unit

package newrelic

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAwsConnectionAssumeRoleExternalIDIsComputed(t *testing.T) {
	credential := resourceNewRelicAwsConnection().Schema["credential"].Elem.(*schema.Resource)
	assumeRole := credential.Schema["assume_role"].Elem.(*schema.Resource)
	externalID := assumeRole.Schema["external_id"]

	if !externalID.Computed || externalID.Optional || externalID.Required {
		t.Fatalf("external_id should be computed only, got computed=%v optional=%v required=%v",
			externalID.Computed, externalID.Optional, externalID.Required)
	}
}

func TestExpandAwsConnectionCredentialDoesNotSendExternalID(t *testing.T) {
	in := []interface{}{map[string]interface{}{
		"assume_role": []interface{}{map[string]interface{}{
			"role_arn":    "arn:aws:iam::123456789012:role/test",
			"external_id": "from-state",
		}},
	}}

	if got := expandAwsConnectionCredential(in).AssumeRole.ExternalId; got != "" {
		t.Fatalf("create input should not carry external_id, got %q", got)
	}
	if got := expandAwsConnectionCredentialUpdate(in).AssumeRole.ExternalId; got != "" {
		t.Fatalf("update input should not carry external_id, got %q", got)
	}
}
