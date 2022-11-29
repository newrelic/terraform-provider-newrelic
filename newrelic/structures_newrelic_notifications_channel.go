package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
)

func expandNotificationChannel(d *schema.ResourceData) notifications.AiNotificationsChannelInput {
	channel := notifications.AiNotificationsChannelInput{
		Name:          d.Get("name").(string),
		DestinationId: d.Get("destination_id").(string),
		Type:          notifications.AiNotificationsChannelType(d.Get("type").(string)),
		Product:       notifications.AiNotificationsProduct(d.Get("product").(string)),
	}
	channel.Properties = expandNotificationChannelProperties(d.Get("property").(*schema.Set).List())

	return channel
}

func expandNotificationChannelUpdate(d *schema.ResourceData) notifications.AiNotificationsChannelUpdate {
	active := d.Get("active").(bool)
	channel := notifications.AiNotificationsChannelUpdate{
		Name:   d.Get("name").(string),
		Active: &active,
	}
	channel.Properties = expandNotificationDestinationProperties(d.Get("property").(*schema.Set).List())

	return channel
}

func expandNotificationChannelProperties(properties []interface{}) []notifications.AiNotificationsPropertyInput {
	props := []notifications.AiNotificationsPropertyInput{}

	for _, p := range properties {
		props = append(props, expandNotificationChannelProperty(p.(map[string]interface{})))
	}

	return props
}

func expandNotificationChannelProperty(cfg map[string]interface{}) notifications.AiNotificationsPropertyInput {
	property := notifications.AiNotificationsPropertyInput{}

	if p, ok := cfg["key"]; ok {
		property.Key = p.(string)
	}

	if p, ok := cfg["value"]; ok {
		property.Value = p.(string)
	}

	if p, ok := cfg["display_value"]; ok {
		property.DisplayValue = p.(string)
	}

	if p, ok := cfg["label"]; ok {
		property.Label = p.(string)
	}

	return property
}

func flattenNotificationChannel(channel *notifications.AiNotificationsChannel, d *schema.ResourceData) error {
	if channel == nil {
		return nil
	}

	var err error

	if err = d.Set("name", channel.Name); err != nil {
		return err
	}

	if err = d.Set("type", channel.Type); err != nil {
		return err
	}

	if err = d.Set("product", channel.Product); err != nil {
		return err
	}

	if err = d.Set("destination_id", channel.DestinationId); err != nil {
		return err
	}

	if err := d.Set("property", flattenNotificationChannelProperties(channel.Properties)); err != nil {
		return err
	}

	if err := d.Set("account_id", channel.AccountID); err != nil {
		return err
	}

	if err := d.Set("status", channel.Status); err != nil {
		return err
	}

	if err := d.Set("active", channel.Active); err != nil {
		return err
	}

	return nil
}

func flattenNotificationChannelProperties(p []notifications.AiNotificationsProperty) []map[string]interface{} {
	properties := []map[string]interface{}{}

	for _, property := range p {
		properties = append(properties, flattenNotificationChannelProperty(property))
	}

	return properties
}

func flattenNotificationChannelProperty(p notifications.AiNotificationsProperty) map[string]interface{} {
	propertyResult := make(map[string]interface{})

	propertyResult["key"] = p.Key
	propertyResult["value"] = p.Value
	propertyResult["display_value"] = p.DisplayValue
	propertyResult["label"] = p.Label

	return propertyResult
}
