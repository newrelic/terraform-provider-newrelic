package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
)

func resourceNewRelicChangeTrackingEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicChangeTrackingEventCreate,
		ReadContext:   schema.NoopContext,
		Delete:        schema.RemoveFromState,
		Schema: map[string]*schema.Schema{
			"entity_search": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Specify the entity to associate with the change tracking event via query.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"query": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "Entity search query string that matches exactly one entity.",
						},
					},
				},
			},
			"category_and_type_data": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The data that defines the category and type of change event.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": {
							Type:        schema.TypeList,
							Optional:    true,
							ForceNew:    true,
							MaxItems:    1,
							Description: "The category and type of the change event.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The category of the change event.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "The type of the change event.",
									},
								},
							},
						},
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
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "A URL to the changelog or a list of changes.",
												},
												"commit": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "The commit identifier.",
												},
												"deep_link": {
													Type:        schema.TypeString,
													Optional:    true,
													ForceNew:    true,
													Description: "A link to the system that generated the deployment.",
												},
												"version": {
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "The version of the deployed software.",
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
													Type:        schema.TypeString,
													Required:    true,
													ForceNew:    true,
													Description: "The identifier of the feature flag.",
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
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the event.",
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
				Description: "A short description of the change suitable for showing in various parts of New Relic One.",
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
			// Computed
			"change_tracking_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique change tracking identifier.",
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

var _ = fmt.Sprintf
var _ = log.Printf
var _ = context.Background
var _ = diag.FromErr

func resourceNewRelicChangeTrackingEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	log.Printf("[INFO] Creating New Relic change tracking event")

	changeTrackingEvent, dataHandlingRules, err := expandChangeTrackingEvent(d)
	if err != nil {
		return diag.FromErr(err)
	}

	result, err := client.ChangeTracking.ChangeTrackingCreateEventWithContext(ctx, changeTrackingEvent, dataHandlingRules)
	if err != nil {
		return diag.FromErr(err)
	}

	if result == nil {
		return diag.Errorf("error creating change tracking event: empty response")
	}

	var changeTrackingID string
	if result.ChangeTrackingEvent != nil {
		switch e := result.ChangeTrackingEvent.(type) {
		case *changetracking.ChangeTrackingDeploymentEvent:
			changeTrackingID = e.ChangeTrackingId
		case *changetracking.ChangeTrackingFeatureFlagEvent:
			changeTrackingID = e.ChangeTrackingId
		case *changetracking.ChangeTrackingGenericEvent:
			changeTrackingID = e.ChangeTrackingId
		case *changetracking.ChangeTrackingEvent:
			changeTrackingID = e.ChangeTrackingId
		}
	}

	if changeTrackingID == "" {
		changeTrackingID = fmt.Sprintf("change-tracking-event-%d", d.Get("timestamp").(int))
	}

	d.SetId(changeTrackingID)
	_ = d.Set("change_tracking_id", changeTrackingID)

	return nil
}
