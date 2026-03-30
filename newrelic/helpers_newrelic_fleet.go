package newrelic

import (
	"context"
	"fmt"
	"strings"

	"github.com/newrelic/newrelic-client-go/v2/pkg/fleetcontrol"
)

// Fleet Management helper functions for Terraform provider

// getOrganizationID retrieves the organization ID from the provider or fetches it from the API
func getOrganizationID(ctx context.Context, providerConfig *ProviderConfig, organizationID string) (string, error) {
	if organizationID != "" {
		return organizationID, nil
	}

	org, err := providerConfig.NewClient.Organization.GetOrganization()
	if err != nil {
		return "", fmt.Errorf("failed to get organization: %w", err)
	}
	return org.ID, nil
}

// parseFleetTags converts tag strings in format "key:value1,value2" into FleetControlTagInput structs
func parseFleetTags(tagStrings []interface{}) ([]fleetcontrol.FleetControlTagInput, error) {
	if len(tagStrings) == 0 {
		return nil, nil
	}

	tags := make([]fleetcontrol.FleetControlTagInput, 0, len(tagStrings))

	for _, tagInterface := range tagStrings {
		tagStr, ok := tagInterface.(string)
		if !ok {
			continue
		}

		// Split on first colon to separate key from values
		parts := strings.SplitN(tagStr, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid tag format '%s': expected 'key:value1,value2'", tagStr)
		}

		key := strings.TrimSpace(parts[0])
		if key == "" {
			return nil, fmt.Errorf("tag key cannot be empty in '%s'", tagStr)
		}

		valueStr := strings.TrimSpace(parts[1])
		if valueStr == "" {
			return nil, fmt.Errorf("tag values cannot be empty in '%s'", tagStr)
		}

		// Split values on comma and trim whitespace
		values := strings.Split(valueStr, ",")
		for i := range values {
			values[i] = strings.TrimSpace(values[i])
		}

		tags = append(tags, fleetcontrol.FleetControlTagInput{
			Key:    key,
			Values: values,
		})
	}

	return tags, nil
}

// mapManagedEntityType converts string to FleetControlManagedEntityType
func mapManagedEntityType(typeStr string) (fleetcontrol.FleetControlManagedEntityType, error) {
	switch strings.ToUpper(typeStr) {
	case "HOST":
		return fleetcontrol.FleetControlManagedEntityTypeTypes.HOST, nil
	case "KUBERNETESCLUSTER":
		return fleetcontrol.FleetControlManagedEntityTypeTypes.KUBERNETESCLUSTER, nil
	default:
		return fleetcontrol.FleetControlManagedEntityType(""), fmt.Errorf(
			"unrecognized managed entity type '%s'", typeStr)
	}
}

// mapOperatingSystemType converts string to FleetControlOperatingSystemType
func mapOperatingSystemType(typeStr string) (fleetcontrol.FleetControlOperatingSystemType, error) {
	switch strings.ToUpper(typeStr) {
	case "LINUX":
		return fleetcontrol.FleetControlOperatingSystemTypeTypes.LINUX, nil
	case "WINDOWS":
		return fleetcontrol.FleetControlOperatingSystemTypeTypes.WINDOWS, nil
	default:
		return fleetcontrol.FleetControlOperatingSystemType(""), fmt.Errorf(
			"unrecognized operating system type '%s'", typeStr)
	}
}
