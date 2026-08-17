terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US, EU, or JP
}

# Verifies WORKFLOW_AUTOMATION notification destination creation.
# The property block is required by the provider schema for all destination
# types. For this type it carries no effect — leave key and value empty.
# Run with: terraform apply (uses NEW_RELIC_ACCOUNT_ID + NEW_RELIC_API_KEY env vars)
resource "newrelic_notification_destination" "workflow_automation_test" {
  name = "tf-test-workflow-automation-destination"
  type = "WORKFLOW_AUTOMATION"

  property {
    key   = ""
    value = ""
  }

  auth_custom_header {
    key   = "Api-Key"
    value = var.nr_api_key
  }
}

variable "nr_api_key" {
  description = "New Relic User API Key for WORKFLOW_AUTOMATION auth header"
  sensitive   = true
}
