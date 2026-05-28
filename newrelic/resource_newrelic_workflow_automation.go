package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workflowautomation"
	"gopkg.in/yaml.v3"
)

const (
	// Scope type constants for workflow automation
	scopeTypeAccount      = "ACCOUNT"
	scopeTypeOrganization = "ORGANIZATION"
)

func resourceNewRelicWorkflowAutomation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicWorkflowAutomationCreate,
		ReadContext:   resourceNewRelicWorkflowAutomationRead,
		UpdateContext: resourceNewRelicWorkflowAutomationUpdate,
		DeleteContext: resourceNewRelicWorkflowAutomationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the workflow automation. Must match the name in the YAML definition.",
			},
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The YAML definition of the workflow automation.",
			},
			"scope_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The scope ID (account ID for ACCOUNT scope, organization ID for ORGANIZATION scope).",
			},
			"scope_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The scope type. Supported values are: ACCOUNT, ORGANIZATION.",
			},
			"definition_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the workflow automation.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the workflow automation.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the workflow automation",
			},
			"yaml": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The yaml of the workflow automation.",
			},
		},
	}
}

func resourceNewRelicWorkflowAutomationValidateScopeType(scopeType string) error {
	// Scope type is required and must not be empty
	if scopeType == "" {
		return fmt.Errorf("scope_type is required and must be specified. Supported values are: %s, %s", scopeTypeAccount, scopeTypeOrganization)
	}

	// Check if scope type is one of the supported values
	if scopeType != scopeTypeAccount && scopeType != scopeTypeOrganization {
		return fmt.Errorf("scope_type '%s' is not supported. Supported values are: %s, %s", scopeType, scopeTypeAccount, scopeTypeOrganization)
	}

	return nil
}

// resourceNewRelicWorkflowAutomationValidateScope validates and returns the scope type and ID
func resourceNewRelicWorkflowAutomationValidateScope(d *schema.ResourceData) (scopeType, scopeID string, err error) {
	scopeType = d.Get("scope_type").(string)
	if err := resourceNewRelicWorkflowAutomationValidateScopeType(scopeType); err != nil {
		return "", "", err
	}

	scopeID = d.Get("scope_id").(string)
	if scopeID == "" {
		return "", "", fmt.Errorf("scope_id is required and must be specified")
	}

	return scopeType, scopeID, nil
}

func resourceNewRelicWorkflowAutomationParseNameFromYAML(yamlContent string) (string, error) {
	var workflow struct {
		Name string `yaml:"name"`
	}
	if err := yaml.Unmarshal([]byte(yamlContent), &workflow); err != nil {
		return "", fmt.Errorf("failed to parse YAML definition: %w", err)
	}
	if workflow.Name == "" {
		return "", fmt.Errorf("name field not found in YAML definition")
	}
	return workflow.Name, nil
}

// resourceNewRelicWorkflowAutomationValidateAndGetName validates that the name
// in the YAML definition matches the name specified in the resource configuration.
// The resource name field is the source of truth.
func resourceNewRelicWorkflowAutomationValidateAndGetName(d *schema.ResourceData) (string, error) {
	// Get name from Terraform resource config (required field, source of truth)
	name := d.Get("name").(string)

	// Get the YAML definition
	definition := d.Get("definition").(string)

	// Parse name from YAML to validate it matches
	yamlName, err := resourceNewRelicWorkflowAutomationParseNameFromYAML(definition)
	if err != nil {
		return "", err
	}

	// Validate that YAML name matches the resource name
	if name != yamlName {
		return "", fmt.Errorf("name in resource configuration (%s) does not match name in YAML definition (%s). The name field in your YAML must match the resource name", name, yamlName)
	}

	// Return the resource name (source of truth)
	return name, nil
}

func resourceNewRelicWorkflowAutomationParseID(id string) (scopeType string, scopeID string, name string, err error) {
	parts := strings.Split(id, "#")
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid resource ID format: expected '<scope_type>#<scope_id>#<workflow_name>', got: %s", id)
	}
	return parts[0], parts[1], parts[2], nil
}

func resourceNewRelicWorkflowAutomationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// Validate that resource name matches YAML definition name
	name, err := resourceNewRelicWorkflowAutomationValidateAndGetName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate and get scope
	scopeType, scopeID, err := resourceNewRelicWorkflowAutomationValidateScope(d)
	if err != nil {
		return diag.FromErr(err)
	}

	definition := d.Get("definition").(string)

	definitionInput := workflowautomation.WorkflowAutomationCreateWorkflowDefinitionInput{
		Yaml: workflowautomation.SecureValue(definition),
	}

	scopeInput := workflowautomation.WorkflowAutomationScopeInput{
		ID:   scopeID,
		Type: workflowautomation.WorkflowAutomationScopeType(scopeType),
	}

	createWorkflowAutomationResult, err := client.WorkflowAutomation.WorkflowAutomationCreateWorkflowDefinition(
		definitionInput,
		scopeInput,
		nil,
	)

	if err != nil {
		return diag.FromErr(err)
	}
	if createWorkflowAutomationResult == nil {
		return diag.FromErr(fmt.Errorf("error creating workflow automation"))
	}

	resourceNewRelicWorkflowAutomationSetValuesToState(d)
	log.Printf("[INFO] Created workflow automation %s (scope: %s/%s)", name, scopeType, scopeID)

	// Read the resource to populate computed fields like description and version
	return resourceNewRelicWorkflowAutomationRead(ctx, d, meta)
}

// setWorkflowStateFromAPI sets the workflow state from the API response
func setWorkflowStateFromAPI(d *schema.ResourceData, yaml, description string, version int, name, scopeType, scopeID string) {
	_ = d.Set("name", name)
	_ = d.Set("scope_type", scopeType)
	_ = d.Set("scope_id", scopeID)

	if yaml != "" {
		_ = d.Set("definition", yaml)
	}
	if description != "" {
		_ = d.Set("description", description)
	}
	if version > 0 {
		_ = d.Set("version", version)
	}

	resourceNewRelicWorkflowAutomationSetValuesToState(d)
}

func resourceNewRelicWorkflowAutomationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	parsedScopeType, parsedScopeID, parsedName, err := resourceNewRelicWorkflowAutomationParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	if name == "" {
		name = parsedName
	}

	scopeType := d.Get("scope_type").(string)
	if scopeType == "" {
		scopeType = parsedScopeType
	}

	scopeID := d.Get("scope_id").(string)
	if scopeID == "" {
		scopeID = parsedScopeID
	}

	if name == "" {
		return diag.FromErr(fmt.Errorf("workflow name is required but not found in state or resource ID"))
	}

	if scopeID == "" {
		return diag.FromErr(fmt.Errorf("scope_id is required but not found in state or resource ID"))
	}

	if err := resourceNewRelicWorkflowAutomationValidateScopeType(scopeType); err != nil {
		return diag.FromErr(err)
	}

	if scopeType == scopeTypeOrganization {
		workflow, err := client.Organization.GetWorkflow(name, 0)
		if err != nil || workflow == nil {
			log.Printf("[WARN] Workflow automation %s not found, removing from state", name)
			d.SetId("")
			return nil
		}

		setWorkflowStateFromAPI(d, string(workflow.Definition.Yaml), workflow.Definition.Description, workflow.Definition.Version, name, scopeType, scopeID)
	} else if scopeType == scopeTypeAccount {
		accountID, err := strconv.Atoi(scopeID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid scope_id format for %s scope (must be numeric): %s", scopeTypeAccount, scopeID))
		}

		workflow, err := client.WorkflowAutomation.GetWorkflow(accountID, name, 0)
		if err != nil || workflow == nil {
			log.Printf("[WARN] Workflow automation %s not found, removing from state", name)
			d.SetId("")
			return nil
		}

		setWorkflowStateFromAPI(d, string(workflow.Definition.Yaml), workflow.Definition.Description, workflow.Definition.Version, name, scopeType, strconv.Itoa(accountID))
	}

	return nil
}

func resourceNewRelicWorkflowAutomationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// Validate that resource name matches YAML definition name
	name, err := resourceNewRelicWorkflowAutomationValidateAndGetName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate and get scope
	scopeType, scopeID, err := resourceNewRelicWorkflowAutomationValidateScope(d)
	if err != nil {
		return diag.FromErr(err)
	}

	definitionInput := workflowautomation.WorkflowAutomationUpdateWorkflowDefinitionInput{
		Yaml: workflowautomation.SecureValue(d.Get("definition").(string)),
	}

	scopeInput := workflowautomation.WorkflowAutomationScopeInput{
		ID:   scopeID,
		Type: workflowautomation.WorkflowAutomationScopeType(scopeType),
	}

	tags := []workflowautomation.WorkflowAutomationTag{}

	updateResult, err := client.WorkflowAutomation.WorkflowAutomationUpdateWorkflowDefinition(
		definitionInput,
		scopeInput,
		tags,
	)

	if err != nil {
		return diag.FromErr(err)
	}
	if updateResult == nil {
		return diag.FromErr(fmt.Errorf("error updating workflow automation"))
	}

	log.Printf("[INFO] Updated workflow automation %s (scope: %s/%s)", name, scopeType, scopeID)

	// Read the resource to populate computed fields like description and version
	return resourceNewRelicWorkflowAutomationRead(ctx, d, meta)
}

func resourceNewRelicWorkflowAutomationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	// Validate that resource name matches YAML definition name
	name, err := resourceNewRelicWorkflowAutomationValidateAndGetName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Validate and get scope
	scopeType, scopeID, err := resourceNewRelicWorkflowAutomationValidateScope(d)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInput := workflowautomation.WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: name,
	}

	scopeInput := workflowautomation.WorkflowAutomationScopeInput{
		ID:   scopeID,
		Type: workflowautomation.WorkflowAutomationScopeType(scopeType),
	}

	_, deleteErr := client.WorkflowAutomation.WorkflowAutomationDeleteWorkflowDefinition(
		deleteInput,
		scopeInput,
	)

	if deleteErr != nil {
		return diag.FromErr(deleteErr)
	}

	d.SetId("")
	log.Printf("[INFO] Deleted workflow automation (scope: %s/%s)", scopeType, scopeID)
	return nil
}

func resourceNewRelicWorkflowAutomationSetValuesToState(
	d *schema.ResourceData,
) {
	name := d.Get("name").(string)
	scopeType := d.Get("scope_type").(string)
	scopeID := d.Get("scope_id").(string)

	resourceID := fmt.Sprintf("%s#%s#%s", scopeType, scopeID, name)
	d.SetId(resourceID)
}
