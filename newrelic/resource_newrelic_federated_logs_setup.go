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

func resourceNewRelicFederatedLogsSetup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFederatedLogsSetupCreate,
		ReadContext:   resourceNewRelicFederatedLogsSetupRead,
		UpdateContext: resourceNewRelicFederatedLogsSetupUpdate,
		DeleteContext: resourceNewRelicFederatedLogsSetupDelete,
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
			"scope_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The scope ID (account ID or organization ID) for the entity.",
			},
			"scope_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "ACCOUNT",
				Description: "The scope type: ACCOUNT or ORGANIZATION.",
			},
		},
	}
}

func resourceNewRelicFederatedLogsSetupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating New Relic Federated Log Setup: name=%s", d.Get("name").(string))

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	scopeType := pipelinecontrol.EntityManagementEntityScope(d.Get("scope_type").(string))
	scopeID := d.Get("scope_id").(string)
	if scopeID == "" {
		scopeID = strconv.Itoa(accountID)
	}

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
			Type: scopeType,
			ID:   scopeID,
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

	return resourceNewRelicFederatedLogsSetupRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsSetupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		if err := d.Set("scope_id", entity.Scope.ID); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("scope_type", string(entity.Scope.Type)); err != nil {
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

func resourceNewRelicFederatedLogsSetupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating New Relic Federated Log Setup: id=%s name=%s", d.Id(), d.Get("name").(string))

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	entityID := d.Id()

	resp, err := client.Pipelinecontrol.GetEntityWithContext(ctx, entityID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.FromErr(fmt.Errorf("entity with ID %s was nil", entityID))
	}

	entity, ok := (*resp).(*pipelinecontrol.EntityManagementFederatedLogSetupEntity)
	if !ok {
		return diag.Errorf("unexpected entity type for ID %s", entityID)
	}

	version := entity.Metadata.Version

	input := pipelinecontrol.EntityManagementFederatedLogSetupEntityUpdateInput{
		CloudProvider:              pipelinecontrol.EntityManagementCloudProvider(d.Get("cloud_provider").(string)),
		CloudProviderRegion:        d.Get("cloud_provider_region").(string),
		DataLocationBucket:         d.Get("data_location_bucket").(string),
		DataProcessingConnectionId: d.Get("data_processing_connection_id").(string),
		Description:                d.Get("description").(string),
		Name:                       d.Get("name").(string),
		NrAccountId:                d.Get("nr_account_id").(string),
		NrRegion:                   pipelinecontrol.EntityManagementNrRegion(d.Get("nr_region").(string)),
		QueryConnectionId:          d.Get("query_connection_id").(string),
		Status:                     pipelinecontrol.EntityManagementFederatedLogSetupStatus(d.Get("status").(string)),
	}

	if v, ok := d.GetOk("data_processing_component_id"); ok {
		input.DataProcessingComponentId = v.(string)
	}

	_, err = client.Pipelinecontrol.EntityManagementUpdateFederatedLogSetupWithContext(ctx, input, entityID, version)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicFederatedLogsSetupRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsSetupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic Federated Log Setup: id=%s", d.Id())

	_, err := client.Pipelinecontrol.EntityManagementDeleteWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
