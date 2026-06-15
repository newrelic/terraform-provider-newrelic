---
layout: "newrelic"
page_title: "New Relic: newrelic_workflow_automation"
sidebar_current: "docs-newrelic-resource-workflow-automation"
description: |-
Create and manage workflow automation in New Relic.
---

# Resource: newrelic\_workflow\_automation

Use this resource to create and manage New Relic Workflow Automation.

Workflow Automation allows you to define automated workflows using YAML definitions. These workflows can scope to either an account or an organization and support various automation steps and configurations.

## **Example Usage**

### **Basic Workflow Automation with ACCOUNT Scope**

```
resource "newrelic_workflow_automation" "test_query" {
  name        = "test_query_workflow"
  scope_id    = "your-account-id"
  scope_type  = "ACCOUNT"
  definition  = <<-YAML
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

### **Workflow Automation with ORGANIZATION Scope**

```
resource "newrelic_workflow_automation" "org_query" {
  name        = "org_query_workflow"
  scope_id    = "your-organization-id"
  scope_type  = "ORGANIZATION"
  definition  = <<-YAML
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

## **Argument Reference**

This resource supports the following arguments:

* **`name`** \- (Required) The name of the workflow. This must match the `name` field in the YAML `definition`. **Important**: Changing this field will force a new resource to be created.
* **`definition`** \- (Required) The YAML definition of the workflow. This must be a valid YAML string that defines the workflow's configuration.
* **`scope_type`** \- (Required) The scope type for the workflow. Must be either `ACCOUNT` or `ORGANIZATION`. **Important**: Changing this field will force a new resource to be created.
* **`scope_id`** \- (Required) The ID of the scope for the workflow. For `ACCOUNT` scope, this is your New Relic account ID (numeric). For `ORGANIZATION` scope, this is your organization ID (string). **Important**: Changing this field will force a new resource to be created.

## **Attributes Reference**

In addition to the arguments above, the resource exports the following attributes:

* **`id`** \- The composite ID of the workflow, with the format `<scope_type>#<scope_id>#<workflow_name>`.
  * Example:
    * ACCOUNT\#123456789\#MyWorkflow
    * ORGANIZATION\#c400b54f-2abd-45b7-987c-dc1920ce701d\#MyWorkflow
* **`description`** \- The description of the workflow, as defined in the YAML.
* **`version`** \- The current version number of the workflow. This number increments with each update to the `definition`.

## [**YAML Definition Structure**](https://docs.newrelic.com/docs/workflow-automation/workflow-automation-apis/definition-schema/)

The `definition` argument accepts a YAML string that defines the workflow's structure and logic.

### **Top-Level Fields**

| Field | Type | Required | Description |
| ----- | ----- | ----- | ----- |
| `name` | String | **Yes** | The name of the workflow. Must match the `name` argument in the Terraform resource. |
| `description` | String | No | A brief summary of what the workflow does. Recommended for clarity. |
| `workflowInputs` | Object | No | A map of input variables that can be passed to the workflow at runtime. |
| `steps` | Array | **Yes** | An ordered array of step objects that define the workflow's logic. |

---

### **Workflow Inputs**

Define inputs to make your workflows more flexible and reusable.

```
workflowInputs:
  inputName:
    type: String       # Supported types include: String, Int, Boolean, List, Map
    defaultValue: "value"  # Optional default value for the input
    required: false    # Optional, defaults to true. If true, a value must be provided.
    validations:       # Optional array of validation rules
      - type: regex
        errorMessage: "Custom error message"
        pattern: ^[a-zA-Z0-9]+$
```

#### Inputs can reference secrets using the syntax: `${{ :secrets:secretName }}`

#### **Input Validation Types**

Workflow inputs support various validation types to ensure data integrity:

##### [Regex Validation:](https://docs.newrelic.com/docs/workflow-automation/workflow-automation-apis/definition-schema/#validation-types)

```
emailDestinationId:
  type: String
  validations:
    - type: regex
      errorMessage: "Must be a valid UUID"
      pattern: ^[a-fA-F0-9]{8}-([a-fA-F0-9]{4}-){3}[a-fA-F0-9]{12}$
```

Integer Range Validation:

```
threshold:
  type: Int
  defaultValue: 5000
  validations:
    - type: minIntValue
      errorMessage: "Minimum value must be at least 100"
      minValue: 100
    - type: maxIntValue
      errorMessage: "Maximum value must be less than 10000"
      maxValue: 10000
```

**Note**: By default, integer variables accept both positive and negative values. If you define a *minIntValue* validation, the field rejects any value below that threshold. For example, setting *minIntValue: 0* prevents negative integers from being entered.

List Type:

```
skipUsers:
  type: List
  defaultValue:
    - "user1@example.com"
    - "user2@example.com"
```

**Note**: This list is case sensitive.

**Step Types**

Each object in the `steps` array defines a single unit of work. The primary step types are [`action`, `wait`, `switch`, and `loop`.](https://docs.newrelic.com/docs/workflow-automation/create-a-workflow-automation/create-your-own/#core-concepts)

#### **Action Steps**

Executes a specific function, such as querying data or sending a notification.

* `type: action` \- Defines the step as an action.
  * `action` \- The specific action to execute (example, `newrelic.nrdb.query`).
  * `version` \- The version of the action to use.
  * `inputs` \- A map of input parameters for the action.

Example:

```
  - name: sendLogs
    type: action
    action: newrelic.ingest.sendLogs
    version: '1'
    inputs:
      logs:
        - message: "This is a test log from Workflow Automation"
```

**Common actions include:**

* `newrelic.nrdb.query` \- Query New Relic database
* `utils.transform.toCSV` \- Transform data to CSV format
* `slack.chat.postMessage` \- Send messages to Slack

**Available actions**
A complete list of available actions, their versions, and required inputs can be found in the [**Workflow Action Catalog**](https://docs.newrelic.com/docs/workflow-automation/actions-catalog/).

#### **Wait steps**

Pauses the workflow for a specified duration or until it receives a signal.

* `type: wait`
  * `seconds` \- The number of seconds to pause.
  * `signals` \- (Optional) An array of signals to wait for. See [SignalWorkflowRun](https://docs.newrelic.com/docs/workflow-automation/workflow-automation-apis/signal-workflow-run/) for more information
  * `next` \- (Optional) The name of the next step. Use `end` to terminate the workflow.

Example:

```
 - name: wait
    type: wait
    seconds: 15
```

#### **Switch steps**

Provides conditional branching logic.

* `type: switch`
  * `switch` \- An array of conditions to evaluate in order.
  * `condition` \- A JQ expression that evaluates to `true` or `false`.
  * `next` \- The name of the step to execute if the condition is true. The `switch` block can also have a top-level `next` field to define the default step if no conditions match.

Example:

```
   - name: hasCompleted
    type: switch
    switch:
      - condition: ${{ .steps.waitForCompletion.outputs.status == "Failed" }}
        next: displayError
      - condition: ${{ .steps.waitForCompletion.outputs.status == "Success" }}
        next: displaySuccess
```

#### **Loop steps**

Iterates over a set of values and executes a sequence of steps for each iteration.

* `type: loop`
  * `for.in` \- A JQ expression that returns an array to iterate over.
  * `steps` \- An array of steps to execute for each iteration. Inside the loop, you can use `next: continue` to skip to the next iteration or `next: break` to exit the loop.

Example:

```
- name: loopStep
    type: loop
    for:
      in: ${{ [range(0; 5)] }}
      steps:
        - name: sendLogs
          type: action
          action: newrelic.ingest.sendLogs
          version: '1'
          inputs:
            logs:
              - message: "This is a test log from Workflow Automation"
          next: continue
```

---

### **Referencing Data in Workflows**

You can dynamically reference data from inputs, secrets, and other steps using JQ-like template expressions.

| Data Source | Syntax | Example |
| ----- | ----- | ----- |
| **Workflow Inputs** | `${{ .workflowInputs.inputName }}` | `${{ .workflowInputs.ccuThreshold }}` |
| **Step Outputs** | `${{ .steps.stepName.outputs.fieldName }}` | `${{ .steps.query1.outputs.results }}` |
| **Loop Elements** | `${{ .steps.loopStepName.loop.element }}` | `${{ .steps.loopStep1.loop.element.email }}` |
| **Secrets** | `${{ :secrets:secretName }}` | `${{ :secrets:myApiKey }}` |

---

### Example YAML Structure

Simple workflow with query and wait:

```
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

Advanced workflow with inputs and multiple actions:

```
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

Import workflow automations using the composite ID format: `<scope_type>#<scope_id>#<workflow_name>`, for example:

```
$ terraform import newrelic_workflow_automation.test_query ACCOUNT#1234567#test_query_workflow
```

```
$ terraform import newrelic_workflow_automation.ccu_governance ACCOUNT#1234567#CCUGovernance
```

For organization-scoped workflows:

```
$ terraform import newrelic_workflow_automation.org_query ORGANIZATION#org-id#org_query_workflow
```

## Important note

### Name Consistency

The `name` field in the Terraform resource must match the name field in the YAML definition. If they don't match, Terraform will return an error during `terraform validate`, `plan`, or `apply`.

Example Error Message:

```
Error: name in resource configuration (my-workflow) does not match name in YAML definition (my-new-workflow). The name field in your YAML must match the resource name
```

For example, this configuration is **correct**:

```
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

```
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

### Scope type

* **ACCOUNT** \- The workflow automation is scoped to a specific New Relic account. Use your numeric account ID as the `scope_id`.
* **ORGANIZATION** \- The workflow automation is scoped to your entire New Relic organization. Use your organization ID string as the `scope_id`.

See, [Create accounts and organizations](https://docs.newrelic.com/docs/accounts/accounts-billing/account-structure/multi-tenancy/org-creation/) on how to create an account or an organization.

### ForceNew attributes

The following attributes, when changed, will force creation of a new resource :

* **`name`** \- Changing the workflow name creates a new resource.
* **`scope_id`** \- Changing the scope ID creates a new resource.
* **`scope_type`** \- Changing between **ACCOUNT** and **ORGANIZATION** creates a new resource.

### YAML validation

The provider validates the YAML definition during plan-and-apply operations:

* The YAML must be valid and parsable.
* The `name` field must be present in the YAML.
* The `name` in the YAML must match the Terraform resource name.

Invalid YAML or missing required fields will result in an error.

Example YAML validation errors:

  YAML Validation Error
  1\. *waitStep* has invalid type "waitAgain". Valid types are:
   action, loop, switch, wait, assign

###   2\. Workflow definition names can not be changed.

### *Versioning*

Each time you update the `definition` of a workflow automation, New Relic automatically increments the `version` attribute. This allows you to track changes to your workflow automation over time.

## **Best practices**

1. **Use Heredoc Syntax**: For multi-line YAML, use the `<<-YAML ... YAML` heredoc syntax for better readability.
2. **External YAML Files**: For complex workflows, store YAML in separate files and reference them using Terraform's `file()` or `templatefile()` function.

```
   resource "newrelic_workflow_automation" "from_file" {
     name       = "workflow-from-file"
     scope_id   = var.account_id
     scope_type = "ACCOUNT"

     definition = file("${path.module}/workflows/my-workflow.yaml")
   }
```

3. **Version Control**: Store workflow YAML definitions in a version control system (like Git) alongside your Terraform code.
4. **Testing**: Always test workflow changes in a non-production environment before applying them to production.
5. **Naming Conventions**: Use clear and consistent naming conventions for your workflows to make them easier to manage.
6. **Manage Secrets Securely**: For sensitive values like API keys or tokens, always use [New Relic secrets](https://docs.newrelic.com/docs/infrastructure/host-integrations/installation/secrets-management/). Avoid hardcoding sensitive information directly in your YAML definitions.

## **Troubleshooting**

#### **Name Mismatch Error**

If you get an error like "`name in resource configuration does not match name in YAML definition"`, ensure the `name` attribute in your Terraform resource exactly matches the `name` field in your YAML, including capitalization and spacing.

#### **Scope ID Format Error**

If you receive an error about invalid `scope_id` format for **ACCOUNT** scope:

* Ensure your `account ID` is numeric (e.g., "**1234567**", not "**account-1234567**").
* For **ACCOUNT** scope, the `scope_id` should be a string representation of your numeric account ID.

#### **Invalid YAML Error**

If you receive a YAML parsing error:

* Validate your YAML syntax using a YAML validator.
* Ensure proper indentation (YAML is indentation-sensitive).
* Check that all required fields are present.

## See also:

[New Relic Workflow Automation Documentation](https://docs.newrelic.com/docs/workflow-automation/)