package newrelic

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/pipelinecontrol"
)

func resourceNewRelicFederatedLogSetup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFederatedLogSetupCreate,
		ReadContext:   resourceNewRelicFederatedLogSetupRead,
		UpdateContext: resourceNewRelicFederatedLogSetupUpdate,
		DeleteContext: resourceNewRelicFederatedLogSetupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The account ID where the federated log setup will be created.",
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_provider_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_location_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_processing_component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_processing_connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nr_account_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nr_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNewRelicFederatedLogSetupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	input := pipelinecontrol.EntityManagementFederatedLogSetupEntityCreateInput{
		CloudProvider:              pipelinecontrol.EntityManagementCloudProvider(d.Get("cloud_provider").(string)),
		CloudProviderRegion:        d.Get("cloud_provider_region").(string),
		DataLocationBucket:         d.Get("data_location_bucket").(string),
		DataProcessingConnectionId: d.Get("data_processing_connection_id").(string),
		Name:                       d.Get("name").(string),
		NrAccountId:                d.Get("nr_account_id").(string),
		NrRegion:                   pipelinecontrol.EntityManagementNrRegion(d.Get("nr_region").(string)),
		QueryConnectionId:          d.Get("query_connection_id").(string),
		Status:                     pipelinecontrol.EntityManagementFederatedLogSetupStatus(d.Get("status").(string)),
		Scope: pipelinecontrol.EntityManagementScopedReferenceInput{
			Type: pipelinecontrol.EntityManagementEntityScopeTypes.ACCOUNT,
			ID:   strconv.Itoa(accountID),
		},
	}

	if v, ok := d.GetOk("data_processing_component_id"); ok {
		input.DataProcessingComponentId = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	resp, err := client.Pipelinecontrol.EntityManagementCreateFederatedLogSetupWithContext(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	// handling the case when account_id is not specified in the configuration,
	// the API returns the account ID with which the entity has been created,
	// which we accordingly set as the account_id attribute(state)
	_ = d.Set("account_id", accountID)

	d.SetId(resp.Entity.ID)

	return resourceNewRelicFederatedLogSetupRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogSetupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic Federated Log Setup for %s", d.Id())

	entityID := d.Id()

	resp, err := client.Pipelinecontrol.GetEntityWithContext(ctx, entityID)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp == nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("entity with ID %s was nil", entityID))
	}

	switch entityType := (*resp).(type) {
	case *pipelinecontrol.EntityManagementFederatedLogSetupEntity:
		entity := (*resp).(*pipelinecontrol.EntityManagementFederatedLogSetupEntity)

		accountIDStr := entity.Scope.ID
		accountIDInt, accountIDIntErr := strconv.Atoi(accountIDStr)
		if accountIDIntErr != nil {
			log.Printf("[ERROR] Failed to convert account ID to integer: %v", accountIDIntErr)
			accountIDInt = selectAccountID(providerConfig, d)
			log.Printf("[INFO] Assigning the value of account_id from the state to prevent a panic: %d", accountIDInt)
		}

		if err := d.Set("account_id", accountIDInt); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("name", entity.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("description", entity.Description); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("cloud_provider", string(entity.CloudProvider)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("cloud_provider_region", entity.CloudProviderRegion); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("data_location_bucket", entity.DataLocationBucket); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nr_account_id", entity.NrAccountId); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nr_region", string(entity.NrRegion)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("status", string(entity.Status)); err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.Errorf("unexpected entity type %T for ID %s", entityType, d.Id())
	}
	return nil
}

func resourceNewRelicFederatedLogSetupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Implement update logic if needed
	return resourceNewRelicFederatedLogSetupRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogSetupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Implement delete logic if needed
	return nil
}
