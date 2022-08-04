package newrelic

import (
	"context"
	"fmt"
	"github.com/newrelic/newrelic-client-go/pkg/workflows"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/ai"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicWorkflowCreate,
		ReadContext:   resourceNewRelicWorkflowRead,
		UpdateContext: resourceNewRelicWorkflowUpdate,
		DeleteContext: resourceNewRelicWorkflowDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "(Required) The name of the workflow.",
			},
			"destination_configurations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Workflow's destination configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "(Required) Destination's channel id.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "(Required) Destination's name.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(listValidWorkflowsDestinationTypes(), false),
							Description:  fmt.Sprintf("(Required) The type of the destination. One of: (%s).", strings.Join(listValidWorkflowsDestinationTypes(), ", ")),
						},
					},
				},
			},
			"enrichments": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Workflow's destination configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "(Required) Enrichment's name.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(listValidWorkflowsEnrichmentTypes(), false),
							Description:  fmt.Sprintf("The type of the enrichment. One of: (%s).", strings.Join(listValidWorkflowsEnrichmentTypes(), ", ")),
						},
						"account_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: "The account id of the enrichment.",
						},
						"configurations": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A set of key-value pairs to represent a destination configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"query": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "NRQL query.",
									},
								},
							},
						},

						// Computed
						"enrichment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enrichment's id.",
						},
					},
				},
			},
			"issues_filter": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Workflow's destination configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "(Required) Filter's name.",
						},
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice(listValidWorkflowsFilterTypes(), false),
							Description:  fmt.Sprintf("(Required) The type of the filter. One of: (%s).", strings.Join(listValidWorkflowsFilterTypes(), ", ")),
						},
						"predicates": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Workflow's destination configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "(Required) predicate's attribute.",
									},
									"operator": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice(listValidWorkflowsOperatorTypes(), false),
										Description:  fmt.Sprintf("The type of the operator. One of: (%s).", strings.Join(listValidWorkflowsOperatorTypes(), ", ")),
									},
									"values": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of predicate values.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},

						// Computed
						"filter_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "filter id.",
						},
					},
				},
			},

			// Optional
			"workflow_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the workflow is enabled.",
			},
			"enrichments_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the enrichments are enabled.",
			},
			"destinations_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the destinations are enabled.",
			},
			"muting_rules_handling": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice(listValidMutingRulesTypes(), false),
				Description:  fmt.Sprintf("The type of the muting rule handling. One of: (%s).", strings.Join(listValidMutingRulesTypes(), ", ")),
			},
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The account id of the workflow.",
			},

			// Computed
			"last_run": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time notification was sent for this workflow.",
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the workflow.",
			},
		},
	}
}

func resourceNewRelicWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	workflowInput, err := expandWorkflow(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic workflow %s", workflowInput.Name)

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	workflowResponse, err := client.Workflows.AiWorkflowsCreateWorkflowWithContext(updatedContext, accountID, *workflowInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(workflowResponse.Workflow.ID)

	return resourceNewRelicWorkflowRead(updatedContext, d, meta)
}

func resourceNewRelicWorkflowRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic workflow id: %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	filters := ai.AiWorkflowsFilters{ID: d.Id()}
	updatedContext := updateContextWithAccountID(ctx, accountID)

	workflowResponse, err := client.Workflows.GetWorkflowsWithContext(updatedContext, accountID, "", filters)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenWorkflow(&workflowResponse.Entities[0], d))
}

func resourceNewRelicWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput, err := expandWorkflowUpdate(d)
	if err != nil {
		return diag.FromErr(err)
	}

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	_, err = client.Workflows.AiWorkflowsUpdateWorkflowWithContext(updatedContext, accountID, *updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicWorkflowRead(updatedContext, d, meta)
}

func resourceNewRelicWorkflowDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic workflow %v", d.Id())

	providerConfig := meta.(*ProviderConfig)
	accountID := selectAccountID(providerConfig, d)
	updatedContext := updateContextWithAccountID(ctx, accountID)

	if _, err := client.Workflows.AiWorkflowsDeleteWorkflowWithContext(updatedContext, accountID, d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Validation function to validate allowed muting rules types
func listValidMutingRulesTypes() []string {
	return []string{
		string(workflows.AiWorkflowsMutingRulesHandlingTypes.NOTIFY_ALL_ISSUES),
		string(workflows.AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_MUTED_ISSUES),
		string(workflows.AiWorkflowsMutingRulesHandlingTypes.DONT_NOTIFY_FULLY_OR_PARTIALLY_MUTED_ISSUES),
	}
}

// Validation function to validate allowed destination types
func listValidWorkflowsDestinationTypes() []string {
	return []string{
		string(workflows.AiWorkflowsDestinationTypeTypes.EMAIL),
		string(workflows.AiWorkflowsDestinationTypeTypes.EVENT_BRIDGE),
		string(workflows.AiWorkflowsDestinationTypeTypes.PAGERDUTY_ACCOUNT_INTEGRATION),
		string(workflows.AiWorkflowsDestinationTypeTypes.PAGERDUTY_SERVICE_INTEGRATION),
		string(workflows.AiWorkflowsDestinationTypeTypes.SERVICE_NOW),
		string(workflows.AiWorkflowsDestinationTypeTypes.WEBHOOK),
		string(workflows.AiWorkflowsDestinationTypeTypes.MOBILE_PUSH),
		string(workflows.AiWorkflowsDestinationTypeTypes.SLACK),
		string(workflows.AiWorkflowsDestinationTypeTypes.JIRA),
	}
}

// Validation function to validate allowed enrichment types
func listValidWorkflowsEnrichmentTypes() []string {
	return []string{
		string(workflows.AiWorkflowsEnrichmentTypeTypes.NRQL),
	}
}

// Validation function to validate allowed filter types
func listValidWorkflowsFilterTypes() []string {
	return []string{
		string(workflows.AiWorkflowsFilterTypeTypes.FILTER),
		string(workflows.AiWorkflowsFilterTypeTypes.VIEW),
	}
}

// Validation function to validate allowed operator types
func listValidWorkflowsOperatorTypes() []string {
	return []string{
		string(workflows.AiWorkflowsOperatorTypes.CONTAINS),
		string(workflows.AiWorkflowsOperatorTypes.DOES_NOT_CONTAIN),
		string(workflows.AiWorkflowsOperatorTypes.DOES_NOT_EQUAL),
		string(workflows.AiWorkflowsOperatorTypes.DOES_NOT_EXACTLY_MATCH),
		string(workflows.AiWorkflowsOperatorTypes.ENDS_WITH),
		string(workflows.AiWorkflowsOperatorTypes.EQUAL),
		string(workflows.AiWorkflowsOperatorTypes.EXACTLY_MATCHES),
		string(workflows.AiWorkflowsOperatorTypes.GREATER_OR_EQUAL),
		string(workflows.AiWorkflowsOperatorTypes.GREATER_THAN),
		string(workflows.AiWorkflowsOperatorTypes.IS),
		string(workflows.AiWorkflowsOperatorTypes.IS_NOT),
		string(workflows.AiWorkflowsOperatorTypes.LESS_OR_EQUAL),
		string(workflows.AiWorkflowsOperatorTypes.LESS_THAN),
		string(workflows.AiWorkflowsOperatorTypes.STARTS_WITH),
	}
}
