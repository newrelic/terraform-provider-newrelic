package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
	"golang.org/x/exp/maps"
)

type authenticationDomainsResponse struct {
	Actor usermanagement.Actor `json:"actor"`
}

type ResourceUser struct {
	id                       string
	name                     string
	email_id                 string
	authentication_domain_id string
	user_type                string
}

type ResourceGroup struct {
	id                       string
	name                     string
	authentication_domain_id string
	user_ids                 []string
}

var userTier = map[string]string{
	"Basic":         "BASIC_USER_TIER",
	"Core":          "CORE_USER_TIER",
	"Full platform": "FULL_USER_TIER",
}

func main() {

	args := os.Args

	client := graphql.NewClient("https://api.newrelic.com/graphql")
	queryData(client, args[1], args[2], args[3])

}

// *UserManagementAuthenticationDomains
func queryData(client *graphql.Client, groupId string, apiKey string, path string) {

	query := `
	query(
		$groupId: [ID!]
	){
		actor {
			organization {
				userManagement {
					authenticationDomains {
					nextCursor
						authenticationDomains {
							name
							id
							groups(id: $groupId) {
								groups {
									id
									displayName
									users {
										users {
											email
											id
											name
											timeZone
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	`

	request := graphql.NewRequest(query)
	request.Var("groupId", groupId)
	request.Header.Set("Api-Key", apiKey)

	response := authenticationDomainsResponse{}
	err := client.Run(context.Background(), request, &response)
	if err != nil {
		panic(err)
	}

	var resourceGroup ResourceGroup
	var resourceUserById map[string]ResourceUser = make(map[string]ResourceUser)

	if len(response.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains) > 0 {

		authDomains := response.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains

		for _, authDomain := range authDomains {
			if len(authDomain.Groups.Groups) > 0 {
				for _, group := range authDomain.Groups.Groups {

					// Create Resource Group
					resourceGroup.id = group.ID
					resourceGroup.name = group.DisplayName
					resourceGroup.authentication_domain_id = authDomain.ID

					for _, user := range group.Users.Users {

						var resourceUser ResourceUser
						resourceUser.id = user.ID
						resourceUser.name = user.Name
						resourceUser.email_id = user.Email
						resourceUser.authentication_domain_id = authDomain.ID

						resourceUserById[user.ID] = resourceUser
						resourceGroup.user_ids = append(resourceGroup.user_ids, user.ID)

					}

				}
			}
		}

	}

	query = `query(
		$userId: [ID!]
	){
		actor {
			organization {
				userManagement {
					authenticationDomains {
						nextCursor
						authenticationDomains {
							name
							id
							users(id: $userId) {
								users {
									email
									id
									name
									type {
										displayName
									}
								}
							}
						}
					}
				}
			}
		}
	}`

	request = graphql.NewRequest(query)
	request.Var("userId", maps.Keys(resourceUserById))
	request.Header.Set("Api-Key", apiKey)

	response = authenticationDomainsResponse{}
	err = client.Run(context.Background(), request, &response)
	if err != nil {
		panic(err)
	}

	if len(response.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains) > 0 {

		authDomains := response.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains

		for _, authDomain := range authDomains {

			if len(authDomain.Users.Users) > 0 {
				for _, user := range authDomain.Users.Users {

					val, ok := resourceUserById[user.ID]
					if ok {
						val.user_type = userTier[user.Type.DisplayName]
						resourceUserById[user.ID] = val
					}

				}
			}
		}

	}

	var fileContent string

	str := `resource "newrelic_group" "%s" {
			name                     = "%s"
			authentication_domain_id = "%s"
			user_ids                 = %s
		}`

	var userIdList string = "["

	resourceName := resourceGroup.name
	resourceName = strings.ReplaceAll(resourceName, " ", "-")

	importCommand := "terraform import newrelic_group." + resourceName + " " + resourceGroup.id + " && "

	for _, user := range resourceUserById {

		str := `resource "newrelic_user" "%s" {
				name                     = "%s"
				email_id                 = "%s"
				authentication_domain_id = "%s"
				user_type = "%s"
			}`

		importCommand += "terraform import newrelic_user." + user.name + " " + user.id + " && "

		userIdList += ("newrelic_user." + user.name + ".id,")

		resourceName := user.name
		resourceName = strings.ReplaceAll(resourceName, " ", "-")

		fileContent += fmt.Sprintf(str, resourceName, user.name, user.email_id, user.authentication_domain_id, user.user_type)
		fileContent += "\n\n"

	}

	importCommand = importCommand[:len(importCommand)-4]

	userIdList = userIdList[:len(userIdList)-1]
	userIdList += "]"

	fileContent += fmt.Sprintf(str, resourceName, resourceGroup.name, resourceGroup.authentication_domain_id, userIdList)
	fileContent += "\n\n"

	createTerraformFile(fileContent, path)

	importTerraformResource(importCommand, path)

}

func createTerraformFile(content string, path string) {

	file, err := os.Create("../../" + path + "/generated.tf")
	if err != nil {
		fmt.Println("Error while creating file: ", err)
		return
	}

	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		fmt.Println("Error while writing file: ", err)
		return
	}

}

func importTerraformResource(command string, path string) {

	cmd := exec.Command("bash", "-c", "cd ../.. && cd "+path+" && "+command)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error executing command: %s\n", err)
		return
	}

	// Print the output of the command
	fmt.Printf("Output:\n%s\n", output)
}
