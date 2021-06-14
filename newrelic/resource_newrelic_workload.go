package newrelic

import (
	"context"
	"fmt"
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

func resourceNewRelicWorkloadCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	createInput := expandWorkloadCreateInput(d)
	accountID := d.Get("account_id").(int)

	log.Printf("[INFO] Creating New Relic One workload %s", createInput.Name)

	created, err := client.Workloads.CreateWorkloadWithContext(ctx, accountID, createInput)
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

	workload, err := client.Workloads.GetWorkloadWithContext(ctx, ids.AccountID, ids.GUID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	return diag.FromErr(flattenWorkload(workload, d))
}

func resourceNewRelicWorkloadUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandWorkloadUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One workload %s", d.Id())

	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.Workloads.UpdateWorkloadWithContext(ctx, ids.GUID, updateInput)
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

	if _, err := client.Workloads.DeleteWorkloadWithContext(ctx, ids.GUID); err != nil {
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
		GUID:      split[2],
	}, nil
}

type workloadIDs struct {
	AccountID int
	ID        int
	GUID      string
}

func (w *workloadIDs) String() string {
	return fmt.Sprintf("%d:%d:%s", w.AccountID, w.ID, w.GUID)
}
