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

func resourceNewRelicFederatedLogsPartition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFederatedLogsPartitionCreate,
		ReadContext:   resourceNewRelicFederatedLogsPartitionRead,
		UpdateContext: resourceNewRelicFederatedLogsPartitionUpdate,
		DeleteContext: resourceNewRelicFederatedLogsPartitionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The account ID where the federated log partition will be created.",
			},
			"data_location_uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URI location of the log partition in object storage.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the log partition.",
			},
			"is_default": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates if this log partition is the default partition.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the log partition.",
			},
			"nr_account_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The NR account ID associated with the federated log partition.",
			},
			"partition_database": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database name associated with the log partition.",
			},
			"partition_table": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The table name associated with the log partition.",
			},
			"retention_duration": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The duration value for retention.",
			},
			"retention_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The time unit for the retention duration (DAYS, WEEKS, MONTHS).",
			},
			"setup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The federated log setup this partition belongs to.",
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The status of the log partition (ACTIVE, CREATING, ERROR, INACTIVE).",
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

func resourceNewRelicFederatedLogsPartitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating New Relic Federated Log Partition: name=%s", d.Get("name").(string))

	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	accountID := selectAccountID(providerConfig, d)

	scopeType := pipelinecontrol.EntityManagementEntityScope(d.Get("scope_type").(string))
	scopeID := d.Get("scope_id").(string)
	if scopeID == "" {
		scopeID = strconv.Itoa(accountID)
	}

	input := pipelinecontrol.EntityManagementFederatedLogPartitionEntityCreateInput{
		DataLocationUri:   d.Get("data_location_uri").(string),
		IsDefault:         d.Get("is_default").(bool),
		Name:              d.Get("name").(string),
		NrAccountId:       d.Get("nr_account_id").(string),
		PartitionDatabase: d.Get("partition_database").(string),
		PartitionTable:    d.Get("partition_table").(string),
		SetupId:           d.Get("setup_id").(string),
		Status:            pipelinecontrol.EntityManagementLogPartitionStatus(d.Get("status").(string)),
		Scope: pipelinecontrol.EntityManagementScopedReferenceInput{
			Type: scopeType,
			ID:   scopeID,
		},
	}

	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}

	if v, ok := d.GetOk("retention_duration"); ok {
		input.RetentionPolicy = &pipelinecontrol.EntityManagementRetentionPolicyCreateInput{
			Duration: v.(int),
			Unit:     pipelinecontrol.EntityManagementRetentionUnit(d.Get("retention_unit").(string)),
		}
	}

	resp, err := client.Pipelinecontrol.EntityManagementCreateFederatedLogPartitionWithContext(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}

	_ = d.Set("account_id", accountID)

	d.SetId(resp.Entity.ID)

	return resourceNewRelicFederatedLogsPartitionRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsPartitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Reading New Relic Federated Log Partition for %s", d.Id())

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
	case *pipelinecontrol.EntityManagementFederatedLogPartitionEntity:
		entity := (*resp).(*pipelinecontrol.EntityManagementFederatedLogPartitionEntity)

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
		if err := d.Set("data_location_uri", entity.DataLocationUri); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("is_default", entity.IsDefault); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("nr_account_id", entity.NrAccountId); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("partition_database", entity.PartitionDatabase); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("partition_table", entity.PartitionTable); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("status", string(entity.Status)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("retention_duration", entity.RetentionPolicy.Duration); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("retention_unit", string(entity.RetentionPolicy.Unit)); err != nil {
			return diag.FromErr(err)
		}
	default:
		return diag.Errorf("unexpected entity type %T for ID %s", entityType, d.Id())
	}
	return nil
}

func resourceNewRelicFederatedLogsPartitionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating New Relic Federated Log Partition: id=%s name=%s", d.Id(), d.Get("name").(string))

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

	entity, ok := (*resp).(*pipelinecontrol.EntityManagementFederatedLogPartitionEntity)
	if !ok {
		return diag.Errorf("unexpected entity type for ID %s", entityID)
	}

	version := entity.Metadata.Version

	input := pipelinecontrol.EntityManagementFederatedLogPartitionEntityUpdateInput{
		DataLocationUri:   d.Get("data_location_uri").(string),
		Description:       d.Get("description").(string),
		IsDefault:         d.Get("is_default").(bool),
		Name:              d.Get("name").(string),
		NrAccountId:       d.Get("nr_account_id").(string),
		PartitionDatabase: d.Get("partition_database").(string),
		PartitionTable:    d.Get("partition_table").(string),
		SetupId:           d.Get("setup_id").(string),
		Status:            pipelinecontrol.EntityManagementLogPartitionStatus(d.Get("status").(string)),
	}

	if v, ok := d.GetOk("retention_duration"); ok {
		input.RetentionPolicy = &pipelinecontrol.EntityManagementRetentionPolicyUpdateInput{
			Duration: v.(int),
			Unit:     pipelinecontrol.EntityManagementRetentionUnit(d.Get("retention_unit").(string)),
		}
	}

	_, err = client.Pipelinecontrol.EntityManagementUpdateFederatedLogPartitionWithContext(ctx, input, entityID, version)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNewRelicFederatedLogsPartitionRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsPartitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Deleting New Relic Federated Log Partition: id=%s", d.Id())

	_, err := client.Pipelinecontrol.EntityManagementDeleteWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
