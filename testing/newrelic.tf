terraform {
  required_providers {
    newrelic = {
      source = "newrelic/newrelic"
    }
  }
}

provider "newrelic" {
  region = "US" # US or EU
}

resource "newrelic_one_dashboard_json" "dashboard_one" {
  json = templatefile("dashboard.json.tftpl", {
    name        = "2180 Multipage Dashboard One",
    description = "A multipage dashboard that helps view multiple things at once.",
    permissions = "PUBLIC_READ_WRITE",
    pages       = ["page_one.json", "page_two.json"]
  })
}

resource "newrelic_one_dashboard_json" "dashboard_two" {
  json = templatefile("dashboard.json.tftpl", {
    name        = "2180 Multipage Dashboard Two",
    description = "A multipage dashboard that helps view multiple things at once.",
    permissions = "PUBLIC_READ_WRITE",
    pages       = ["page_two.json", "page_three.json"]
  })
}

resource "newrelic_one_dashboard_json" "dashboard_three" {
  json = templatefile("dashboard.json.tftpl", {
    name        = "2180 Multipage Dashboard Three",
    description = "A multipage dashboard that helps view multiple things at once.",
    permissions = "PUBLIC_READ_WRITE",
    pages       = ["page_one.json", "page_two.json", "page_three.json"]
  })
}