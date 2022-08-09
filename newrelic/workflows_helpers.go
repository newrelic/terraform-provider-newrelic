package newrelic

import (
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/workflows"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

// Builds an array of typed workflow create response errors based on the GraphQL `response.errors` array.
func buildAiWorkflowsCreateResponseError(errors []workflows.AiWorkflowsCreateResponseError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		diagErrors = append(diagErrors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
		})
	}
	return diagErrors
}

// Builds an array of typed workflow update response errors based on the GraphQL `response.errors` array.
func buildAiWorkflowsUpdateResponseError(errors []workflows.AiWorkflowsUpdateResponseError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		diagErrors = append(diagErrors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
		})
	}
	return diagErrors
}

// Builds an array of typed workflow delete response errors based on the GraphQL `response.errors` array.
func buildAiWorkflowsDeleteResponseError(errors []workflows.AiWorkflowsDeleteResponseError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		diagErrors = append(diagErrors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
		})
	}
	return diagErrors
}
