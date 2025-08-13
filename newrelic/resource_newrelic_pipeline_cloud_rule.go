package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entitymanagement"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"
	"log"
	"strconv"
)

func resourceNewRelicPipelineCloudRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicPipelineCloudRuleCreate,
		ReadContext:   resourceNewRelicPipelineCloudRuleRead,
		DeleteContext: resourceNewRelicPipelineCloudRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The account ID where the Pipeline Cloud rule will be created.",
			},
			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The name of the rule. This must be unique within an account.",
			},
			"nrql": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The NRQL query that defines which data will be processed by this pipeline cloud rule.",
			},
			"description": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Provides additional information about the rule.",
			},
			"entity_version": {
				Type:     schema.TypeInt,
				Computed: true,
				// Add this line to hide the attribute from the user
				Deprecated: "This attribute is internal and managed by the provider.",
			},
		},
	}
}

func resourceNewRelicPipelineCloudRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	description := d.Get("description").(string)
	nrql := d.Get("nrql")
	name := d.Get("name").(string)

	createInput := entitymanagement.EntityManagementPipelineCloudRuleEntityCreateInput{
		Description: description,
		NRQL:        nrdb.NRQL(nrql.(string)),
		Name:        name,
		Scope: entitymanagement.EntityManagementScopedReferenceInput{
			Type: entitymanagement.EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   strconv.Itoa(accountID),
		},
	}

	created, err := client.EntityManagement.EntityManagementCreatePipelineCloudRuleWithContext(ctx, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(created.Entity.ID)
	return resourceNewRelicPipelineCloudRuleRead(ctx, d, meta)
}

func resourceNewRelicPipelineCloudRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic Pipeline Cloud Rule for %s", d.Id())

	accountID := selectAccountID(providerConfig, d)
	ruleID := d.Id()

	resp, err := client.EntityManagement.GetEntityWithContext(ctx, ruleID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("entity with ruleID %s was nil", ruleID))
	}

	switch entityType := (*resp).(type) {
	case *entitymanagement.EntityManagementPipelineCloudRuleEntity:
		entity := (*resp).(*entitymanagement.EntityManagementPipelineCloudRuleEntity)

		d.SetId(ruleID)
		if err := d.Set("account_id", accountID); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("name", entity.Name); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("description", entity.Description); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("nrql", entity.NRQL); err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("entity_version", entity.Metadata.Version); err != nil {
			return diag.FromErr(err)
		}
	default:
		// This handles cases where the GUID belongs to a different type of New Relic entity.
		return diag.Errorf("unexpected entity type %T for ID %s", entityType, d.Id())
	}
	return nil
}

func resourceNewRelicPipelineCloudRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic Pipeline Cloud rule entity with rule id %s", d.Id())

	ruleID := d.Id()
	version := d.Get("entity_version").(int)

	result, err := client.EntityManagement.EntityManagementDelete(ruleID, version)
	if err != nil {
		return diag.FromErr(err)
	}

	if result.ID == "" {
		return diag.FromErr(fmt.Errorf("error in deleting entity with ruleID %s", ruleID))
	}
	return nil
}
