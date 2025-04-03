package newrelic

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

func resourceNewRelicGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNewRelicGroupCreate,
		ReadContext:   resourceNewRelicGroupRead,
		UpdateContext: resourceNewRelicGroupUpdate,
		DeleteContext: resourceNewRelicGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Description:  "The name of the group.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"authentication_domain_id": {
				Type:         schema.TypeString,
				Description:  "The ID of the authentication domain the group will belong to.",
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
				ForceNew:     true,
				// ForceNew has been added as the authentication_domain_id of a group cannot be updated post creation
				// This is because the `authenticationDomainId` field does not exist in the userManagementUpdateGroup mutation on NerdGraph
			},
			"user_ids": {
				Type:        schema.TypeSet,
				Description: "IDs of users to be added to the group.",
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNewRelicGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	name := d.Get("name").(string)
	if name == "" {
		return diag.FromErr(fmt.Errorf("`name` cannot be an empty string"))
	}

	authenticationDomainID := d.Get("authentication_domain_id").(string)
	if authenticationDomainID == "" {
		return diag.FromErr(fmt.Errorf("`authentication_domain_id` cannot be an empty string"))
	}

	log.Println("[INFO] sending request to create a group with the specified configuration")
	createGroupResponse, err := client.UserManagement.UserManagementCreateGroupWithContext(
		ctx,
		usermanagement.UserManagementCreateGroup{
			AuthenticationDomainId: authenticationDomainID,
			DisplayName:            name,
		})

	if err != nil {
		return diag.FromErr(err)
	}

	if createGroupResponse == nil {
		return diag.Errorf("error: failed to create group")
	}

	createdGroupID := createGroupResponse.Group.ID
	d.SetId(createdGroupID)
	log.Printf("[INFO] successfully created a group, ID: %s\n", createdGroupID)

	usersList := d.Get("user_ids")
	if usersList == nil {
		log.Println("[INFO] no users specified in the configuration to add to the group")
		return nil
	}

	ul := usersList.(*schema.Set).List()
	// the above would still only cause the list to have elements of type interface{} while we need string elements

	var usersListCleaned []string
	for _, u := range ul {
		if str, ok := u.(string); ok {
			usersListCleaned = append(usersListCleaned, str)
		}

	}

	if len(usersListCleaned) == 0 {
		log.Println("[INFO] no users specified in the configuration to add to the group")
		return nil
	}

	log.Printf("[INFO] sending request to add users %v to the created group %s\n", usersListCleaned, createdGroupID)
	_, err = addUsersToGroup(ctx, client, createdGroupID, usersListCleaned)

	if err != nil {
		// _ = d.Set("user_ids", schema.NewSet(schema.HashString, []interface{}{}))
		// _ = d.Set("user_ids", nil)
		return diag.FromErr(err)
	}

	log.Printf("[INFO] successfully added the following users to the group %s: %v\n", createdGroupID, usersListCleaned)
	return nil
}

func resourceNewRelicGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	authDomainID := ""
	groupID := d.Id()
	var groupName string
	var userListFetched []string

	// authenticationDomainIDs := []string{authDomainID}
	groupIDs := []string{groupID}

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		getUsersInGroupsResponse, err := client.UserManagement.UserManagementGetGroupsWithUsersWithContext(
			ctx,
			[]string{},
			groupIDs,
			"",
		)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getUsersInGroupsResponse == nil {
			return resource.RetryableError(fmt.Errorf("error fetching group: trying again"))
		}

		for _, a := range getUsersInGroupsResponse.AuthenticationDomains {
			for _, g := range a.Groups.Groups {
				if g.ID == groupID {
					authDomainID = a.ID
					groupName = g.DisplayName
					for _, u := range g.Users.Users {
						userListFetched = append(userListFetched, u.ID)
					}
				}
			}
		}

		return nil
	})

	err := d.Set("authentication_domain_id", authDomainID)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("name", groupName)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(userListFetched) != 0 {
		err = d.Set("user_ids", userListFetched)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return nil
}

func resourceNewRelicGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	log.Println("[INFO] updating the group with the specified configuration")
	groupID := d.Id()

	oldName, newName := d.GetChange("name")
	olN := oldName.(string)
	nlN := newName.(string)

	if nlN == "" {
		return diag.FromErr(fmt.Errorf("name cannot be an empty string"))
	}

	if olN != nlN {
		name := d.Get("name").(string)

		updateGroupResponse, err := client.UserManagement.UserManagementUpdateGroupWithContext(
			ctx,
			usermanagement.UserManagementUpdateGroup{
				DisplayName: name,
				ID:          d.Id(),
			})

		if err != nil {
			return diag.FromErr(err)
		}

		if updateGroupResponse == nil {
			return diag.Errorf("error: failed to update group")
		}

		log.Println("[INFO] updated the group successfully")
	}

	oldUsersList, newUsersList := d.GetChange("user_ids")

	if oldUsersList == nil && newUsersList == nil {
		log.Println("[INFO] no users specified in the configuration (both previously, and currently) to update the group with")
		return nil

	}
	ol := oldUsersList.(*schema.Set).List()
	nl := newUsersList.(*schema.Set).List()
	// the above would still only cause the list to have elements of type interface{} while we need string elements

	var oldUsersListCleaned []string
	var newUsersListCleaned []string
	for _, o := range ol {
		if str, ok := o.(string); ok {
			oldUsersListCleaned = append(oldUsersListCleaned, str)
		}
	}
	for _, n := range nl {
		if str, ok := n.(string); ok {
			newUsersListCleaned = append(newUsersListCleaned, str)
		}
	}

	if len(oldUsersListCleaned) == 0 && len(newUsersListCleaned) == 0 {
		log.Println("[INFO] no users specified in the configuration to create the group")
		return nil
	}
	log.Println("[INFO] find diffs between the two using d.GetChange(), add the right users, remove the right ones")
	oldUsersMap := make(map[string]bool)
	newUsersMap := make(map[string]bool)

	for _, item := range oldUsersListCleaned {
		oldUsersMap[item] = true
	}

	for _, item := range newUsersListCleaned {
		newUsersMap[item] = true
	}

	deletedUsers := usersDiffChecker(oldUsersListCleaned, newUsersMap)
	addedUsers := usersDiffChecker(newUsersListCleaned, oldUsersMap)

	if len(addedUsers) != 0 {
		log.Printf("[INFO] adding users %v to group %s\n", groupID, addedUsers)
		_, err := addUsersToGroup(ctx, client, groupID, addedUsers)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] successfully added the following users to the group %s: %v\n", groupID, addedUsers)
	}

	if len(deletedUsers) != 0 {
		log.Printf("[INFO] removing users %v from group %s\n", groupID, addedUsers)
		_, err := removeUsersFromGroup(ctx, client, groupID, deletedUsers)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] successfully removed the following users from the group %s: %v\n", groupID, deletedUsers)
	}

	return nil

}

func resourceNewRelicGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	groupID := d.Id()

	deleteGroupResponse, err := client.UserManagement.UserManagementDeleteGroupWithContext(ctx, usermanagement.UserManagementDeleteGroup{ID: groupID})
	if err != nil {
		return diag.FromErr(err)
	}
	if deleteGroupResponse == nil {
		return diag.FromErr(fmt.Errorf("error: failed to delete group, no response returned from NerdGraph"))
	}

	log.Printf("[INFO] successfully deleted the group with ID: %s\n", groupID)
	return nil
}

func addUsersToGroup(ctx context.Context, client *newrelic.NewRelic, groupID string, userIDs []string) (*usermanagement.UserManagementAddUsersToGroupsPayload, error) {
	log.Printf("[INFO] sending request to add user IDs %v to group %s\n", userIDs, groupID)
	addUsersToGroupResponse, err := client.UserManagement.UserManagementAddUsersToGroupsWithContext(
		ctx,
		usermanagement.UserManagementUsersGroupsInput{
			GroupIds: []string{groupID},
			UserIDs:  userIDs,
		},
	)

	if err != nil {
		return nil, err
	}

	if addUsersToGroupResponse == nil {
		return nil, fmt.Errorf("error: failed to add users to the created group")
	}

	return addUsersToGroupResponse, nil
}

func removeUsersFromGroup(ctx context.Context, client *newrelic.NewRelic, groupID string, userIDs []string) (*usermanagement.UserManagementRemoveUsersFromGroupsPayload, error) {
	log.Printf("[INFO] sending request to remove user IDs %v from group %s\n", userIDs, groupID)
	removeUsersFromGroupResponse, err := client.UserManagement.UserManagementRemoveUsersFromGroupsWithContext(
		ctx,
		usermanagement.UserManagementUsersGroupsInput{
			GroupIds: []string{groupID},
			UserIDs:  userIDs,
		},
	)

	if err != nil {
		return nil, err
	}

	if removeUsersFromGroupResponse == nil {
		return nil, fmt.Errorf("error: failed to remove users from the created group")
	}

	return removeUsersFromGroupResponse, nil
}

func usersDiffChecker(usersListCleaned []string, usersMap map[string]bool) []string {
	var usersDiff []string
	for _, item := range usersListCleaned {
		if _, exists := usersMap[item]; !exists {
			usersDiff = append(usersDiff, item)
		}
	}

	return usersDiff
}
