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

resource "newrelic_application_settings" "app" {
   guid = "MzgwNjUyNnxBUE18QVBQTElDQVRJT058NTM1Mzg4OTcz"
   name = "testing name alias 2"
   app_apdex_threshold = "0.5"
   enable_real_user_monitoring = true
   transaction_tracer{
     explain_query_plans{
       query_plan_threshold_value = "0.5"
       query_plan_threshold_type = "VALUE"
     }

   }
#   end_user_apdex_threshold = "0.8"
    error_collector{
      expected_error_classes = ["errr1"]
    }
    tracer_type ="NONE"
   enable_thread_profiler = false
}


