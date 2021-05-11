---
layout: "newrelic"
page_title: "New Relic Terraform Provider Migration guide for newrelic_dashboard resource"
sidebar_current: "docs-newrelic-provider-dashboard-migration-guide"
description: |-
  Use this guide to migrate from the deprecated 'newrelic_dashboard' resource onto the new 'newrelic_one_dashboard' resource.
---

## New Relic Terraform Provider Migration guide for newrelic_dashboard resource

This guide describes the process of migrating your `newrelic_dashboard` to the new `newrelic_one_dashboard` format. The `newrelic_dashboard` has been deprecated and will stop working in the near future.

For more information check out [this Github issue](https://github.com/newrelic/terraform-provider-newrelic/issues/1297).

### Requirements

* Latest version of the New Relic CLI: https://github.com/newrelic/newrelic-cli
* Python 3

### Process

To help you migrate your `newrelic_dashboard` to the new `newrelic_one_dashboard` resource we provide two options:

- [Manual process](#manual-process): This guide will take you step by step through the process to migrate your dashboards
- [Automated process](#automated-process): This guide and Python 3 code will automate a part of the process to migrate your dashboards

Which process to choose depends on the number of Terraform `newrelic_dashboard` you need to move and the complexicity of your environment. In general we advice to run the manual process at least once so you gain confidence with the process, before switching to the automated task.

#### Manual process

Go through the following steps for each `newrelic_dashboard` resource you wish to migrate.

##### 1. Get the new dashboard JSON

As a first step we need to gather all the JSON definitions of the dashboards you want to migrate. You can do this either though the UI, or through the New Relic CLI.

To get the definition through to UI go to the dashboard you want to migrate, and click on the `Copy JSON to clipboard`. You can find the button on the top right in the same place as the time picker and the `Share` button. Once you have the JSON save it to a file.

Alternatively you can do it programmaticly with the New Relic CLI command below. Make sure you replace `[[DASHBOARD_GUID]]` with your own. You can find the GUID in the New Relic interface when look at the dashboard tags, or through the [CLI](https://github.com/newrelic/newrelic-cli/blob/main/docs/cli/newrelic_entity_search.md).

`newrelic nerdgraph query 'query($guid: EntityGuid!) { actor { entity(guid: $guid) { ... on DashboardEntity { name description pages { name description  widgets { id visualization { id } layout { column row height width } title rawConfiguration }}}}}}' --variables '{ "guid": "[[DASHBOARD_GUID]]" }' | newrelic utils jq '.actor.entity'`

##### 2. Convert to new Terraform resource

The next step is converting the JSON we captured in the previous step to the new `newrelic_one_dashboard` resource format. Use the following command to achieve this. Make sure to replace `[[NAME_OF_RESOURCE]]` with the name of the `newrelic_dashboard` resource you're replacing.

`cat your-dashboard.json | newrelic utils terraform dashboard --label [[NAME_OF_RESOURCE]] > newdashboard.tf`

At this point you need to migrate any variables, loops or other dynamic elements from your old `newrelic_dashboard` resource to the new `newrelic_one_dashboard`. So make sure to double check the generated resource definition.

##### 3. Import the dashboard

Once the `newrelic_one_dashboard` is in a good state we need to import the new state, and delete the old one. The first command will import your dashboard into the Terraform state, and it tries to link it with your newly designed resource. It will fail if it can't find the resource. Make sure you replace the `[[NAME_OF_RESOURCE]]` and `[[DASHBOARD_GUID]]` with right one.

`terraform import newrelic_one_dashboard.[[NAME_OF_RESOURCE]] [[DASHBOARD_GUID]]`

After the import has completed succesfully, you can delete the old `newrelic_dashboard` from the Terraform state. Make sure you replace `[[NAME_OF_RESOURCE]]` with the right resource name.

`terraform state rm 'newrelic_dashboard.[[NAME_OF_RESOURCE]]`

#### Automated process

To automate the process New Relic provides a Python script to read all dashboards from your existing `terraform.tfstate` and automatically generates the HCL, removes the old dashboard from state and imports the new. You still need to manually implement any dynamic elements like variables and loops.

Make sure you take a backup of your state and terraform resources before running running the code below. To run the code you need Python 3 and a CLI. Simply copy the code into a file, for example `migrate.py` and run it with python `python3 migrate.py` inside your Terraform diractory.

```python
#!/usr/bin/env python3
import json
import base64
import os
import sys

def main():
    # Open terraform state and read as json
    with open('terraform.tfstate', 'r') as file:
        data = json.load(file)

    # Iterating through the state and find newrelic_dashboard resources
    for resource in data['resources']:
        if resource['type'] == 'newrelic_dashboard':
            for instance in resource['instances']:
                process(resource['name'], instance)

def runCommand(command):
    print("Running: %s" % command)
    result = os.system(command)
    if result > 0:
        print("Terraform command failed, check output")
        sys.exit(result)

def process(resourceName, dashboard):
    # Get all data for the UUID
    dashboardUrl = dashboard['attributes']['dashboard_url']
    dashboardUrlData = dashboardUrl.split('/')
    accountId = dashboardUrlData[4]
    dashboardId = dashboardUrlData[6]

    # Generate the UUID
    guid = str(
        base64.b64encode(
            "{0}|VIZ|DASHBOARD|{1}"
            .format(accountId, dashboardId)
            .encode("utf-8")
        )
        , "utf-8"
    )[:-1]

    # Generate the HCL
    runCommand('newrelic nerdgraph query \'query($guid: EntityGuid!) { actor { entity(guid: $guid) { ... on DashboardEntity { name description permissions pages { name description widgets { id visualization { id } layout { column row height width } title rawConfiguration }}}}}}\' --variables \'{ "guid": "%s" }\' | newrelic utils jq \'.actor.entity\' | newrelic utils terraform dashboard --label %s > %s.tf' % (guid, resourceName, resourceName))

    # Import the new dashboard
    runCommand("terraform import newrelic_one_dashboard.{0} {1}".format(resourceName, guid))

    # Remove the old from state
    runCommand("terraform state rm 'newrelic_dashboard.{0}'".format(resourceName))

# Lets go
main()
```

### Run Terraform

Now that you've migrated all your dashboards, the only step that remains is running your Terraform. Some changes are expected when comparing the new resources to the old. If you encounter any problems or errors, don't hesitate to reach out on [Github](https://github.com/newrelic/terraform-provider-newrelic).
