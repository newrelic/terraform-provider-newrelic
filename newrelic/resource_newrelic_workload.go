package newrelic

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/common"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

func resourceNewRelicWorkload() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicWorkloadCreate,
		ReadContext:   resourceNewRelicWorkloadRead,
		UpdateContext: resourceNewRelicWorkloadUpdate,
		DeleteContext: resourceNewRelicWorkloadDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The New Relic account ID where you want to create the workload.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The workload's name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Relevant information about the workload.",
			},
			"entity_guids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "A list of entity GUIDs manually assigned to this workload.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"entity_search_query": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A list of search queries that define a dynamic workload.",
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The query.",
						},
					},
				},
			},
			"scope_account_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "A list of account IDs that will be used to get entities from.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"status_config_automatic": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "An input object used to represent an automatic status configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the automatic status configuration is enabled or not.",
						},
						"remaining_entities_rule_rollup": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "The input object used to represent a rollup strategy.",
							Elem:        WorkloadremainingEntitiesRuleSchemaElem(),
						},
						"rules": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "A list of rules.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entity_guids": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "A list of entity GUIDs composing the rule.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"nrql_query": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "A list of entity search queries used to retrieve the entities that compose the rule.",
										ForceNew:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"query": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The query.",
												},
											},
										},
									},
									"rollup": {
										Type:        schema.TypeSet,
										Required:    true,
										Description: "The input object used to represent a rollup strategy.",
										Elem:        WorkloadRuleRollupInputSchemaElem(),
									},
								},
							},
						},
					},
				},
			},
			"status_config_static": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "A list of static status configurations. You can only configure one static status for a workload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A description that provides additional details about the status of the workload.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Whether the static status configuration is enabled or not.",
						},
						"status": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The status of the workload.",
							ValidateFunc: validation.StringInSlice(listValidWorkloadStatuses(), false),
						},
						"summary": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A short description of the status of the workload.",
						},
					},
				},
			},

			"workload_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique entity identifier of the workload.",
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the workload in New Relic.",
			},
			"composite_entity_search_query": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The composite query used to compose a dynamic workload.",
			},
			"permalink": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the workload.",
			},
		},
	}
}

func WorkloadRuleRollupInputSchemaElem() *schema.Resource {
	s := WorkloadRollupInputSchemaElem()
	return &schema.Resource{
		Schema: s,
	}
}

func WorkloadremainingEntitiesRuleSchemaElem() *schema.Resource {
	s := WorkloadRollupInputSchemaElem()

	s["group_by"] = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The grouping to be applied to the remaining entities.",
		ValidateFunc: validation.StringInSlice(listValidWorkloadGroupBy(), false),
	}

	return &schema.Resource{
		Schema: s,
	}
}

func WorkloadRollupInputSchemaElem() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rollup": {
			Type:        schema.TypeSet,
			Required:    true,
			Description: "The input object used to represent a rollup strategy.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"strategy": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The rollup strategy that is applied to a group of entities.",
						ValidateFunc: validation.StringInSlice(listValidWorkloadStrategy(), false),
					},
					"threshold_type": {
						Type:         schema.TypeString,
						Optional:     true,
						Description:  "Type of threshold defined for the rule. This is an optional field that only applies when strategy is WORST_STATUS_WINS. Use a threshold to roll up the worst status only after a certain amount of entities are not operational.",
						ValidateFunc: validation.StringInSlice(listValidWorkloadRuleThresholdType(), false),
					},
					"threshold_value": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Threshold value defined for the rule. This optional field is used in combination with thresholdType. If the threshold type is null, the threshold value will be ignored.",
					},
				},
			},
		},
	}
}

func resourceNewRelicWorkloadCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	createInput := expandWorkloadCreateInput(d)
	accountID := d.Get("account_id").(int)

	log.Printf("[INFO] Creating New Relic One workload %s", createInput.Name)

	created, err := client.Workloads.WorkloadCreateWithContext(ctx, accountID, createInput)
	if err != nil {
		return diag.FromErr(err)
	}

	ids := workloadIDs{
		AccountID: accountID,
		ID:        created.ID,
		GUID:      created.GUID,
	}

	d.SetId(ids.String())
	return resourceNewRelicWorkloadRead(ctx, d, meta)
}

func resourceNewRelicWorkloadRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Entities.GetEntityWithContext(ctx, ids.GUID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	switch workload := (*resp).(type) {
	case *entities.WorkloadEntity:
		err := setWorkloadAttributes(d, map[string]string{
			"name": workload.Name,
		})
		if err != nil {
			diag.FromErr(err)
		}
	}
	return nil
}

func resourceNewRelicWorkloadUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandWorkloadUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One workload %s", d.Id())

	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Workloads.WorkloadUpdateWithContext(ctx, ids.GUID, updateInput)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(ids.String())

	return resourceNewRelicWorkloadRead(ctx, d, meta)
}

func resourceNewRelicWorkloadDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One workload %s", d.Id())

	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.Workloads.WorkloadDeleteWithContext(ctx, common.EntityGUID(ids.GUID)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func parseWorkloadIDs(ids string) (*workloadIDs, error) {
	split := strings.Split(ids, ":")

	accountID, err := strconv.ParseInt(split[0], 10, 32)
	if err != nil {
		return nil, err
	}

	workloadID, err := strconv.ParseInt(split[1], 10, 32)
	if err != nil {
		return nil, err
	}

	return &workloadIDs{
		AccountID: int(accountID),
		ID:        int(workloadID),
		GUID:      common.EntityGUID(split[2]),
	}, nil
}

type workloadIDs struct {
	AccountID int
	ID        int
	GUID      common.EntityGUID
}

func (w *workloadIDs) String() string {
	return fmt.Sprintf("%d:%d:%s", w.AccountID, w.ID, w.GUID)
}
