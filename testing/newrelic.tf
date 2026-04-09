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

# Example: Linux Host Fleet
resource "newrelic_fleet" "linux_production" {
  name                = "Production Linux Hosts"
  managed_entity_type = "HOST"
  operating_system    = "LINUX"
  description         = "Fleet for managing production Linux hosts"

  tags = [
    "environment:production",
    "team:platform,ops"
  ]
}

# Example: Kubernetes Cluster Fleet
resource "newrelic_fleet" "k8s_clusters" {
  name                = "Production Kubernetes Clusters"
  managed_entity_type = "KUBERNETESCLUSTER"
  description         = "Fleet for managing production K8s clusters"

  tags = [
    "environment:production",
    "platform:kubernetes"
  ]
}

# Example: Infrastructure Agent Configuration
resource "newrelic_fleet_configuration" "infra_config" {
  name                = "Production Infrastructure Config"
  agent_type          = "INFRASTRUCTURE"
  managed_entity_type = "HOST"

  configuration_content = <<-EOT
    log:
      level: info
      file: /var/log/newrelic-infra/newrelic-infra.log
    metrics:
      enabled: true
      system_sample_rate: 15
    integrations:
      - name: nri-docker
        enabled: true
  EOT
}

# Example: Kubernetes Agent Configuration
resource "newrelic_fleet_configuration" "k8s_config" {
  name                = "Production K8s Config"
  agent_type          = "KUBERNETES"
  managed_entity_type = "KUBERNETESCLUSTER"

  configuration_content = <<-EOT
    cluster:
      enabled: true
      name: production-cluster
    prometheus:
      enabled: true
      scrape_interval: 30s
    logging:
      enabled: true
      level: info
  EOT
}

