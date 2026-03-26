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
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The name of the workflow automation. If not specified, it will be extracted from the YAML definition.",
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
		return fmt.Errorf("scope_type is required and must be specified. Supported values are: ACCOUNT, ORGANIZATION")
	}

	// Check if scope type is one of the supported values
	if scopeType != "ACCOUNT" && scopeType != "ORGANIZATION" {
		return fmt.Errorf("scope_type '%s' is not supported. Supported values are: ACCOUNT, ORGANIZATION", scopeType)
	}

	return nil
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

func resourceNewRelicWorkflowAutomationGetOrParseWorkflowName(d *schema.ResourceData) (string, error) {
	// Get name from Terraform config if provided
	name := d.Get("name").(string)

	// Get the YAML definition
	definition := d.Get("definition").(string)

	// Parse name from YAML
	parsedName, err := resourceNewRelicWorkflowAutomationParseNameFromYAML(definition)
	if err != nil {
		return "", err
	}

	// If name is provided in Terraform config, validate it matches the YAML
	if name != "" && name != parsedName {
		return "", fmt.Errorf("name in Terraform config (%s) does not match name in YAML definition (%s)", name, parsedName)
	}

	// Use the parsed name from YAML and set it in the state
	_ = d.Set("name", parsedName)

	return parsedName, nil
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

	// Parse and validate the workflow name from YAML
	name, err := resourceNewRelicWorkflowAutomationGetOrParseWorkflowName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	definition := d.Get("definition").(string)

	// Get and validate scope type (required field)
	scopeType := d.Get("scope_type").(string)
	if validationErr := resourceNewRelicWorkflowAutomationValidateScopeType(scopeType); validationErr != nil {
		return diag.FromErr(validationErr)
	}

	// Get and validate scope ID (required field)
	scopeID := d.Get("scope_id").(string)
	if scopeID == "" {
		return diag.FromErr(fmt.Errorf("scope_id is required and must be specified"))
	}

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

	if scopeType == "ORGANIZATION" {
		workflow, err := client.Organization.GetWorkflow(name, 0)
		if err != nil || workflow == nil {
			log.Printf("[WARN] Workflow automation %s not found, removing from state", name)
			d.SetId("")
			return nil
		}

		_ = d.Set("name", name)
		_ = d.Set("scope_type", scopeType)
		_ = d.Set("scope_id", scopeID)

		if workflow.Definition.Yaml != "" {
			_ = d.Set("definition", string(workflow.Definition.Yaml))
		}
		if workflow.Definition.Description != "" {
			_ = d.Set("description", workflow.Definition.Description)
		}
		if workflow.Definition.Version > 0 {
			_ = d.Set("version", workflow.Definition.Version)
		}

		resourceNewRelicWorkflowAutomationSetValuesToState(d)
	} else if scopeType == "ACCOUNT" {
		accountID, err := strconv.Atoi(scopeID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("invalid scope_id format for ACCOUNT scope (must be numeric): %s", scopeID))
		}

		workflow, err := client.WorkflowAutomation.GetWorkflow(accountID, name, 0)
		if err != nil || workflow == nil {
			log.Printf("[WARN] Workflow automation %s not found, removing from state", name)
			d.SetId("")
			return nil
		}

		_ = d.Set("name", name)
		_ = d.Set("scope_type", scopeType)
		_ = d.Set("scope_id", strconv.Itoa(accountID))

		if workflow.Definition.Yaml != "" {
			_ = d.Set("definition", string(workflow.Definition.Yaml))
		}
		if workflow.Definition.Description != "" {
			_ = d.Set("description", workflow.Definition.Description)
		}
		if workflow.Definition.Version > 0 {
			_ = d.Set("version", workflow.Definition.Version)
		}

		resourceNewRelicWorkflowAutomationSetValuesToState(d)
	}

	return nil
}

func resourceNewRelicWorkflowAutomationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	name, err := resourceNewRelicWorkflowAutomationGetOrParseWorkflowName(d)
	if err != nil {
		return diag.FromErr(err)
	}

	scopeType := d.Get("scope_type").(string)
	if validationErr := resourceNewRelicWorkflowAutomationValidateScopeType(scopeType); validationErr != nil {
		return diag.FromErr(validationErr)
	}

	scopeID := d.Get("scope_id").(string)
	if scopeID == "" {
		return diag.FromErr(fmt.Errorf("scope_id is required and must be specified"))
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

	scopeType := d.Get("scope_type").(string)
	if err := resourceNewRelicWorkflowAutomationValidateScopeType(scopeType); err != nil {
		return diag.FromErr(err)
	}

	scopeID := d.Get("scope_id").(string)
	if scopeID == "" {
		return diag.FromErr(fmt.Errorf("scope_id is required and must be specified"))
	}

	deleteInput := workflowautomation.WorkflowAutomationDeleteWorkflowDefinitionInput{
		Name: d.Get("name").(string),
	}

	scopeInput := workflowautomation.WorkflowAutomationScopeInput{
		ID:   scopeID,
		Type: workflowautomation.WorkflowAutomationScopeType(scopeType),
	}

	_, err := client.WorkflowAutomation.WorkflowAutomationDeleteWorkflowDefinition(
		deleteInput,
		scopeInput,
	)

	if err != nil {
		return diag.FromErr(err)
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
