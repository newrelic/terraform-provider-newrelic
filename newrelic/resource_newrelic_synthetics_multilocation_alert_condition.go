package newrelic

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
)

// syntheticsMultiLocationConditionTermSchema returns the schema used for a critial or warning term priority.
func syntheticsMultiLocationConditionTermSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"threshold": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The minimum number of monitor locations that must be concurrently failing before a violation is opened.",
			},
		},
	}
}

func resourceNewRelicSyntheticsMultiLocationAlertCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicSyntheticsMultiLocationAlertConditionCreate,
		Read:   resourceNewRelicSyntheticsMultiLocationAlertConditionRead,
		Update: resourceNewRelicSyntheticsMultiLocationAlertConditionUpdate,
		Delete: resourceNewRelicSyntheticsMultiLocationAlertConditionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of this condition.",
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the policy where this condition will be used.",
			},
			"runbook_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Runbook URL to display in notifications.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Set whether to enable the alert condition. Defaults to true.",
			},
			"entities": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "The GUIDs of the Synthetics monitors to alert on.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"critical": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Elem:        syntheticsMultiLocationConditionTermSchema(),
				Description: "A condition term with priority set to critical.",
			},
			"warning": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Optional:    true,
				Elem:        syntheticsMultiLocationConditionTermSchema(),
				Description: "A condition term with priority set to warning.",
			},
			"violation_time_limit_seconds": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 3600, 7200, 14400, 28800, 43200, 86400}),
				Description:  "The maximum number of seconds a violation can remain open before being closed by the system.  Must be one of: 0, 3600, 7200, 14400, 28800, 43200, 86400",
			},
		},
	}
}

func resourceNewRelicSyntheticsMultiLocationAlertConditionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	policyID := d.Get("policy_id").(int)
	condition, err := expandMultiLocationSyntheticsCondition(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating New Relic Alerts multi-location failure condition %s", condition.Name)

	condition, err = client.Alerts.CreateMultiLocationSyntheticsCondition(*condition, policyID)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{policyID, condition.ID}))

	return resourceNewRelicSyntheticsMultiLocationAlertConditionRead(d, meta)
}

func resourceNewRelicSyntheticsMultiLocationAlertConditionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic Alerts multi-location failure condition %s", d.Id())

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]
	id := ids[1]

	_, err = client.Alerts.GetPolicy(policyID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	condition, err := client.Alerts.GetMultiLocationSyntheticsCondition(policyID, id)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenMultiLocationSyntheticsCondition(condition, d)
}

func resourceNewRelicSyntheticsMultiLocationAlertConditionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	condition, err := expandMultiLocationSyntheticsCondition(d)
	if err != nil {
		return err
	}

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	condition.ID = id

	log.Printf("[INFO] Udpating New Relic Alerts multi-location failure condition %d", id)

	_, err = client.Alerts.UpdateMultiLocationSyntheticsCondition(*condition)
	if err != nil {
		return err
	}

	return resourceNewRelicSyntheticsMultiLocationAlertConditionRead(d, meta)
}

func resourceNewRelicSyntheticsMultiLocationAlertConditionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	id := ids[1]

	log.Printf("[INFO] Deleting New Relic Alerts multi-location failure condition %d", id)

	_, err = client.Alerts.DeleteMultiLocationSyntheticsCondition(id)

	if err != nil {
		return err
	}

	return nil
}
