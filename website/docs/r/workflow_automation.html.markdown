---
layout: "newrelic"
page_title: "New Relic: newrelic_workflow_automation"
sidebar_current: "docs-newrelic-resource-workflow-automation"
description: |-
Create and manage workflow automation in New Relic.
---

# Resource: newrelic\_workflow\_automation

Use this resource to create and manage New Relic Workflow Automation.

Workflow Automation allows you to define automated workflows using YAML definitions. These workflows can be scoped to either an account or an organization and support various automation steps and configurations.

## Example Usage

### Basic Workflow Automation with ACCOUNT Scope

```hcl
resource "newrelic_workflow_automation" "test_query" {
  name       = "test_query_workflow"
  scope_id   = "your-account-id"
  scope_type = "ACCOUNT"

  definition = <<-YAML
name: test_query_workflow
description: Simple workflow that queries NRDB and waits
steps:
  - name: queryNrdb
    type: action
    action: newrelic.nrdb.query
    version: 1
    inputs:
      query: SELECT count(*) from Log LIMIT 10
  - name: wait
    type: wait
    seconds: 3
    signals: []
    next: end
YAML
}
```

### Workflow Automation with ORGANIZATION Scope

```hcl
resource "newrelic_workflow_automation" "org_query" {
  name       = "org_query_workflow"
  scope_id   = "your-organization-id"
  scope_type = "ORGANIZATION"

  definition = <<-YAML
name: org_query_workflow
description: Organization-level workflow that queries NRDB
steps:
  - name: queryNrdb
    type: action
    action: newrelic.nrdb.query
    version: 1
    inputs:
      query: SELECT count(*) from Transaction LIMIT 10
  - name: wait
    type: wait
    seconds: 5
    signals: []
    next: end
YAML
}
```

### Advanced Workflow with Inputs and Slack Integration

This example demonstrates a more complex workflow that queries NRDB, transforms results to CSV, and posts them to Slack:

```hcl
resource "newrelic_workflow_automation" "nrql_slack_report" {
  name       = "nrqlSlackReport"
  scope_id   = "1234567"
  scope_type = "ACCOUNT"

  definition = <<-YAML
name: nrqlSlackReport

workflowInputs:
  nrql:
    type: String
  accountIds:
    type: List
  channel:
    type: String
  slackToken:
    type: String
    defaultValue: "${{ :secrets:slackToken }}"

steps:
  - name: queryForLog
    type: action
    action: newrelic.nrdb.query
    version: 1
    inputs:
      accountIds: "${{ .workflowInputs.accountIds }}"
      query: "${{ .workflowInputs.nrql }}"

  - name: generateCSV
    type: action
    action: utils.transform.toCSV
    version: 1
    inputs:
      data: '${{ .steps.queryForLog.outputs.results }}'

  - name: postCSV
    type: action
    action: slack.chat.postMessage
    version: 1
    inputs:
      channel: ${{ .workflowInputs.channel }}
      text: "NRQL Results"
      attachment:
        filename: 'results.csv'
        content: "${{ .steps.generateCSV.outputs.csv }}"
      token: ${{ .workflowInputs.slackToken }}
YAML
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workflow automation. This must match the `name` field in the YAML definition provided in the `definition` argument. **Important**: Changes to this field will force a new resource to be created.
* `definition` - (Required) The YAML definition of the workflow automation. This should be a valid YAML string that includes a `name` field matching the resource `name` argument, and defines the workflow steps and configuration.
* `scope_id` - (Required) The scope ID for the workflow automation. For `ACCOUNT` scope, this should be your New Relic account ID (numeric). For `ORGANIZATION` scope, this should be your organization ID (string). **Important**: Changes to this field will force a new resource to be created.
* `scope_type` - (Required) The scope type for the workflow automation. Must be either `ACCOUNT` or `ORGANIZATION`. **Important**: Changes to this field will force a new resource to be created.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the workflow automation. The ID format is `<scope_type>#<scope_id>#<workflow_name>`.
* `description` - The description of the workflow automation, as defined in the YAML definition.
* `version` - The current version number of the workflow automation.

## YAML Definition Structure

The `definition` argument accepts a YAML string that defines the workflow automation. The YAML must include the following fields:

### Required Fields

* `name` - (Required) The name of the workflow. This **must** match the `name` argument in the Terraform resource.
* `description` - (Optional but recommended) A description of what the workflow automation does.
* `steps` - (Required) An array of steps that define the workflow automation logic.

### Workflow Inputs

Workflows can define inputs that can be passed when the workflow is executed:

```yaml
workflowInputs:
  inputName:
    type: String       # Can be String, Number, Boolean, List, etc.
    defaultValue: "value"  # Optional default value
```

Inputs can reference secrets using the syntax: `${{ :secrets:secretName }}`

### Step Types

Each step in the `steps` array can be of different types:

#### Action Steps

Action steps execute specific actions like querying NRDB, transforming data, or sending notifications:

* `type: action` - Executes an action
  * `action` - The action to execute (e.g., `newrelic.nrdb.query`, `utils.transform.toCSV`, `slack.chat.postMessage`)
  * `version` - The version of the action to use
  * `inputs` - Input parameters for the action

Common actions include:
- `newrelic.nrdb.query` - Query New Relic database
- `utils.transform.toCSV` - Transform data to CSV format
- `slack.chat.postMessage` - Send messages to Slack
- And many more...

#### Wait Steps

* `type: wait` - Pauses the workflow execution
  * `seconds` - Number of seconds to wait
  * `signals` - Optional array of signals to wait for
  * `next` - Optional next step name (use "end" to terminate)

### Referencing Data in Workflows

You can reference data from previous steps and inputs using template expressions:

* Workflow inputs: `${{ .workflowInputs.inputName }}`
* Step outputs: `${{ .steps.stepName.outputs.fieldName }}`
* Secrets: `${{ :secrets:secretName }}`

### Example YAML Structure

**Simple workflow with query and wait:**

```yaml
name: example-workflow
description: Query logs and wait
steps:
  - name: queryNrdb
    type: action
    action: newrelic.nrdb.query
    version: 1
    inputs:
      query: SELECT count(*) from Log LIMIT 10
  - name: wait
    type: wait
    seconds: 5
    next: end
```

**Advanced workflow with inputs and multiple actions:**

```yaml
name: advanced-workflow
workflowInputs:
  accountId:
    type: String
  query:
    type: String
steps:
  - name: runQuery
    type: action
    action: newrelic.nrdb.query
    version: 1
    inputs:
      accountIds: "${{ .workflowInputs.accountId }}"
      query: "${{ .workflowInputs.query }}"
  - name: transformData
    type: action
    action: utils.transform.toCSV
    version: 1
    inputs:
      data: '${{ .steps.runQuery.outputs.results }}'
```

## Import

Workflow automations can be imported using the composite ID format: `<scope_type>#<scope_id>#<workflow_name>`, e.g.

```bash
$ terraform import newrelic_workflow_automation.test_query ACCOUNT#1234567#test_query_workflow
```

For workflows with complex names:

```bash
$ terraform import newrelic_workflow_automation.nrql_slack_report ACCOUNT#1234567#nrqlSlackReport
```

For organization-scoped workflows:

```bash
$ terraform import newrelic_workflow_automation.org_query ORGANIZATION#org-id#org_query_workflow
```

## Important Notes

### Name Consistency

The `name` field in the Terraform resource **must** match the `name` field in the YAML definition. If they don't match, Terraform will return an error during plan or apply.

For example, this configuration is **correct**:

```hcl
resource "newrelic_workflow_automation" "example" {
  name       = "my-workflow"
  scope_id   = "1234567"
  scope_type = "ACCOUNT"

  definition = <<-YAML
name: my-workflow      # This matches the resource name
description: Example workflow
steps:
  - name: waitStep
    type: wait
    seconds: 10
YAML
}
```

This configuration is **incorrect** and will fail:

```hcl
resource "newrelic_workflow_automation" "example" {
  name       = "my-workflow"
  scope_id   = "1234567"
  scope_type = "ACCOUNT"

  definition = <<-YAML
name: different-name   # This doesn't match the resource name - ERROR!
description: Example workflow
steps:
  - name: waitStep
    type: wait
    seconds: 10
YAML
}
```

### Scope Types

* **ACCOUNT** - The workflow automation is scoped to a specific New Relic account. Use your numeric account ID as the `scope_id`.
* **ORGANIZATION** - The workflow automation is scoped to your entire New Relic organization. Use your organization ID string as the `scope_id`.

### ForceNew Attributes

The following attributes will force a new resource to be created if changed:
* `name` - Changing the workflow name creates a new workflow.
* `scope_id` - Changing the scope ID creates a new workflow.
* `scope_type` - Changing between ACCOUNT and ORGANIZATION creates a new workflow.

### YAML Validation

The provider validates the YAML definition during plan and apply operations:
* The YAML must be valid and parseable.
* The `name` field must be present in the YAML.
* The `name` in the YAML must match the Terraform resource `name`.

Invalid YAML or missing required fields will result in an error.

## Additional Information

For more details about New Relic Workflow Automation, please refer to the [New Relic Workflow Automation documentation](https://docs.newrelic.com/docs/workflow-automation/).

### Versioning

Each time you update the `definition` of a workflow automation, New Relic automatically increments the `version` attribute. This allows you to track changes to your workflow automation over time.

### Best Practices

1. **Use Heredoc Syntax**: For multi-line YAML definitions, use the heredoc syntax (`<<-YAML ... YAML`) for better readability.

2. **External YAML Files**: For complex workflows, consider storing your YAML in separate files and using Terraform's `file()` or `templatefile()` functions:

   ```hcl
   resource "newrelic_workflow_automation" "from_file" {
     name       = "workflow-from-file"
     scope_id   = var.account_id
     scope_type = "ACCOUNT"

     definition = file("${path.module}/workflows/my-workflow.yaml")
   }
   ```

3. **Version Control**: Store your workflow YAML definitions in version control alongside your Terraform configuration.

4. **Testing**: Test workflow automation changes in a non-production account before applying to production.

5. **Naming Conventions**: Use consistent naming conventions for your workflows to make them easier to manage and identify.

### Troubleshooting

#### Name Mismatch Error

If you receive an error like "name in resource configuration does not match name in YAML definition", ensure that:
* The `name` attribute in your Terraform resource matches exactly with the `name` field in your YAML definition.
* There are no extra spaces or different capitalization between the two names.

#### Scope ID Format Error

If you receive an error about invalid scope_id format for ACCOUNT scope:
* Ensure your account ID is numeric (e.g., "1234567", not "account-1234567").
* For ACCOUNT scope, the scope_id should be a string representation of your numeric account ID.

#### Invalid YAML Error

If you receive a YAML parsing error:
* Validate your YAML syntax using a YAML validator.
* Ensure proper indentation (YAML is indentation-sensitive).
* Check that all required fields are present.

## See Also

* [New Relic Workflow Automation Documentation](https://docs.newrelic.com/docs/workflow-automation/)