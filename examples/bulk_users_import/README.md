### Overview
The `run_bulk_users_import.sh` script is a Bash script used to run a Go program that imports users in bulk (belonging to the specified group) into your Terraform configuration, to be saved to the state as `newrelic_user` resources (and `newrelic_group`, for the group with the ID specified), so as to facilitate controlling the imported user and group via Terraform, using resources in the New Relic Terraform Provider.

This is specifically useful in cases where a huge number of users were created in the New Relic One UI, added to a group, and would now like to be controlled via the `newrelic_user` and `newrelic_group` resources respectively, along with future users who would be added to the group via these resources in the New Relic Terraform Provider.

 The script works as follows - 
- Fetch users from the group with the ID specified,
- Get details of all of such users, in alignment with expected arguments of the `newrelic_user` resource,
- Write the attributes and values of each user (and the group specified) to strings in HCL (Terraform format),
- Write the generated Terraform configuration to `generated.tf` in the filepath specified, and
- Run a `terraform import` command in the filepath specified to also import all of these users (and the group) to the current Terraform state.


#### Arguments
--groupId <groupId>: The ID of the group to which the users will be imported. This is a required flag.
--apiKey <apiKey>: The User API key used for authentication. This is an optional flag, and can only be skipped if your environment has a `NEW_RELIC_API_KEY` that can be defaulted to.
--filePath <filePath>: The path to the file containing the user data to be imported. This is a required flag.

#### Example Usage
```sh
bash run_bulk_users_import.sh --groupId XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX --filePath ../../testing
```

```sh
bash run_bulk_users_import.sh --groupId XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX --apiKey XXXX-XXXXXXXXXXXXXXXX --filePath ../../testing
```



