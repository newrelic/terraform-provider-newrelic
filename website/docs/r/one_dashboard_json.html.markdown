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
- `account_id` - (Optional) Determines the New Relic account where the dashboard will be created. Defaults to the account associated with the API key used.

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
### Dashboard With Multiple Pages

The following example demonstrates creating a dashboard with multiple pages, with each page comprising a couple of widgets.
```hcl
resource "newrelic_one_dashboard_json" "foo" {
  json = templatefile("dashboard.json.tftpl", {
    account_id = 123456
  })
}
```

`dashboard.json.tftpl`
```tftpl
{
  "name": "Multi - Page Dashboard Sample",
  "description": "An example that demonstrates creating a dashboard with multiple widgets, across a couple of pages.",
  "permissions": "PUBLIC_READ_WRITE",
  "pages": [
    {
      "name": "Memory Metrics",
      "description": "Widgets displaying metrics on memory utilization.",
      "widgets": [
        {
          "title": "Memory Utilization",
          "layout": {
            "column": 1,
            "row": 1,
            "width": 4,
            "height": 3
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
                "accountIds": [
                  "${account_id}"
                ],
                "query": "FROM Metric SELECT average(apm.service.memory.physical) as avgMem WHERE appName='sampleApp' TIMESERIES 2 days since 2 months ago"
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
    },
    {
      "name": "Transaction Metrics",
      "description": "Widgets displaying metrics on Transactions.",
      "widgets": [
        {
          "title": "Total Transaction Count Across Apps",
          "layout": {
            "column": 1,
            "row": 1,
            "width": 4,
            "height": 3
          },
          "linkedEntityGuids": null,
          "visualization": {
            "id": "viz.line"
          },
          "rawConfiguration": {
            "colors": {
              "seriesOverrides": [
                {
                  "color": "#418ba4",
                  "seriesName": "sampleApp"
                }
              ]
            },
            "facet": {
              "showOtherSeries": false
            },
            "legend": {
              "enabled": true
            },
            "nrqlQueries": [
              {
                "accountIds": [
                  "${account_id}"
                ],
                "query": "select count(*) from Transaction facet appName since 1 month ago TIMESERIES 1 day"
              }
            ],
            "platformOptions": {
              "ignoreTimeRange": false
            },
            "yAxisLeft": {
              "zero": true
            }
          }
        },
        {
          "title": "Response Headers Summary",
          "layout": {
            "column": 5,
            "row": 1,
            "width": 4,
            "height": 3
          },
          "linkedEntityGuids": null,
          "visualization": {
            "id": "viz.billboard"
          },
          "rawConfiguration": {
            "facet": {
              "showOtherSeries": false
            },
            "nrqlQueries": [
              {
                "accountIds": [
                  "${account_id}"
                ],
                "query": "SELECT count(*) from Transaction facet response.headers.contentType since 2 months ago"
              }
            ],
            "platformOptions": {
              "ignoreTimeRange": false
            }
          }
        }
      ]
    },
    {
      "name": "Log Metrics",
      "description": "Widgets displaying metrics on Logs.",
      "widgets": [
        {
          "title": "Log Tracker",
          "layout": {
            "column": 1,
            "row": 1,
            "width": 4,
            "height": 3
          },
          "linkedEntityGuids": null,
          "visualization": {
            "id": "viz.stacked-bar"
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
                "accountIds": [
                  "${account_id}"
                ],
                "query": "SELECT count(*) from Log since 48 hours ago TIMESERIES 3 hours"
              }
            ],
            "platformOptions": {
              "ignoreTimeRange": false
            }
          }
        }
      ]
    }
  ],
  "variables": []
}
```

### Configuring Multiple Dashboards with a Fixed Set of Pages

The following example demonstrates using a pre-defined set of pages across a variable number of dashboards, by specifying a list of the required pages as arguments in the dashboards to be created. This helps reuse pages and the widgets they comprise, across multiple dashboards.

```hcl
resource "newrelic_one_dashboard_json" "dashboard_one" {
  json = templatefile("dashboard.json.tftpl", {
    name        = "Multipage Dashboard One",
    description = "The first sample multipage dashboard in a set of three.",
    permissions = "PUBLIC_READ_WRITE",
    pages       = ["page_one.json", "page_two.json"]
  })
}

resource "newrelic_one_dashboard_json" "dashboard_two" {
  json = templatefile("dashboard.json.tftpl", {
    name        = "Multipage Dashboard Two",
    description = "The second sample multipage dashboard in a set of three.",
    permissions = "PUBLIC_READ_WRITE",
    pages       = ["page_two.json", "page_three.json"]
  })
}

resource "newrelic_one_dashboard_json" "dashboard_three" {
  json = templatefile("dashboard.json.tftpl", {
    name        = "Multipage Dashboard Three",
    description = "The third sample multipage dashboard in a set of three.",
    permissions = "PUBLIC_READ_WRITE",
    pages       = ["page_one.json", "page_two.json", "page_three.json"]
  })
}
```

`dashboard.json.tftpl`
```tftpl
{
  "name": "${name}",
  "description": "${description}",
  "permissions": "${permissions}",
  "pages": [
    %{ for index, page_name in pages }
    %{ if index!=0 }, %{ endif }
    ${ file("${page_name}") }
    %{ endfor }
  ]
}
```

`page_one.json`
```json
{
  "name": "Memory Metrics",
  "description": "Widgets displaying metrics on memory utilization.",
  "widgets": [
    {
      "title": "Memory Utilization",
      "layout": {
        "column": 1,
        "row": 1,
        "width": 4,
        "height": 3
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
            "accountIds": [
              account_id
            ],
            "query": "FROM Metric SELECT average(apm.service.memory.physical) as avgMem WHERE appName='sampleApp' TIMESERIES 2 days since 2 months ago"
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
```

`page_two.json`
```json
{
  "name": "Transaction Metrics",
  "description": "Widgets displaying metrics on Transactions.",
  "widgets": [
    {
      "title": "Total Transaction Count Across Apps",
      "layout": {
        "column": 1,
        "row": 1,
        "width": 4,
        "height": 3
      },
      "linkedEntityGuids": null,
      "visualization": {
        "id": "viz.line"
      },
      "rawConfiguration": {
        "colors": {
          "seriesOverrides": [
            {
              "color": "#418ba4",
              "seriesName": "sampleApp"
            }
          ]
        },
        "facet": {
          "showOtherSeries": false
        },
        "legend": {
          "enabled": true
        },
        "nrqlQueries": [
          {
            "accountIds": [
              account_id
            ],
            "query": "select count(*) from Transaction facet appName since 1 month ago TIMESERIES 1 day"
          }
        ],
        "platformOptions": {
          "ignoreTimeRange": false
        },
        "yAxisLeft": {
          "zero": true
        }
      }
    },
    {
      "title": "Response Headers Summary",
      "layout": {
        "column": 5,
        "row": 1,
        "width": 4,
        "height": 3
      },
      "linkedEntityGuids": null,
      "visualization": {
        "id": "viz.billboard"
      },
      "rawConfiguration": {
        "facet": {
          "showOtherSeries": false
        },
        "nrqlQueries": [
          {
            "accountIds": [
              account_id
            ],
            "query": "SELECT count(*) from Transaction facet response.headers.contentType since 2 months ago"
          }
        ],
        "platformOptions": {
          "ignoreTimeRange": false
        }
      }
    }
  ]
} 
```

`page_three.json`
```json
{
  "name": "Log Metrics",
  "description": "Widgets displaying metrics on Logs.",
  "widgets": [
    {
      "title": "Log Tracker",
      "layout": {
        "column": 1,
        "row": 1,
        "width": 4,
        "height": 3
      },
      "linkedEntityGuids": null,
      "visualization": {
        "id": "viz.stacked-bar"
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
            "accountIds": [
              account_id
            ],
            "query": "SELECT count(*) from Log since 48 hours ago TIMESERIES 3 hours"
          }
        ],
        "platformOptions": {
          "ignoreTimeRange": false
        }
      }
    }
  ]
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
                  account_id
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

### More Complex Examples

The following examples show more intricate use cases of creating dashboards from JSON files, using this resource.
- [This example](https://github.com/newrelic-experimental/nr-terraform-json-dashboard-examples/blob/main/dash_composed.tf) illustrates the use of a variable list of items to create a dashboard, that may be used iteratively to populate queries and other arguments of widgets, using Terraform template files.
- [This example](https://github.com/newrelic-experimental/nr-terraform-json-dashboard-examples/blob/main/dash_nrql_composed.tf) elaborates on the use of an apt Terraform configuration with additional dependencies, to instrument the use of values obtained from a GraphQL API response iteratively to configure widgets in the dashboard for each item in the response, using the Terraform `jsondecode` function.

More of such examples may be found in ths [GitHub repository](https://github.com/newrelic-experimental/nr-terraform-json-dashboard-examples/tree/main).

## Import

New Relic dashboards can be imported using their GUID, e.g.

```bash
$ terraform import newrelic_one_dashboard_json.my_dashboard <dashboard GUID>
```
