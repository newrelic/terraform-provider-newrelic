package newrelic

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceNewRelicCustomEvents() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicCustomEventsCreate,
		Read:   schema.Noop,
		Update: resourceNewRelicCustomEventsUpdate,
		Delete: schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"event": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem:     eventSchema(),
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
				Description: "The event's name",
				ValidateFunc: validation.StringMatch(
					regexp.MustCompile(`[a-zA-Z0-9_:]+`),
					"only alphanumeric characters, underscores, and colons allowed for event type",
				),
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Description: "Unix epoch timestamp in either seconds or milliseconds",
				Optional:    true,
			},
			"attribute": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				MaxItems: 255,
				Elem:     eventValueSchema(),
			},
		},
	}
}

func eventValueSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				Description:  "The name for the attribute",
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
			},
			"value": {
				Type:        schema.TypeString,
				Description: "The value for the attribute",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type for attribute value. Accepted values are string, int, or float. Defaults to string.",
				Optional:    true,
			},
		},
	}
}

type Event struct {
	Type       string
	Timestamp  *int
	Attributes []map[string]interface{}
}

func (e *Event) MarshalJSON() ([]byte, error) {
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

func resourceNewRelicCustomEventsCreate(d *schema.ResourceData, meta interface{}) error {
	var payload []*Event

	if v, ok := d.GetOkExists("event"); ok {
		events := v.(*schema.Set).List()
		payload = make([]*Event, len(events))

		for i, event := range v.(*schema.Set).List() {
			attrs := event.(map[string]interface{})["attribute"].(*schema.Set).List()
			eventPayload := &Event{
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
						return fmt.Errorf("unable to convert value %q to an int", val)
					}
					val = f
				case "float":
					f, err := strconv.ParseFloat(val.(string), 64)
					if err != nil {
						return fmt.Errorf("unable to convert value %q to a float", val)
					}
					val = f
				case "string": // do nothing
				case "": // do nothing
				default:
					return fmt.Errorf("%q is not a valid type for an attribute value", valType)
				}

				eventPayload.Attributes[i] = map[string]interface{}{key: val}
			}
			payload[i] = eventPayload
		}
	}

	b, err := json.Marshal(payload)
	if err != nil {
		log.Print(err)
	}
	log.Printf("%+v", string(b))

	return nil
}

func resourceNewRelicCustomEventsUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}
