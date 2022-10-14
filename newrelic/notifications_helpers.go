package newrelic

import (
	"fmt"

	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func notificationsPropertySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification property key.",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == "source"
				},
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Notification property value.",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == "terraform"
				},
			},
			"label": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notification property label.",
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					return oldValue == "terraform-source-internal"
				},
			},
			"display_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Notification property display key.",
			},
		},
	}
}

// Builds an array of typed notifications error interface based on the GraphQL `response.errors` array.
func buildAiNotificationsErrors(errors []ai.AiNotificationsError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		if len(err.Fields) > 0 {
			diagErrors = append(diagErrors, buildAiNotificationsDataValidationError(ai.AiNotificationsDataValidationError{
				Details: err.Details,
				Fields:  err.Fields,
			}))
		} else if len(err.Dependencies) > 0 {
			diagErrors = append(diagErrors, buildAiNotificationsConstraintError(ai.AiNotificationsConstraintError{
				Name:         err.Name,
				Dependencies: err.Dependencies,
			}))
		} else {
			diagErrors = append(diagErrors, buildAiNotificationsResponseError(ai.AiNotificationsResponseError{
				Description: err.Description,
				Details:     err.Details,
				Type:        err.Type,
			}))
		}
	}

	return diagErrors
}

// Builds data validation error based on the GraphQL `response.error`.
func buildAiNotificationsDataValidationError(err ai.AiNotificationsDataValidationError) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("%s", err.Fields),
	}
}

// Builds constrain error based on the GraphQL `response.error`.
func buildAiNotificationsConstraintError(err ai.AiNotificationsConstraintError) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("%s: %s", err.Name, err.Dependencies),
	}
}

// Builds response error based on the GraphQL `response.error`.
func buildAiNotificationsResponseError(err ai.AiNotificationsResponseError) diag.Diagnostic {
	return diag.Diagnostic{
		Severity: diag.Error,
		Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
	}
}

// Builds an array of typed notifications response errors based on the GraphQL `response.errors` array.
func buildAiNotificationsResponseErrors(errors []notifications.AiNotificationsResponseError) diag.Diagnostics {
	var diagErrors diag.Diagnostics
	for _, err := range errors {
		diagErrors = append(diagErrors, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("%s: %s", string(err.Type), err.Description),
		})
	}
	return diagErrors
}

func createMonitoringProperty() notifications.AiNotificationsPropertyInput {
	return notifications.AiNotificationsPropertyInput{
		Key:   "source",
		Value: "terraform",
		Label: "terraform-source-internal",
	}
}
