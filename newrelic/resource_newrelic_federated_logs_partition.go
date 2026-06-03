package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/federatedlogs"
)

func resourceNewRelicFederatedLogsPartition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFederatedLogsPartitionCreate,
		ReadContext:   resourceNewRelicFederatedLogsPartitionRead,
		UpdateContext: resourceNewRelicFederatedLogsPartitionUpdate,
		DeleteContext: resourceNewRelicFederatedLogsPartitionDelete,
		CustomizeDiff: validateFederatedLogsPartitionDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where the federated logs partition will live. Defaults to the provider's account_id. Changing this after creation is rejected by the API.",
			},
			"setup_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the federated log setup this partition belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the partition.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the partition.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the partition is active. When false, log routing to this partition is turned off.",
			},
			"storage": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Storage details for this partition. Cannot be changed after creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"table": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The table name associated with the partition.",
						},
						"data_location_uri": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URI location of the partition in object storage.",
						},
					},
				},
			},
			"data_retention_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"duration": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The duration value for retention.",
						},
						"unit": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{
									string(federatedlogs.FederatedLogsRetentionUnitTypes.DAYS),
									string(federatedlogs.FederatedLogsRetentionUnitTypes.WEEKS),
									string(federatedlogs.FederatedLogsRetentionUnitTypes.MONTHS),
								},
								false,
							),
							Description: "The time unit for the retention duration.",
						},
					},
				},
			},
			"forwarder_configuration": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Forwarder configuration for this partition. Type must match the parent setup's forwarder type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{string(federatedlogs.FederatedLogsForwarderTypeTypes.PIPELINE_CONTROL)},
								false,
							),
							Description: "The type of forwarder. Must match the parent setup's forwarder type.",
						},
						"pipeline_control": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"partition_rule": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expression": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "OTTL expression for routing logs to this partition.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"is_default": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether this is the default partition for the setup.",
			},
			"lifecycle_status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     statusDetailSchema(),
			},
			"health_check": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"last_updated_at":   {Type: schema.TypeString, Computed: true},
						"query_connection":  {Type: schema.TypeList, Computed: true, Elem: statusDetailSchema()},
						"end2end_data_flow": {Type: schema.TypeList, Computed: true, Elem: statusDetailSchema()},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNewRelicFederatedLogsPartitionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	input := federatedlogs.FederatedLogsCreatePartitionInput{
		Name:    d.Get("name").(string),
		Storage: expandFederatedLogsPartitionStorage(d.Get("storage").([]interface{})),
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("data_retention_policy"); ok {
		input.DataRetentionPolicy = expandFederatedLogsRetentionPolicy(v.([]interface{}))
	}
	if v, ok := d.GetOk("forwarder_configuration"); ok {
		input.ForwarderConfiguration = expandFederatedLogsPartitionForwarderConfig(v.([]interface{}))
	}

	setupID := d.Get("setup_id").(string)

	resp, err := client.Federatedlogs.FederatedLogsCreatePartitionWithContext(ctx, accountID, input, setupID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.Errorf("federatedLogsCreatePartition returned an empty response")
	}

	d.SetId(resp.Partition.ID)
	_ = d.Set("account_id", accountID)
	return resourceNewRelicFederatedLogsPartitionRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsPartitionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Reading New Relic Federated Logs Partition %s", d.Id())
	resp, err := client.Federatedlogs.GetPartitionWithContext(ctx, accountID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil || resp.ID == "" {
		log.Printf("[WARN] Federated Logs Partition %s not found", d.Id())
		d.SetId("")
		return nil
	}

	if err := flattenFederatedLogsPartitionIntoState(d, resp); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNewRelicFederatedLogsPartitionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	input := federatedlogs.FederatedLogsUpdatePartitionInput{}
	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		input.Description = d.Get("description").(string)
	}
	if d.HasChange("active") {
		input.Active = getBoolPointer(d.Get("active").(bool))
	}
	if d.HasChange("data_retention_policy") {
		input.DataRetentionPolicy = expandFederatedLogsRetentionPolicy(d.Get("data_retention_policy").([]interface{}))
	}
	if d.HasChange("forwarder_configuration") {
		input.ForwarderConfiguration = expandFederatedLogsPartitionForwarderConfig(d.Get("forwarder_configuration").([]interface{}))
	}

	log.Printf("[INFO] Updating New Relic Federated Logs Partition %s", d.Id())
	if _, err := client.Federatedlogs.FederatedLogsUpdatePartitionWithContext(ctx, accountID, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}
	return resourceNewRelicFederatedLogsPartitionRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsPartitionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Partition deletion is a soft-delete via the wrapper update
	// mutation, transitioning lifecycleStatus to DELETING.
	// The entity is not removed outright.
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	accountID := selectAccountID(providerConfig, d)

	input := federatedlogs.FederatedLogsUpdatePartitionInput{
		LifecycleStatus: &federatedlogs.FederatedLogsLifecycleStatusInput{
			Status: federatedlogs.FederatedLogsLifecycleStateTypes.DELETING,
		},
	}
	if _, err := client.Federatedlogs.FederatedLogsUpdatePartitionWithContext(ctx, accountID, d.Id(), input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to mark partition %s as DELETING: %w", d.Id(), err))
	}

	return nil
}

// validateFederatedLogsPartitionDiff guards fields that can't be updated in place.
func validateFederatedLogsPartitionDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.Id() == "" {
		return nil
	}
	if d.HasChange("setup_id") {
		return immutableFieldError("setup_id",
			"a partition is permanently tied to its parent setup. Recreate the resource to attach it to a different setup")
	}
	if d.HasChange("storage.0.table") {
		return immutableFieldError("storage.table", federatedLogsPartitionRecreateHint)
	}
	if d.HasChange("storage.0.data_location_uri") {
		return immutableFieldError("storage.data_location_uri", federatedLogsPartitionRecreateHint)
	}
	return nil
}
