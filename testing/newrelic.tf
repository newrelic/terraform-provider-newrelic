terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "staging" # US or EU
}

variable "global_filter" {
  type    = string
  default = "WHERE email NOT LIKE '%@newrelic.com' AND email not like '%.testinator.com' AND email not like '%@datanerd.us' AND env != staging"
}

variable "flows" {
  type = list(object({
    installConfigName = string
    children = list(object({
      installFlow = string
      stages      = list(object({
        label = string
        nrql_query = string
      }))
    }))
  }))
  default = [
    {
      installConfigName = "PHP"
      children = [
        {
          installFlow = "tarInstall"
          stages  = [
            {
              label = "Check prerequisites"
              nrql_query = <<EOT
                FROM TessenAction
                SELECT funnel(customer_user_id,
                WHERE
                  installConfigName = 'linux'
                  AND installFlow = 'infraAgent'
                  AND stageLabel = 'Install the infrastructure agent'
                  AND eventName IN ('DEVEX_InstallFramework_viewstage') AS 'Install the infrastructure agent',
                WHERE
                  installConfigName = 'linux'
                  AND installFlow = 'infraAgent'
                  AND stageLabel = 'Install the infrastructure agent'
                  AND eventName IN ('DEVEX_InstallFramework_copyTerminal') AS 'Copy any helper command to send data',
                WHERE
                  installConfigName = 'linux'
                  AND installFlow = 'infraAgent'
                  AND stageLabel = 'Install the infrastructure agent'
                  AND eventName IN ('DEVEX_InstallFramework_clickCompleteStage') AS 'Continue'
                ) WHERE email NOT LIKE '%@newrelic.com' AND email not like '%.testinator.com' AND email not like '%@datanerd.us' AND env != staging
                SINCE 7 days ago
              EOT
            },
            # {
            #   label = "Select package manager"
            # },
            # {
            #   label = "Install the PHP Agent"
            # },
            # {
            #   label = "Configure the PHP agent"
            # },
            # {
            #   label = "Restart your application"
            # },
            # {
            #   label = "Optional: connect logs and infrastructure"
            # },
            # {
            #   label = "Test the connection"
            # }
          ]
        }
      ]
    }
  ]
}

resource "newrelic_one_dashboard" "install_flow_funnels" {
  # provider = newrelic.virtuoso_engineering
  name        = "Install Framework Funnels"
  permissions = "public_read_only"

  page {
    name = "Overview"

    widget_markdown {
      title  = ""
      row    = 1
      column = 1
      width  = 3
      height = 1

      text = "# Write a summary of what will be here"
    }
  }

  dynamic "page" {
    for_each = var.flows

    content {
      name = page.value.installConfigName

      dynamic "widget_markdown" {
        for_each = page.value.children

        content {
            title  = ""
            row    = 1
            column = 1
            width  = 3
            height = 1
            text = "# Test"
        }
      }
    }
  }
}
