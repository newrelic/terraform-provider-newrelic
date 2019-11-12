package newrelic

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func resourceNewRelicDashboard() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicDashboardCreate,
		Read:   resourceNewRelicDashboardRead,
		Update: resourceNewRelicDashboardUpdate,
		Delete: resourceNewRelicDashboardDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"icon": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "bar-chart",
			},
			"visibility": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "all",
				ValidateFunc: validation.StringInSlice([]string{"owner", "all"}, false),
			},
			"dashboard_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"editable": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "editable_by_all",
				ValidateFunc: validation.StringInSlice([]string{"read_only", "editable_by_owner", "editable_by_all", "all"}, false),
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_types": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Required: true,
							Set:      schema.HashString,
						},
						"attributes": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Set:      schema.HashString,
						},
					},
				},
			},
			"widget": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 60,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"widget_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Required: true,
						},
						"visualization": {
							Type:     schema.TypeString,
							Required: true,
						},
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"row": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"column": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"notes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"nrql": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"source": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"threshold_red": {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: Float64AtLeast(0),
						},
						"threshold_yellow": {
							Type:         schema.TypeFloat,
							Optional:     true,
							ValidateFunc: Float64AtLeast(0),
						},
						"drilldown_dashboard_id": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"duration": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"end_time": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"raw_metric_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"facet": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"order_by": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"limit": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"entity_ids": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"metric": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"units": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"scope": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"values": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"compare_with": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"offset_duration": {
										Type:     schema.TypeString,
										Required: true,
									},
									"presentation": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"color": {
													Type:     schema.TypeString,
													Required: true,
												},
												"name": {
													Type:     schema.TypeString,
													Required: true,
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
		},
	}
}

func resourceNewRelicDashboardCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	dashboard, err := expandDashboard(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating New Relic dashboard: %s", dashboard.Title)

	dashboard, err = client.CreateDashboard(*dashboard)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(dashboard.ID))

	return resourceNewRelicDashboardRead(d, meta)
}

func resourceNewRelicDashboardRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic dashboard %s", d.Id())

	dashboardID, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard, err := client.GetDashboard(dashboardID)
	if err != nil {
		if err == newrelic.ErrNotFound {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenDashboard(dashboard, d)
}

func resourceNewRelicDashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client
	dashboard, err := expandDashboard(d)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	dashboard.ID = id
	log.Printf("[INFO] Updating New Relic dashboard %d", id)

	_, err = client.UpdateDashboard(*dashboard)
	if err != nil {
		return err
	}

	return resourceNewRelicDashboardRead(d, meta)
}

func resourceNewRelicDashboardDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting New Relic dashboard %v", id)

	if err := client.DeleteDashboard(id); err != nil {
		if err == newrelic.ErrNotFound {
			return nil
		}
		return err
	}

	return nil
}
