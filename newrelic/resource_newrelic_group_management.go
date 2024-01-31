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

func resourceNewRelicGroupManagement() *schema.Resource {
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
				Type:        schema.TypeString,
				Description: "The ID of the authentication domain the user will belong to.",
				Required:    true,
				ForceNew:    true,
			},
			"users": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "IDs of users to be added to the group",
				Default:     nil,
			},
		},
	}
}

func resourceNewRelicGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient
	log.Println("[INFO] creating a group with the specified configuration")

	name := d.Get("name").(string)
	if name == "" {
		return diag.FromErr(fmt.Errorf("`name` cannot be an empty string"))
	}

	authenticationDomainId := d.Get("authentication_domain_id").(string)

	createGroupResponse, err := client.UserManagement.UserManagementCreateGroup(
		usermanagement.UserManagementCreateGroup{
			AuthenticationDomainId: authenticationDomainId,
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
	log.Println("[INFO] created the group successfully")

	usersList := d.Get("users")
	if usersList == nil {
		log.Println("[INFO] no users specified in the configuration to create the group")
		//_ = d.Set("users", schema.NewSet(schema.HashString, []interface{}{}))
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
		log.Println("[INFO] no users specified in the configuration to create the group")
		return nil
	}

	addUsersToGroupResponse, err := client.UserManagement.UserManagementAddUsersToGroups(
		usermanagement.UserManagementUsersGroupsInput{
			GroupIds: []string{createdGroupID},
			UserIDs:  usersListCleaned,
		},
	)

	if err != nil {
		return diag.FromErr(err)
	}

	if addUsersToGroupResponse == nil {
		return diag.Errorf("error: failed to add users to the created group")
	}

	return nil
}

func resourceNewRelicGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	authDomainID := d.Get("authentication_domain_id").(string)
	groupID := d.Id()
	var groupName string
	var userListFetched []string

	authenticationDomainIDs := []string{authDomainID}
	groupIDs := []string{groupID}

	retryErr := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		getUsersInGroupsResponse, err := client.UserManagement.GetUsersInGroups(authenticationDomainIDs, groupIDs)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if getUsersInGroupsResponse == nil {
			return resource.RetryableError(fmt.Errorf("error fetching group: trying again"))
		}

		for _, a := range getUsersInGroupsResponse.AuthenticationDomains {
			if a.ID == authDomainID {
				for _, g := range a.Groups.Groups {
					if g.ID == groupID {
						groupName = g.DisplayName
						for _, u := range g.Users.Users {
							userListFetched = append(userListFetched, u.ID)
						}
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
		err = d.Set("users", userListFetched)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	//else {
	//	err = d.Set("users", nil)
	//	if err != nil {
	//		return diag.FromErr(err)
	//	}
	//}

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

		updateGroupResponse, err := client.UserManagement.UserManagementUpdateGroup(
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

		//_ := updateGroupResponse.Group.ID
		log.Println("[INFO] updated the group successfully")
	}

	// usersList := d.Get("users")

	oldUsersList, newUsersList := d.GetChange("users")

	if oldUsersList == nil && newUsersList == nil {
		log.Println("[INFO] no users specified in the configuration (both previously, and currently) to update the group with")
		//_ = d.Set("users", schema.NewSet(schema.HashString, []interface{}{}))
		return nil

	} else {
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
		} else {
			if len(oldUsersListCleaned) == 0 && len(newUsersListCleaned) != 0 {
				log.Println("[INFO] new users have been added to the group in the update process. ADDING USERS TO THE GROUP")
				_, err := addUsersToGroup(client, groupID, newUsersListCleaned)
				//addUsersToGroupResponse, err := client.UserManagement.UserManagementAddUsersToGroups(
				//	usermanagement.UserManagementUsersGroupsInput{
				//		GroupIds: []string{groupID},
				//		UserIDs:  newUsersListCleaned,
				//	},
				//)
				if err != nil {
					return diag.FromErr(err)
				}
				//if addUsersToGroupResponse == nil {
				//	return diag.Errorf("error: failed to add users to the created group")
				//}
			} else if len(oldUsersListCleaned) != 0 && len(newUsersListCleaned) == 0 {
				log.Println("[INFO] all users in the group have been deleted in the update process. REMOVING USERS FROM THE GROUP")
				_, err := removeUsersFromGroup(client, groupID, oldUsersListCleaned)
				//removeUsersFromGroupResponse, err := client.UserManagement.UserManagementRemoveUsersFromGroups(
				//	usermanagement.UserManagementUsersGroupsInput{
				//		GroupIds: []string{groupID},
				//		UserIDs:  oldUsersListCleaned,
				//	},
				//)
				if err != nil {
					return diag.FromErr(err)
				}
				//
				//if removeUsersFromGroupResponse == nil {
				//	return diag.Errorf("error: failed to remove users from the created group")
				//}
			} else {
				log.Println("[INFO] find diffs between the two using d.GetChange(), add the right users, remove the right ones")
				oldUsersMap := make(map[string]bool)
				newUsersMap := make(map[string]bool)

				for _, item := range oldUsersListCleaned {
					oldUsersMap[item] = true
				}

				for _, item := range newUsersListCleaned {
					newUsersMap[item] = true
				}

				var deletedUsers []string
				for _, item := range oldUsersListCleaned {
					if _, exists := newUsersMap[item]; !exists {
						deletedUsers = append(deletedUsers, item)
					}
				}

				var addedUsers []string
				for _, item := range newUsersListCleaned {
					if _, exists := oldUsersMap[item]; !exists {
						addedUsers = append(addedUsers, item)
					}
				}

				_, err := addUsersToGroup(client, groupID, addedUsers)
				//addUsersToGroupResponse, err := client.UserManagement.UserManagementAddUsersToGroups(
				//	usermanagement.UserManagementUsersGroupsInput{
				//		GroupIds: []string{groupID},
				//		UserIDs:  addedUsers,
				//	},
				//)

				if err != nil {
					return diag.FromErr(err)
				}

				//if addUsersToGroupResponse == nil {
				//	return diag.Errorf("error: failed to add users to the created group")
				//}

				_, err = removeUsersFromGroup(client, groupID, deletedUsers)
				//removeUsersFromGroupResponse, err := client.UserManagement.UserManagementRemoveUsersFromGroups(
				//	usermanagement.UserManagementUsersGroupsInput{
				//		GroupIds: []string{groupID},
				//		UserIDs:  deletedUsers,
				//	},
				//)

				if err != nil {
					return diag.FromErr(err)
				}

				//if removeUsersFromGroupResponse == nil {
				//	return diag.Errorf("error: failed to remove users from the created group")
				//}

				log.Println("Deleted users:", deletedUsers)
				log.Println("Added users:", addedUsers)
			}

		}
		return nil
	}
}

func resourceNewRelicGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	deleteGroupResponse, err := client.UserManagement.UserManagementDeleteGroup(usermanagement.UserManagementDeleteGroup{ID: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}
	if deleteGroupResponse == nil {
		return diag.FromErr(fmt.Errorf("error: failed to delete group, no response returned from NerdGraph"))
	}

	return nil
}

func addUsersToGroup(client *newrelic.NewRelic, groupID string, userIDs []string) (*usermanagement.UserManagementAddUsersToGroupsPayload, error) {
	log.Printf("[INFO] processing request to add user IDs %v to group %s\n", userIDs, groupID)
	addUsersToGroupResponse, err := client.UserManagement.UserManagementAddUsersToGroups(
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

func removeUsersFromGroup(client *newrelic.NewRelic, groupID string, userIDs []string) (*usermanagement.UserManagementRemoveUsersFromGroupsPayload, error) {
	log.Printf("[INFO] processing request to remove add user IDs %v from group %s\n", userIDs, groupID)
	removeUsersFromGroupResponse, err := client.UserManagement.UserManagementRemoveUsersFromGroups(
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
