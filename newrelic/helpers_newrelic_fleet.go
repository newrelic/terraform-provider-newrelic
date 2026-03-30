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
// parseAgentSpecs parses agent specification from Terraform schema into FleetControlAgentInput structs
func parseAgentSpecs(agentSpecs []interface{}) ([]fleetcontrol.FleetControlAgentInput, error) {
	if len(agentSpecs) == 0 {
		return nil, nil
	}

	agents := make([]fleetcontrol.FleetControlAgentInput, 0, len(agentSpecs))

	for _, agentInterface := range agentSpecs {
		agentMap, ok := agentInterface.(map[string]interface{})
		if !ok {
			continue
		}

		agentType, _ := agentMap["agent_type"].(string)
		version, _ := agentMap["version"].(string)
		configVersionIDs, _ := agentMap["configuration_version_ids"].([]interface{})

		if agentType == "" || version == "" {
			return nil, fmt.Errorf("agent_type and version are required for each agent")
		}

		var configVersionList []fleetcontrol.FleetControlConfigurationVersionListInput
		if len(configVersionIDs) > 0 {
			configVersionList = make([]fleetcontrol.FleetControlConfigurationVersionListInput, 0, len(configVersionIDs))
			for _, versionIDInterface := range configVersionIDs {
				if versionID, ok := versionIDInterface.(string); ok && versionID != "" {
					configVersionList = append(configVersionList, fleetcontrol.FleetControlConfigurationVersionListInput{
						ID: versionID,
					})
				}
			}
		}

		agents = append(agents, fleetcontrol.FleetControlAgentInput{
			AgentType:                agentType,
			Version:                  version,
			ConfigurationVersionList: configVersionList,
		})
	}

	return agents, nil
}

// validateAgentVersionsForFleet validates that agent versions are compatible with the fleet type
// Wildcard "*" version is only allowed for KUBERNETESCLUSTER fleets, not HOST fleets
func validateAgentVersionsForFleet(ctx context.Context, client *ProviderConfig, fleetID string, agents []fleetcontrol.FleetControlAgentInput) error {
	// Fetch the fleet entity to check its managed entity type
	entityInterface, err := client.NewClient.FleetControl.GetEntityWithContext(ctx, fleetID)
	if err != nil {
		return fmt.Errorf("failed to fetch fleet details for validation: %w", err)
	}

	if entityInterface == nil {
		return fmt.Errorf("fleet with ID '%s' not found", fleetID)
	}

	// Type assert to fleet entity
	fleetEntity, ok := (*entityInterface).(*fleetcontrol.EntityManagementFleetEntity)
	if !ok {
		return fmt.Errorf("entity '%s' is not a fleet", fleetID)
	}

	// Check if this is a HOST fleet
	isHostFleet := string(fleetEntity.ManagedEntityType) == "HOST"

	// If it's a HOST fleet, validate that no agent uses "*" as version
	if isHostFleet {
		for _, agent := range agents {
			if agent.Version == "*" {
				return fmt.Errorf(
					"agent version '*' (wildcard) is not supported for HOST fleets. "+
						"Please specify an explicit version (e.g., '1.70.0'). "+
						"Wildcard versions are only supported for KUBERNETESCLUSTER fleets. "+
						"Fleet '%s' is of type: %s",
					fleetID, string(fleetEntity.ManagedEntityType))
			}
		}
	}

	// KUBERNETESCLUSTER fleets allow "*", so no validation needed
	return nil
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

// mapScopeType converts string to FleetControlEntityScope
//
//nolint:unused
func mapScopeType(typeStr string) (fleetcontrol.FleetControlEntityScope, error) {
	switch strings.ToUpper(typeStr) {
	case "ACCOUNT":
		return fleetcontrol.FleetControlEntityScopeTypes.ACCOUNT, nil
	case "ORGANIZATION":
		return fleetcontrol.FleetControlEntityScopeTypes.ORGANIZATION, nil
	default:
		return fleetcontrol.FleetControlEntityScope(""), fmt.Errorf(
			"unrecognized scope type '%s'", typeStr)
	}
}

// mapOperatingSystemType converts string to FleetControlOperatingSystemType
func mapOperatingSystemType(osStr string) (fleetcontrol.FleetControlOperatingSystemType, error) {
	switch strings.ToUpper(osStr) {
	case "LINUX":
		return fleetcontrol.FleetControlOperatingSystemTypeTypes.LINUX, nil
	case "WINDOWS":
		return fleetcontrol.FleetControlOperatingSystemTypeTypes.WINDOWS, nil
	default:
		return fleetcontrol.FleetControlOperatingSystemType(""), fmt.Errorf(
			"unrecognized operating system type '%s'", osStr)
	}
}
