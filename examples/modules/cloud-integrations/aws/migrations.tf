moved {
  from = newrelic_cloud_aws_link_account.newrelic_cloud_integration_pull 
  to = newrelic_cloud_aws_link_account.newrelic_cloud_integration_pull[0]
}

moved {
  from = newrelic_cloud_aws_integrations.newrelic_cloud_integration_pull 
  to = newrelic_cloud_aws_integrations.newrelic_cloud_integration_pull[0]
}

moved {
  from = newrelic_cloud_aws_link_account.newrelic_cloud_integration_push
  to = newrelic_cloud_aws_link_account.newrelic_cloud_integration_push[0]
}