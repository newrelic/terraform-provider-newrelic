package newrelic

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

func resourceNewRelicWorkload() *schema.Resource {
	userReference := &schema.Schema{
		Type:     schema.TypeList,
		MaxItems: 1,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"email": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"gravatar": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"id": {
					Type:     schema.TypeInt,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}

	return &schema.Resource{
		Create: resourceNewRelicWorkloadCreate,
		Read:   resourceNewRelicWorkloadRead,
		Update: resourceNewRelicWorkloadUpdate,
		Delete: resourceNewRelicWorkloadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The account the workload belongs to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The moment when the object was created, represented in milliseconds since the Unix epoch.",
			},
			"created_by": userReference,
			"entity": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "A list of entities manually assigned to this workload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"guid": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique entity identifier in New Relic.",
						},
					},
				},
			},
			"entity_search_query": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "A list of search queries that define a dynamic workload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the query.",
						},
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The query.",
						},
					},
				},
			},
			"guid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique entity identifier of the workload in New Relic.",
			},
			"workload_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The unique entity identifier of the workload.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The workload's name.",
			},
			"permalink": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the workload.",
			},
			"scope_accounts": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "Accounts that will be used to get entities from.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_ids": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
							MinItems:    1,
							Description: "A list of accounts that will be used to get entities from.",
						},
					},
				},
			},
		},
	}
}

func resourceNewRelicWorkloadCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	workload := expandWorkload(d)
	accountID := d.Get("account.0.id").(int)

	log.Printf("[INFO] Creating New Relic One workload %s", workload.Name)

	createInput := workloads.CreateInput{
		Name: workload.Name,
	}

	if len(workload.Entities) > 0 {
		var entityGUIDs []string
		for _, e := range workload.Entities {
			entityGUIDs = append(entityGUIDs, *e.GUID)
		}

		createInput.EntityGUIDs = entityGUIDs
	}

	if len(workload.EntitySearchQueries) > 0 {
		var entitySearchQueries []workloads.EntitySearchQueryInput
		for _, q := range workload.EntitySearchQueries {
			queryInput := workloads.EntitySearchQueryInput{
				Name:  &q.Name,
				Query: q.Query,
			}

			entitySearchQueries = append(entitySearchQueries, queryInput)
		}

		createInput.EntitySearchQueries = entitySearchQueries
	}

	if len(workload.ScopeAccounts.AccountIDs) > 0 {
		var scopeAccounts workloads.ScopeAccountsInput
		scopeAccounts.AccountIDs = append(scopeAccounts.AccountIDs, workload.ScopeAccounts.AccountIDs...)

		createInput.ScopeAccountsInput = &scopeAccounts
	}

	created, err := client.Workloads.CreateWorkload(accountID, createInput)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d:%d:%s", accountID, created.ID, created.GUID))
	return resourceNewRelicWorkloadRead(d, meta)
}

func resourceNewRelicWorkloadRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return err
	}

	workload, err := client.Workloads.GetWorkload(ids.AccountID, ids.ID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenWorkload(workload, d)
}

func resourceNewRelicWorkloadUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	workload := expandWorkload(d)

	log.Printf("[INFO] Updating New Relic One workload %s", d.Id())

	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return err
	}

	updateInput := workloads.UpdateInput{
		Name: &workload.Name,
	}

	if len(workload.Entities) > 0 {
		var entityGUIDs []string
		for _, e := range workload.Entities {
			entityGUIDs = append(entityGUIDs, *e.GUID)
		}

		updateInput.EntityGUIDs = entityGUIDs
	}

	if len(workload.EntitySearchQueries) > 0 {
		var entitySearchQueries []workloads.EntitySearchQueryInput
		for _, q := range workload.EntitySearchQueries {
			queryInput := workloads.EntitySearchQueryInput{
				Name:  &q.Name,
				Query: q.Query,
			}

			entitySearchQueries = append(entitySearchQueries, queryInput)
		}

		updateInput.EntitySearchQueries = entitySearchQueries
	}

	if len(workload.ScopeAccounts.AccountIDs) > 0 {
		var scopeAccounts workloads.ScopeAccountsInput
		scopeAccounts.AccountIDs = append(scopeAccounts.AccountIDs, workload.ScopeAccounts.AccountIDs...)

		updateInput.ScopeAccountsInput = &scopeAccounts
	}

	_, err = client.Workloads.UpdateWorkload(ids.GUID, updateInput)
	if err != nil {
		return err
	}

	return resourceNewRelicWorkloadRead(d, meta)
}

func resourceNewRelicWorkloadDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One workload %s", d.Id())

	ids, err := parseWorkloadIDs(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Workloads.DeleteWorkload(ids.GUID); err != nil {
		return err
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
