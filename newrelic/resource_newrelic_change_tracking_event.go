package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

func resourceNewRelicChangeTrackingEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicChangeTrackingEventCreate,
		ReadContext:   schema.NoopContext,
		Delete:        schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The category of the change event.",
			},
			"category_and_type_data": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The data that defines the category and type of change event.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_fields": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "Container for category-related input fields.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deployment": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Deployment-related fields.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"changelog": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"commit": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"deep_link": {
													Type:     schema.TypeString,
													Optional: true,
													ForceNew: true,
												},
												"version": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
											},
										},
									},
									"feature_flag": {
										Type:        schema.TypeList,
										Optional:    true,
										ForceNew:    true,
										MaxItems:    1,
										Description: "Feature flag-related fields.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"feature_flag_id": {
													Type:     schema.TypeString,
													Required: true,
													ForceNew: true,
												},
											},
										},
									},
								},
							},
						},
						"kind": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "The category and type of the change event.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"type": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"custom_attributes": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Custom attributes as a JSON string.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the event.",
			},
			"entity_search": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Specify the entity to associate with the change tracking event via query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "An identifier used to correlate account-wide changes across entities.",
			},
			"short_description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A short description of the change event.",
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The start time of the change tracking event as the number of milliseconds since the Unix epoch. Defaults to now.",
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name or identifier of the person responsible for the change.",
			},
			"validation_flags": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Flags for validation.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(listValidChangeTrackingEventValidationFlags(), false),
				},
			},
		},
	}
}

func listValidChangeTrackingEventValidationFlags() []string {
	return []string{
		string(changetracking.ChangeTrackingValidationFlagTypes.ALLOW_CUSTOM_CATEGORY_OR_TYPE),
		string(changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH),
		string(changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_REST_API_FAILURES),
	}
}

var (
	_ = fmt.Sprintf
	_ = log.Printf
	_ = context.Background
	_ = diag.FromErr
	_ = changetracking.ChangeTrackingValidationFlagTypes
	_ = common.EntityGUID("")
	_ = nrtime.EpochMilliseconds{}
)

func resourceNewRelicChangeTrackingEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	changeTrackingEvent, dataHandlingRules, err := expandChangeTrackingEvent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating New Relic change tracking event")

	result, err := client.ChangeTracking.ChangeTrackingCreateEventWithContext(ctx, changeTrackingEvent, dataHandlingRules)
	if err != nil {
		return diag.FromErr(err)
	}

	if result == nil {
		return diag.Errorf("error creating change tracking event: empty response")
	}

	if result.ChangeTrackingEvent != nil {
		switch e := result.ChangeTrackingEvent.(type) {
		case *changetracking.ChangeTrackingDeploymentEvent:
			d.SetId(e.ChangeTrackingId)
		case *changetracking.ChangeTrackingFeatureFlagEvent:
			d.SetId(e.ChangeTrackingId)
		case *changetracking.ChangeTrackingGenericEvent:
			d.SetId(e.ChangeTrackingId)
		case *changetracking.ChangeTrackingEvent:
			d.SetId(e.ChangeTrackingId)
		default:
			d.SetId(fmt.Sprintf("change-tracking-event-%d", d.Get("timestamp").(int)))
		}
	} else {
		d.SetId(fmt.Sprintf("change-tracking-event-%d", d.Get("timestamp").(int)))
	}

	return nil
}
