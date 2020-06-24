package newrelic

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
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

func expandMultiLocationSyntheticsCondition(d *schema.ResourceData) (*alerts.MultiLocationSyntheticsCondition, error) {
	condition := alerts.MultiLocationSyntheticsCondition{
		Name:                      d.Get("name").(string),
		Enabled:                   d.Get("enabled").(bool),
		ViolationTimeLimitSeconds: d.Get("violation_time_limit_seconds").(int),
	}

	if attr, ok := d.GetOk("runbook_url"); ok {
		condition.RunbookURL = attr.(string)
	}

	terms, err := expandMultiLocationSyntheticsConditionTerms(d)
	if err != nil {
		return nil, err
	}

	condition.Terms = terms

	var entities []string
	for _, x := range d.Get("entities").([]interface{}) {
		entities = append(entities, x.(string))
	}

	condition.Entities = entities

	return &condition, nil
}

func expandMultiLocationSyntheticsConditionTerms(d *schema.ResourceData) ([]alerts.MultiLocationSyntheticsConditionTerm, error) {
	var expandedTerms []alerts.MultiLocationSyntheticsConditionTerm

	if critical, ok := d.GetOk("critical"); ok {
		x := critical.([]interface{})
		// A critical attribute is a list, but is limited to a single item in the shema.
		if len(x) > 0 {
			single := x[0].(map[string]interface{})

			criticalTerm, err := expandMultiLocationSyntheticsConditionTerm(single, "critical")
			if err != nil {
				return nil, err
			}
			if criticalTerm != nil {
				expandedTerms = append(expandedTerms, *criticalTerm)
			}
		}
	}

	if warning, ok := d.GetOk("warning"); ok {
		x := warning.([]interface{})
		// A warning attribute is a list, but is limited to a single item in the shema.
		if len(x) > 0 {
			single := x[0].(map[string]interface{})

			warningTerm, err := expandMultiLocationSyntheticsConditionTerm(single, "warning")
			if err != nil {
				return nil, err
			}

			if warningTerm != nil {
				expandedTerms = append(expandedTerms, *warningTerm)
			}
		}
	}

	return expandedTerms, nil
}

func expandMultiLocationSyntheticsConditionTerm(term map[string]interface{}, priority string) (*alerts.MultiLocationSyntheticsConditionTerm, error) {
	// required
	threshold := term["threshold"].(int)

	return &alerts.MultiLocationSyntheticsConditionTerm{
		Priority:  strings.ToLower(priority),
		Threshold: threshold,
	}, nil
}

func flattenMultiLocationSyntheticsCondition(condition *alerts.MultiLocationSyntheticsCondition, d *schema.ResourceData) error {
	ids, err := parseIDs(d.Id(), 2)
	if err != nil {
		return err
	}

	policyID := ids[0]

	// d.Set("policy_id", policyID)
	d.Set("name", condition.Name)
	d.Set("runbook_url", condition.RunbookURL)
	d.Set("enabled", condition.Enabled)
	d.Set("violation_time_limit_seconds", condition.ViolationTimeLimitSeconds)
	d.Set("entities", condition.Entities)
	d.Set("policy_id", policyID)

	for _, term := range condition.Terms {
		switch term.Priority {
		case "critical":
			terms := []map[string]interface{}{
				{
					"threshold": term.Threshold,
				},
			}
			if err := d.Set("critical", terms); err != nil {
				return fmt.Errorf("[DEBUG] Error setting synthetics multi-location alert condition `critical`: %v", err)
			}
		case "warning":
			terms := []map[string]interface{}{
				{
					"threshold": term.Threshold,
				},
			}
			if err := d.Set("warning", terms); err != nil {
				return fmt.Errorf("[DEBUG] Error setting synthetics multi-location alert condition `warning`: %v", err)
			}
		}
	}

	return nil
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
