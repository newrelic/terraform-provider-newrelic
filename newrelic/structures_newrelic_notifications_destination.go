package newrelic

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/notifications"
)

func expandNotificationDestinationInput(d *schema.ResourceData) (*notifications.DestinationInput, error) {
	destination := notifications.DestinationInput{
		Name: d.Get("name").(string),
		Type: notifications.DestinationType(d.Get("type").(string)),
	}

	var auth, authOk = d.GetOk("auth")

	if !authOk {
		return nil, errors.New("notification destination requires an auth attribute")
	}

	if authOk {
		a, err := expandNotificationDestinationAuth(auth)
		if err != nil {
			return nil, err
		}

		destination.Auth = *a
	}

	err := validateDestinationAuth(destination.Auth)
	if err != nil {
		return nil, err
	}

	properties, propertiesOk := d.GetOk("properties")

	if !propertiesOk {
		return nil, errors.New("notification destination requires a properties attribute")
	}

	if propertiesOk {
		var destinationProperty map[string]interface{}

		// Need to FIX!!!
		x := properties.([]interface{})
		if len(x) > 0 {
			if x[0] != nil {
				destinationProperty = x[0].(map[string]interface{})
			}
		}

		if val, err := expandNotificationDestinationProperty(destinationProperty); err == nil {
			destination.Properties = append(destination.Properties, *val)
		}
	}

	return &destination, nil
}

func expandNotificationDestinationAuth(authList interface{}) (*notifications.AiNotificationsCredentialsInput, error) {
	auth := notifications.AiNotificationsCredentialsInput{}
	var authConfig map[string]interface{}
	list := authList.(*schema.Set).List()
	authArr := make([]int, len(list))

	if len(authArr) > 0 {
		authConfig = list[0].(map[string]interface{})

		if typeAuth, ok := authConfig["type"]; ok {
			auth.Type = notifications.AuthType(typeAuth.(string))
		}

		if auth.Type == notifications.AuthTypes.Token {
			if token, ok := authConfig["token"]; ok {
				a, err := expandNotificationDestinationTokenAuth(token)
				if err != nil {
					return nil, err
				}

				auth.Token = *a
			}
		}

		if auth.Type == notifications.AuthTypes.Basic {
			if basic, ok := authConfig["basic"]; ok {
				a, err := expandNotificationDestinationBasicAuth(basic)
				if err != nil {
					return nil, err
				}

				auth.Basic = *a
			}
		}

	}

	return &auth, nil
}

func expandNotificationDestinationBasicAuth(basicAuthList interface{}) (*notifications.BasicAuth, error) {
	basicAuth := notifications.BasicAuth{}
	var basicAuthConfig map[string]interface{}
	list := basicAuthList.(*schema.Set).List()
	basicAuthArr := make([]int, len(list))

	if len(basicAuthArr) > 0 {
		basicAuthConfig = list[0].(map[string]interface{})

		if user, ok := basicAuthConfig["user"]; ok {
			basicAuth.User = user.(string)
		}

		if password, ok := basicAuthConfig["password"]; ok {
			basicAuth.Password = notifications.SecureValue(password.(string))
		}
	}

	return &basicAuth, nil
}

func expandNotificationDestinationTokenAuth(tokenAuthList interface{}) (*notifications.TokenAuth, error) {
	tokenAuth := notifications.TokenAuth{}
	var tokenAuthConfig map[string]interface{}
	list := tokenAuthList.(*schema.Set).List()
	tokenAuthArr := make([]int, len(list))

	if len(tokenAuthArr) > 0 {
		tokenAuthConfig = list[0].(map[string]interface{})

		if prefix, ok := tokenAuthConfig["prefix"]; ok {
			tokenAuth.Prefix = prefix.(string)
		}

		if token, ok := tokenAuthConfig["token"]; ok {
			tokenAuth.Token = notifications.SecureValue(token.(string))
		}
	}

	return &tokenAuth, nil
}

func expandNotificationDestinationProperty(cfg map[string]interface{}) (*notifications.PropertyInput, error) {
	property := notifications.PropertyInput{}

	if propertyKey, ok := cfg["key"]; ok {
		property.Key = propertyKey.(string)
	}

	if propertyValue, ok := cfg["value"]; ok {
		property.Value = propertyValue.(string)
	}

	if propertyDisplayValue, ok := cfg["display_value"]; ok {
		property.DisplayValue = propertyDisplayValue.(string)
	}

	if propertyLabel, ok := cfg["label"]; ok {
		property.Label = propertyLabel.(string)
	}

	return &property, nil
}

func flattenNotificationDestination(destination *notifications.Destination, d *schema.ResourceData) error {
	if destination == nil {
		return nil
	}

	var err error

	if err = d.Set("name", destination.Name); err != nil {
		return err
	}

	if err = d.Set("type", destination.Type); err != nil {
		return err
	}

	if err := d.Set("auth", flattenNotificationDestinationAuth(&destination.Auth)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting notification auth: %#v", err)
	}

	properties, propertiesErr := flattenNotificationDestinationProperties(&destination.Properties)
	if propertiesErr != nil {
		return propertiesErr
	}

	if err := d.Set("properties", properties); err != nil {
		return err
	}

	return nil
}

func flattenNotificationDestinationProperties(p *[]notifications.Property) ([]map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	var properties []map[string]interface{}

	for _, property := range *p {
		if val, err := flattenNotificationDestinationProperty(&property); err == nil {
			properties = append(properties, val)
		}
	}

	return properties, nil
}

func flattenNotificationDestinationProperty(p *notifications.Property) (map[string]interface{}, error) {
	if p == nil {
		return nil, nil
	}

	propertyResult := make(map[string]interface{})

	propertyResult["key"] = p.Key
	propertyResult["value"] = p.Value

	if p.DisplayValue != "" {
		propertyResult["display_value"] = p.DisplayValue
	}

	if p.Label != "" {
		propertyResult["label"] = p.Label
	}

	return propertyResult, nil
}

func flattenNotificationDestinationAuth(a *notifications.Auth) interface{} {
	authConfig := map[string]interface{}{
		"authType": *a.AuthType,
	}

	if *a.AuthType == notifications.AuthTypes.Basic {
		authConfig["user"] = *a.User
	}

	if *a.AuthType == notifications.AuthTypes.Token {
		authConfig["prefix"] = *a.Prefix
	}

	return authConfig
}

func validateDestinationAuth(auth notifications.AiNotificationsCredentialsInput) error {
	if auth.Type == notifications.AuthTypes.Token && (auth.Token.Token == "" || auth.Token.Prefix == "") {
		return errors.New("token and prefix is required when using token auth type")
	}

	if auth.Type == notifications.AuthTypes.Basic && (auth.Basic.User == "" || auth.Basic.Password == "") {
		return errors.New("user and password is required when using basic auth type")
	}

	return nil
}
