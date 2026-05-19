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

func resourceNewRelicFederatedLogsSetup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFederatedLogsSetupCreate,
		ReadContext:   resourceNewRelicFederatedLogsSetupRead,
		UpdateContext: resourceNewRelicFederatedLogsSetupUpdate,
		DeleteContext: resourceNewRelicFederatedLogsSetupDelete,
		CustomizeDiff: validateFederatedLogsSetupDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the federated log setup.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the federated log setup.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether the setup is active. When false, log routing to this setup is turned off.",
			},
			"storage": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Storage configuration for this setup. Cannot be changed after creation.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_location_bucket": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object storage bucket where log data is stored.",
						},
						"database": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The database name associated with the federated log setup.",
						},
						"data_ingest_connection_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The connection manager entity GUID used for writing data.",
						},
						"query_connection_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The connection manager entity GUID used by query workers for reading data.",
						},
						"cloud_provider_configuration": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"provider": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice(
											[]string{string(federatedlogs.FederatedLogsCloudProviderTypes.AWS)},
											false,
										),
										Description: "The cloud provider. Currently only AWS is supported.",
									},
									"region": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The cloud provider region.",
									},
								},
							},
						},
					},
				},
			},
			"default_partition": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Default partition created alongside this setup. Logs that do not match any specific partition rule are routed here.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"table": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The table name associated with the default partition.",
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
					},
				},
			},
			"forwarder": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Forwarder configuration for processing and routing logs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								[]string{string(federatedlogs.FederatedLogsForwarderTypeTypes.PIPELINE_CONTROL)},
								false,
							),
							Description: "The type of forwarder. Currently only PIPELINE_CONTROL is supported.",
						},
						"pipeline_control": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fleet_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The fleet entity GUID used for deploying the pipeline configuration.",
									},
									"routing_rule": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"expression": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "OTTL expression for routing logs to this setup.",
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
			"default_partition_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default partition entity ID created alongside this setup.",
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

func resourceNewRelicFederatedLogsSetupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	input := federatedlogs.FederatedLogsCreateSetupInput{
		Name:             d.Get("name").(string),
		Storage:          expandFederatedLogsSetupStorage(d.Get("storage").([]interface{})),
		DefaultPartition: expandFederatedLogsDefaultPartition(d.Get("default_partition").([]interface{})),
	}
	if v, ok := d.GetOk("description"); ok {
		input.Description = v.(string)
	}
	if v, ok := d.GetOk("forwarder"); ok {
		input.Forwarder = expandFederatedLogsForwarder(v.([]interface{}))
	}

	log.Printf("[INFO] Creating New Relic Federated Logs Setup %s", input.Name)
	resp, err := client.Federatedlogs.FederatedLogsCreateSetupWithContext(ctx, input)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		return diag.Errorf("federatedLogsCreateSetup returned an empty response")
	}

	d.SetId(resp.Setup.ID)
	return resourceNewRelicFederatedLogsSetupRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsSetupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Federated Logs Setup %s", d.Id())
	resp, err := client.Federatedlogs.GetSetupWithContext(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	if resp == nil {
		d.SetId("")
		return nil
	}

	if err := flattenFederatedLogsSetupIntoState(d, resp); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNewRelicFederatedLogsSetupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	input := federatedlogs.FederatedLogsUpdateSetupInput{}
	if d.HasChange("name") {
		input.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		input.Description = d.Get("description").(string)
	}
	if d.HasChange("active") {
		input.Active = getBoolPointer(d.Get("active").(bool))
	}
	// The connection IDs are mutable per FederatedLogsUpdateSetupInput; bucket /
	// database / cloud_provider_configuration are immutable and are guarded
	// against accidental updates by the federatedLogsImmutableDiff CustomizeDiff.
	if d.HasChange("storage.0.data_ingest_connection_id") {
		input.DataIngestConnectionId = d.Get("storage.0.data_ingest_connection_id").(string)
	}
	if d.HasChange("storage.0.query_connection_id") {
		input.QueryConnectionId = d.Get("storage.0.query_connection_id").(string)
	}
	if d.HasChange("forwarder") {
		input.Forwarder = expandFederatedLogsForwarder(d.Get("forwarder").([]interface{}))
	}

	log.Printf("[INFO] Updating New Relic Federated Logs Setup %s", d.Id())
	if _, err := client.Federatedlogs.FederatedLogsUpdateSetupWithContext(ctx, d.Id(), input); err != nil {
		return diag.FromErr(err)
	}
	return resourceNewRelicFederatedLogsSetupRead(ctx, d, meta)
}

func resourceNewRelicFederatedLogsSetupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Per IDD, setup deletion is a soft-delete that transitions the entity to
	// the DELETING lifecycle state. The API cascades the DELETING state
	// to the default partition automatically; we don't need a separate
	// updatePartition call. Validation (no non-default partitions exist)
	// are handled server-side.
	client := meta.(*ProviderConfig).NewClient

	input := federatedlogs.FederatedLogsUpdateSetupInput{
		LifecycleStatus: federatedlogs.FederatedLogsLifecycleStatusInput{
			Status: federatedlogs.FederatedLogsLifecycleStateTypes.DELETING,
		},
	}
	if _, err := client.Federatedlogs.FederatedLogsUpdateSetupWithContext(ctx, d.Id(), input); err != nil {
		return diag.FromErr(fmt.Errorf("failed to mark setup %s as DELETING: %w", d.Id(), err))
	}

	return nil
}

// validateFederatedLogsSetupDiff guards fields that the wrapper API doesn't
// allow updating in place — i.e. fields accepted on FederatedLogsCreateSetupInput
// but absent from FederatedLogsUpdateSetupInput.
func validateFederatedLogsSetupDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	if d.Id() == "" {
		// Fresh create — every field "changes" from zero, which is fine.
		return nil
	}
	if d.HasChange("storage.0.data_location_bucket") {
		return immutableFieldError("storage.data_location_bucket", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("storage.0.database") {
		return immutableFieldError("storage.database", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("storage.0.cloud_provider_configuration.0.provider") {
		return immutableFieldError("storage.cloud_provider_configuration.provider", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("storage.0.cloud_provider_configuration.0.region") {
		return immutableFieldError("storage.cloud_provider_configuration.region", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("default_partition.0.storage.0.table") {
		return immutableFieldError("default_partition.storage.table", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("default_partition.0.storage.0.data_location_uri") {
		return immutableFieldError("default_partition.storage.data_location_uri", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("default_partition.0.data_retention_policy.0.duration") {
		return immutableFieldError("default_partition.data_retention_policy.duration", federatedLogsSetupRecreateHint)
	}
	if d.HasChange("default_partition.0.data_retention_policy.0.unit") {
		return immutableFieldError("default_partition.data_retention_policy.unit", federatedLogsSetupRecreateHint)
	}
	return nil
}
