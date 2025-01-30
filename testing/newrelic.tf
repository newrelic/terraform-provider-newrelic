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

resource "newrelic_application_settings_copy" "app" {
   guid = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTQ2NjMxMDAz"
   apm_config {
     enable_server_side_config = true
   }
}


