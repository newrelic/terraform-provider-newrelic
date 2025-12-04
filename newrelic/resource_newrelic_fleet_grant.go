package newrelic

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/newrelic/newrelic-client-go/v2/pkg/authorizationmanagement"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

func resourceNewRelicFleetGrant() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicFleetGrantCreate,
		ReadContext:   resourceNewRelicFleetGrantRead,
		UpdateContext: resourceNewRelicFleetGrantUpdate,
		DeleteContext: resourceNewRelicFleetGrantDelete,

		Schema: map[string]*schema.Schema{
			"fleet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The fleet ID.",
			},
			"organization_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The organization ID.",
			},
			"grant": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "A grant for the fleet.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the grant.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The group ID.",
						},
						"role_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The role ID.",
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceNewRelicFleetGrantCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)
	grants := d.Get("grant").(*schema.Set).List()

	var successfulGrantIDs []string
	var successfulGrants []interface{}
	var diags diag.Diagnostics

	for _, grantData := range grants {
		grantMap := grantData.(map[string]interface{})
		groupID := grantMap["group_id"].(string)
		roleID := grantMap["role_id"].(int)

		input := collateAuthorizationManagementGrantAccessRequest(roleID, fleetID, groupID)

		result, err := client.AuthorizationManagement.AuthorizationManagementGrantAccess(input)
		if err != nil {
			// Log the error but continue with other grants
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Failed to grant access for group %s and role %d", groupID, roleID),
				Detail:   err.Error(),
			})
			continue
		}

		if len(result.AccessGrants) > 0 {
			successfulGrantIDs = append(successfulGrantIDs, result.AccessGrants[0].ID)
			successfulGrants = append(successfulGrants, map[string]interface{}{
				"group_id": groupID,
				"role_id":  roleID,
				"id":       result.AccessGrants[0].ID,
			})
		}
	}

	if len(successfulGrantIDs) == 0 {
		return diag.Errorf("no grants were successfully created")
	}

	// Encode the composite ID using XOR cipher + Base64
	compositeID := encodeCompositeID(successfulGrantIDs)
	d.SetId(compositeID)

	// Set only the successful grants in state
	if err := d.Set("grant", successfulGrants); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	// Fetch and store organization ID for future use
	organization, err := client.Organization.GetOrganization()
	if err != nil {
		// Log warning but don't fail the operation
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Failed to fetch organization ID",
			Detail:   fmt.Sprintf("Grants were successfully created, but failed to fetch organization ID: %v", err),
		})
	} else {
		if err := d.Set("organization_id", organization.ID); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Failed to set organization ID in state",
				Detail:   err.Error(),
			})
		}
	}

	return diags
}

func resourceNewRelicFleetGrantRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	var err error
	// grantIDs, err := decodeCompositeID(d.Id())
	//if err != nil {
	//	return diag.FromErr(fmt.Errorf("failed to decode resource ID: %w", err))
	//}

	// Try to get organization ID from state first
	organizationID := d.Get("organization_id").(string)
	if organizationID == "" {
		// If not in state, fetch it from API
		organization, err := client.Organization.GetOrganization()
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to fetch organization information: %v", err))
		}
		organizationID = organization.ID
		// Update state with organization ID for future reads
		_ = d.Set("organization_id", organizationID)
	}

	var foundGrants []interface{}

	filter := customeradministration.MultiTenantAuthorizationGrantFilterInputExpression{
		OrganizationId: &customeradministration.MultiTenantAuthorizationGrantOrganizationIdInputFilter{
			Eq: organizationID,
		},
		ScopeId: &customeradministration.MultiTenantAuthorizationGrantScopeIdInputFilter{
			Eq: d.Get("fleet_id").(string),
		},
	}

	grants, err := client.CustomerAdministration.GetGrants(
		"",
		filter,
		[]customeradministration.MultiTenantAuthorizationGrantSortInput{},
	)

	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to fetch grants: %v", err))
	}

	for _, grant := range grants.Items {
		foundGrant := map[string]interface{}{
			"id":       strconv.Itoa(grant.ID),
			"group_id": grant.Grantee.ID,
			"role_id":  grant.Role.ID,
		}

		foundGrants = append(foundGrants, foundGrant)
	}

	if len(foundGrants) == 0 {
		d.SetId("")
		return nil
	}

	if err := d.Set("grant", foundGrants); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceNewRelicFleetGrantUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)

	// Check if grants have changed
	// This case will not be encountered though; as all other attributes != grant are ForceNew
	if !d.HasChange("grant") {
		return nil
	}

	oldGrantsInterface, newGrantsInterface := d.GetChange("grant")
	oldGrants := oldGrantsInterface.(*schema.Set).List()
	newGrants := newGrantsInterface.(*schema.Set).List()

	// Create maps for efficient lookup
	oldGrantsMap := make(map[string]map[string]interface{})
	newGrantsMap := make(map[string]map[string]interface{})

	// Build map of old grants: key is "groupID:roleID"
	for _, grantData := range oldGrants {
		grantMap := grantData.(map[string]interface{})
		groupID := grantMap["group_id"].(string)
		roleID := grantMap["role_id"].(int)
		key := fmt.Sprintf("%s:%d", groupID, roleID)
		oldGrantsMap[key] = grantMap
	}

	// Build map of new grants: key is "groupID:roleID"
	for _, grantData := range newGrants {
		grantMap := grantData.(map[string]interface{})
		groupID := grantMap["group_id"].(string)
		roleID := grantMap["role_id"].(int)
		key := fmt.Sprintf("%s:%d", groupID, roleID)
		newGrantsMap[key] = grantMap
	}

	// Find grants to add (in new but not in old)
	var grantsToAdd []map[string]interface{}
	for key, grant := range newGrantsMap {
		if _, exists := oldGrantsMap[key]; !exists {
			grantsToAdd = append(grantsToAdd, grant)
		}
	}

	// Find grants to remove (in old but not in new)
	var grantsToRemove []map[string]interface{}
	for key, grant := range oldGrantsMap {
		if _, exists := newGrantsMap[key]; !exists {
			grantsToRemove = append(grantsToRemove, grant)
		}
	}

	var diags diag.Diagnostics
	var successfulGrantIDs []string
	var successfulGrants []interface{}

	// First, collect all existing grants that are still valid (not being removed)
	for key, grant := range newGrantsMap {
		if _, exists := oldGrantsMap[key]; exists {
			// Grant exists in both old and new, keep it
			successfulGrants = append(successfulGrants, grant)
			if grantID, ok := grant["id"].(string); ok && grantID != "" {
				successfulGrantIDs = append(successfulGrantIDs, grantID)
			}
		}
	}

	// Add new grants
	if len(grantsToAdd) > 0 {
		for _, grantMap := range grantsToAdd {
			groupID := grantMap["group_id"].(string)
			roleID := grantMap["role_id"].(int)

			input := collateAuthorizationManagementGrantAccessRequest(roleID, fleetID, groupID)

			result, err := client.AuthorizationManagement.AuthorizationManagementGrantAccess(input)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Failed to grant access for group %s and role %d", groupID, roleID),
					Detail:   err.Error(),
				})
				continue
			}

			if len(result.AccessGrants) > 0 {
				successfulGrantIDs = append(successfulGrantIDs, result.AccessGrants[0].ID)
				successfulGrants = append(successfulGrants, map[string]interface{}{
					"group_id": groupID,
					"role_id":  roleID,
					"id":       result.AccessGrants[0].ID,
				})
			}
		}
	}

	// Remove old grants
	if len(grantsToRemove) > 0 {
		for _, grantMap := range grantsToRemove {
			groupID := grantMap["group_id"].(string)
			roleID := grantMap["role_id"].(int)

			input := collateAuthorizationManagementRevokeAccessRequest(roleID, fleetID, groupID)

			_, err := client.AuthorizationManagement.AuthorizationManagementRevokeAccess(input)
			if err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  fmt.Sprintf("Failed to revoke access for group %s and role %d", groupID, roleID),
					Detail:   err.Error(),
				})
				continue
			}
		}
	}

	// Update the state with the new grants
	if len(successfulGrantIDs) == 0 {
		return diag.Errorf("no grants remain after update")
	}

	// Encode the composite ID with all current grant IDs
	compositeID := encodeCompositeID(successfulGrantIDs)
	d.SetId(compositeID)

	// Set the grants in state
	if err := d.Set("grant", successfulGrants); err != nil {
		diags = append(diags, diag.FromErr(err)...)
	}

	return diags
}

func resourceNewRelicFleetGrantDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	fleetID := d.Get("fleet_id").(string)
	grants := d.Get("grant").(*schema.Set).List()

	var diags diag.Diagnostics
	successCount := 0

	for _, grantData := range grants {
		grantMap := grantData.(map[string]interface{})
		groupID := grantMap["group_id"].(string)
		roleID := grantMap["role_id"].(int)

		input := collateAuthorizationManagementRevokeAccessRequest(roleID, fleetID, groupID)

		_, err := client.AuthorizationManagement.AuthorizationManagementRevokeAccess(input)
		if err != nil {
			// Log the error but continue with other grants
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("Failed to revoke access for group %s and role %d", groupID, roleID),
				Detail:   err.Error(),
			})
			continue
		}

		successCount++
	}

	// If at least one grant was successfully revoked, clear the ID
	// If no grants were successfully revoked, return an error
	if successCount == 0 && len(grants) > 0 {
		return diag.Errorf("failed to revoke any grants")
	}

	d.SetId("")
	return diags
}

// encodeCompositeID encodes a list of grant IDs into an obfuscated composite ID
// The encoding uses XOR cipher followed by Base64 encoding to make the ID non-readable
func encodeCompositeID(grantIDs []string) string {
	// Join grant IDs with a delimiter
	joinedIDs := strings.Join(grantIDs, ":")

	// Create a key from a hash for XOR operation
	hash := sha256.Sum256([]byte("newrelic-fleet-grant-v1"))
	key := hash[:16] // Use first 16 bytes as key

	// XOR encode the joined IDs
	encoded := make([]byte, len(joinedIDs))
	for i := 0; i < len(joinedIDs); i++ {
		encoded[i] = joinedIDs[i] ^ key[i%len(key)]
	}

	// Convert to hex then Base64 for additional obfuscation
	hexEncoded := hex.EncodeToString(encoded)
	return base64.StdEncoding.EncodeToString([]byte(hexEncoded))
}

// decodeCompositeID decodes an obfuscated composite ID back to a list of grant IDs
// This reverses the encoding process used in encodeCompositeID
func decodeCompositeID(compositeID string) ([]string, error) {
	// Base64 decode
	hexEncoded, err := base64.StdEncoding.DecodeString(compositeID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	// Hex decode
	encoded, err := hex.DecodeString(string(hexEncoded))
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}

	// Create the same key used for encoding
	hash := sha256.Sum256([]byte("newrelic-fleet-grant-v1"))
	key := hash[:16]

	// XOR decode
	decoded := make([]byte, len(encoded))
	for i := 0; i < len(encoded); i++ {
		decoded[i] = encoded[i] ^ key[i%len(key)]
	}

	// Split back into individual grant IDs
	joinedIDs := string(decoded)
	grantIDs := strings.Split(joinedIDs, ":")

	return grantIDs, nil
}

func collateEntityManagementAccessGrants(roleID int, fleetID string) []authorizationmanagement.AuthorizationManagementEntityAccessGrants {
	return []authorizationmanagement.AuthorizationManagementEntityAccessGrants{
		{
			Entity: authorizationmanagement.AuthorizationManagementEntity{
				Type: "fleet",
				ID:   fleetID,
			},
			RoleId: strconv.Itoa(roleID),
		},
	}
}

func collateAuthorizationManagementGrantAccessRequest(roleID int, fleetID string, groupID string) authorizationmanagement.AuthorizationManagementGrantAccess {
	return authorizationmanagement.AuthorizationManagementGrantAccess{
		EntityAccessGrants: collateEntityManagementAccessGrants(roleID, fleetID),
		GroupId:            groupID,
	}
}

func collateAuthorizationManagementRevokeAccessRequest(roleID int, fleetID string, groupID string) authorizationmanagement.AuthorizationManagementRevokeAccess {
	return authorizationmanagement.AuthorizationManagementRevokeAccess{
		EntityAccessGrants: collateEntityManagementAccessGrants(roleID, fleetID),
		GroupId:            groupID,
	}
}
