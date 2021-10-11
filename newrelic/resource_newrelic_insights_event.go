package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNewRelicInsightsEvent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicInsightsEventCreate,
		ReadContext:   schema.NoopContext,
		Delete:        schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"event": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     eventSchema(),
				ForceNew: true,
			},
		},
	}
}

func eventSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The event's name. Can be a combination of alphanumeric characters, underscores, and colons.",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`^[a-zA-Z0-9_: ]+$`),
					"only alphanumeric characters, underscores, and colons allowed for event type",
				),
				ForceNew: true,
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Description: "Must be a Unix epoch timestamp. You can define timestamps either in seconds or in milliseconds.",
				Optional:    true,
				ForceNew:    true,
			},
			"attribute": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				MaxItems:    255,
				Elem:        eventValueSchema(),
				ForceNew:    true,
				Description: "An attribute to include in your event payload. Multiple attribute blocks can be defined for an event.",
			},
		},
	}
}

func eventValueSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Description:  "The name of the attribute.",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				ForceNew:     true,
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value of the attribute.",
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Type:         schema.TypeString,
				Description:  "Specify the type for the attribute value. This is useful when passing integer or float values to Insights. Allowed values are string, int, or float. Defaults to string.",
				ValidateFunc: validation.StringInSlice([]string{"", "int", "float", "string"}, true),
				Optional:     true,
				ForceNew:     true,
			},
		},
	}
}

// InsightsEvent represents an Insights event
type InsightsEvent struct {
	Type       string
	Timestamp  *int
	Attributes []map[string]interface{}
}

// MarshalJSON implements a custom marshal method for InsightsEvent
func (e *InsightsEvent) MarshalJSON() ([]byte, error) {
	event := map[string]interface{}{
		"eventType": e.Type,
	}
	if e.Timestamp != nil {
		event["timestamp"] = *e.Timestamp
	}
	for _, attr := range e.Attributes {
		for k, v := range attr {
			event[k] = v
		}
	}

	b, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func resourceNewRelicInsightsEventCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*ProviderConfig).InsightsInsertClient
	var eventsPayload []*InsightsEvent

	if v, ok := d.GetOkExists("event"); ok {
		events := v.(*schema.Set).List()
		eventsPayload = make([]*InsightsEvent, len(events))

		for i, event := range v.(*schema.Set).List() {
			attrs := event.(map[string]interface{})["attribute"].(*schema.Set).List()
			eventPayload := &InsightsEvent{
				Type:       event.(map[string]interface{})["type"].(string),
				Attributes: make([]map[string]interface{}, len(attrs)),
			}
			if timestamp := event.(map[string]interface{})["timestamp"].(int); timestamp > 0 {
				eventPayload.Timestamp = &timestamp
			}
			for i, attr := range attrs {
				key := attr.(map[string]interface{})["key"].(string)
				val := attr.(map[string]interface{})["value"]

				switch valType := attr.(map[string]interface{})["type"]; valType {
				case "int":
					f, err := strconv.Atoi(val.(string))
					if err != nil {
						return diag.Errorf("unable to convert value %q to an int", val)
					}
					val = f
				case "float":
					f, err := strconv.ParseFloat(val.(string), 64)
					if err != nil {
						return diag.Errorf("unable to convert value %q to a float", val)
					}
					val = f
				case "string": // do nothing
				case "": // do nothing
				default:
					return diag.Errorf("%q is not a valid type for an attribute value", valType)
				}

				eventPayload.Attributes[i] = map[string]interface{}{key: val}
			}
			eventsPayload[i] = eventPayload
		}
	}

	if err := client.PostEvent(eventsPayload); err != nil {
		return diag.Errorf("error occurreed while posting events to Insights: %q", err)
	}

	d.SetId(fmt.Sprintf("%d", rand.Int()))
	return nil
}
