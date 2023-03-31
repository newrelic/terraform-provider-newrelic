---
layout: "newrelic"
page_title: "New Relic: newrelic_one_dashboard_json"
sidebar_current: "docs-newrelic-resource-one-dashboard-json"
description: |-
  Create and manage dashboards from a JSON file.
---

# Resource: newrelic_one_dashboard_json

## Example Usage: Create a New Relic One Dashboard from a JSON file

```hcl
resource "newrelic_one_dashboard_json" "foo" {
   json = file("dashboard.json")
}
```

## Argument Reference

The following arguments are supported:

- `json` - (Required) The JSON export of a dashboard. [The JSON can be exported from the UI](https://docs.newrelic.com/docs/query-your-data/explore-query-data/dashboards/dashboards-charts-import-export-data/#dashboards)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

- `guid` - The unique entity identifier of the dashboard in New Relic.
- `permalink` - The URL for viewing the dashboard.
- `updated_at` - The date and time when the dashboard was last updated.

## Additional Examples

### Template

Below is an example how you can use [templatefile](https://www.terraform.io/language/functions/templatefile) to dynamically generate pages based on a list. We also replace the `account_id` which is usually hardcoded in the json with a variable.

```hcl
resource "newrelic_one_dashboard_json" "bar" {
   json = templatefile("dashboard.json.tftpl", {
      account_id = 1234567,
      applications = [
        {
            name = "Application one",
            appName = "app1",
        },
        {
            name = "Application two",
            appName = "app2",
        },
      ],
   })
}
```

`dashboard.json.tftpl`
```json
{
  "name": "Applications",
  "description": null,
  "permissions": "PUBLIC_READ_WRITE",
  "pages": [
  %{for index, application in applications}
  %{ if index!=0 }, %{ endif }
    {
      "name": "${application.name}",
      "description": null,
      "widgets": [
        {
          "title": "",
          "layout": {
            "column": 1,
            "row": 1,
            "width": 12,
            "height": 5
          },
          "linkedEntityGuids": null,
          "visualization": {
            "id": "viz.line"
          },
          "rawConfiguration": {
            "facet": {
              "showOtherSeries": false
            },
            "legend": {
              "enabled": true
            },
            "nrqlQueries": [
              {
                "accountId": "${account_id}",
                "query": "FROM Transaction SELECT average(duration) WHERE appName = '${application.appName}' TIMESERIES "
              }
            ],
            "platformOptions": {
              "ignoreTimeRange": false
            },
            "yAxisLeft": {
              "zero": true
            }
          }
        }
      ]
    }
  %{ endfor }
  ],
  "variables": []
}
```
### Setting Thresholds

The following example demonstrates setting thresholds on a billboard widget.

`dashboard.json`
```json
{
  "name" : "Sample",
  "permissions" : "PUBLIC_READ_WRITE",
  "pages" : [
    {
      "name" : "Sample Page",
      "description" : "A guide to the metrics of daily transactions on the website.",
      "widgets" : [
        {
          "title" : "Transaction Failure Tracker",
          "layout" : {
            "column" : 1,
            "row" : 1,
            "width" : 3,
            "height" : 5
          },
          "visualization" : {
            "id" : "viz.billboard"
          },
          "rawConfiguration" : {
            "nrqlQueries" : [
              {
                "accountIds" : [
                  {Your-Account-ID}
                ],
                "query" : "SELECT count(*) from Transaction where httpResponseCode!=200 since 1 hour ago"
              }
            ],
            "thresholds" : [
              {
                "alertSeverity" : "WARNING",
                "value" : 15
              },
              {
                "alertSeverity" : "CRITICAL",
                "value" : 40
              }
            ]
          }
        }
      ]
    }
  ]
}
```


## Import

New Relic dashboards can be imported using their GUID, e.g.

```bash
$ terraform import newrelic_one_dashboard_json.my_dashboard <dashboard GUID>
```
