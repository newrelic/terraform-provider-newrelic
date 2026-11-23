//go:build unit

package newrelic

import (
	"bytes"
	"encoding/json"
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

func TestExpandAwsConnectionCredentialOmitsExternalID(t *testing.T) {
	for name, externalID := range map[string]interface{}{
		"value in state": "from-state",
		"null":           nil,
	} {
		in := []interface{}{map[string]interface{}{
			"assume_role": []interface{}{map[string]interface{}{
				"role_arn":    "arn:aws:iam::123456789012:role/test",
				"external_id": externalID,
			}},
		}}

		for kind, input := range map[string]interface{}{
			"create": expandAwsConnectionCredential(in),
			"update": expandAwsConnectionCredentialUpdate(in),
		} {
			b, err := json.Marshal(input)
			if err != nil {
				t.Fatalf("%s/%s: %v", name, kind, err)
			}
			if bytes.Contains(b, []byte("externalId")) {
				t.Fatalf("%s/%s: externalId must be absent, got %s", name, kind, b)
			}
		}
	}
}
