
resource "newrelic_alert_channel" "vo-my-routingkey" {
  name = "my-routingkey"
  type = "victorops"
  config {
    key       = "xxxxxx"
    route_key = "my-routingkey"
  }
}

resource "newrelic_alert_channel" "slack-my-channel" {
  name = "my-channe"
  type = "slack"
  config {
    url = "xxxxxx"
    channel = "my-channel"
  }
}