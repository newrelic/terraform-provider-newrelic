package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/machinebox/graphql"
	"golang.org/x/exp/maps"
)

func main() {
	groupId, apiKey, filePath := parseFlags()

	client := graphql.NewClient(NERDGRAPH_API_ENDPOINT)

	logInfo("main", "Fetching group and users...")
	resourceGroup, resourceUserById := fetchGroupAndUsers(client, groupId, apiKey)

	logInfo("main", "Updating user details...")
	updateUserDetails(client, apiKey, resourceUserById)

	logInfo("main", "Generating Terraform files...")
	fileContent, importCommand := generateTerraformFiles(resourceGroup, resourceUserById)

	logInfo("main", "Creating Terraform file...")
	createTerraformFile(fileContent, filePath)

	logInfo("main", "Importing Terraform resources...")
	importTerraformResources(importCommand, filePath)
}

func parseFlags() (string, string, string) {
	logInfo("parseFlags", "Parsing command line arguments...")
	groupId := flag.String("groupId", "", "The group ID")
	apiKey := flag.String("apiKey", os.Getenv("NEW_RELIC_API_KEY"), "The API key")
	filePath := flag.String("filePath", "", "The file path")

	flag.Parse()

	if *groupId == "" || *filePath == "" {
		log.Fatalf("[ERROR] parseFlags: groupId and filePath are required arguments")
	}

	logInfo("parseFlags", "Finished parsing command line arguments.")
	return *groupId, *apiKey, *filePath
}

func fetchGroupAndUsers(client *graphql.Client, groupId string, apiKey string) (ResourceGroup, map[string]ResourceUser) {
	logInfo("fetchGroupAndUsers", "Starting to fetch group and users...")
	query := GET_GROUPS_AND_USERS_QUERY
	response := authenticationDomainsResponse{}
	err := RunGraphQLRequest(client, query, map[string]interface{}{"groupId": groupId}, apiKey, &response)
	if err != nil {
		log.Fatalf("[ERROR] fetchGroupAndUsers: Error fetching group and users: %v", err)
	}

	var resourceGroup ResourceGroup
	resourceUserById := make(map[string]ResourceUser)

	if len(response.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains) > 0 {
		authDomains := response.Actor.Organization.UserManagement.AuthenticationDomains.AuthenticationDomains
		for _, authDomain := range authDomains {
			if len(authDomain.Groups.Groups) > 0 {
				for _, group := range authDomain.Groups.Groups {
					resourceGroup.id = group.ID
					resourceGroup.name = group.DisplayName
					resourceGroup.authentication_domain_id = authDomain.ID

					for _, user := range group.Users.Users {
						resourceUser := ResourceUser{
							id:                       user.ID,
							name:                     user.Name,
							email_id:                 user.Email,
							authentication_domain_id: authDomain.ID,
						}
						resourceUserById[user.ID] = resourceUser
						resourceGroup.user_ids = append(resourceGroup.user_ids, user.ID)
					}
				}
			}
		}
	}

	logInfo("fetchGroupAndUsers", "Finished fetching group and users.")
	return resourceGroup, resourceUserById
}

func updateUserDetails(client *graphql.Client, apiKey string, resourceUserById map[string]ResourceUser) {
	logInfo("updateUserDetails", "Starting to update user details...")
	query := GET_USER_DETAILS_QUERY
	response := authenticationDomainsResponse{}
	err := RunGraphQLRequest(client, query, map[string]interface{}{"userId": maps.Keys(resourceUserById)}, apiKey, &response)
	if err != nil {
		log.Fatalf("[ERROR] updateUserDetails: Error updating user details: %v", err)
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

	logInfo("updateUserDetails", "Finished updating user details.")
}

func generateTerraformFiles(resourceGroup ResourceGroup, resourceUserById map[string]ResourceUser) (string, string) {
	logInfo("generateTerraformFiles", "Starting to generate Terraform files...")
	var fileContent string

	groupStr := `resource "newrelic_group" "%s" {
			name                     = "%s"
			authentication_domain_id = "%s"
			user_ids                 = %s
		}`

	var userIdList string = "["

	resourceName := resourceGroup.name
	resourceName = strings.ReplaceAll(resourceName, " ", "-")

	importCommand := "terraform import newrelic_group." + resourceName + " " + resourceGroup.id + " && "

	for _, user := range resourceUserById {
		userStr := `resource "newrelic_user" "%s" {
				name                     = "%s"
				email_id                 = "%s"
				authentication_domain_id = "%s"
				user_type = "%s"
			}`

		importCommand += "terraform import newrelic_user." + user.name + " " + user.id + " && "
		userIdList += ("newrelic_user." + user.name + ".id,")

		resourceName := user.name
		resourceName = strings.ReplaceAll(resourceName, " ", "-")

		fileContent += fmt.Sprintf(userStr, resourceName, user.name, user.email_id, user.authentication_domain_id, user.user_type)
		fileContent += "\n\n"
	}

	importCommand = importCommand[:len(importCommand)-4]
	userIdList = userIdList[:len(userIdList)-1]
	userIdList += "]"

	importCommand += " && terraform fmt "

	fileContent += fmt.Sprintf(groupStr, resourceName, resourceGroup.name, resourceGroup.authentication_domain_id, userIdList)
	fileContent += "\n\n"

	logInfo("generateTerraformFiles", "Finished generating Terraform files.")
	return fileContent, importCommand
}

func createTerraformFile(content string, path string) {
	logInfo("createTerraformFile", "Starting to create Terraform file...")
	filePath := path + "/generated.tf"

	if _, err := os.Stat(filePath); err == nil {
		log.Printf("[ERROR] createTerraformFile: File already exists at path: %s", filePath)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("[ERROR] createTerraformFile: Error while creating file: %v", err)
	}

	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		log.Fatalf("[ERROR] createTerraformFile: Error while writing file: %v", err)
	}

	logInfo("createTerraformFile", "Finished creating Terraform file.")
}

func importTerraformResources(command string, path string) {
	logInfo("importTerraformResources", "Starting to import Terraform resources...")
	cmd := exec.Command("bash", "-c", "cd "+path+" && "+command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("[ERROR] importTerraformResources: Error executing command: %v", err)
	}

	log.Printf("[INFO] importTerraformResources: Output:\n%s\n", output)
	logInfo("importTerraformResources", "Finished importing Terraform resources.")
}

func logInfo(functionName, message string) {
	log.Printf("[INFO] %s: %s", functionName, message)
}
