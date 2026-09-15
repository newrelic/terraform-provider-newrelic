# ── Rule 1: APM alert coverage ────────────────────────────────────────────────
# Drift scenario A: toggle enabled = true → false in the UI, then terraform plan
# should detect the disabled state and plan to re-enable it.
resource "newrelic_scorecard_rule" "alert_coverage" {
  name         = "pnandula-drift-rule-alert-coverage"
  description  = "APM services must have alerts configured"
  enabled      = true
  run_interval = 1440

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(latest(alertSeverity) != 'NOT_CONFIGURED', 1, 0) AS 'score' FROM Entity WHERE type = 'APM-APPLICATION' FACET id LIMIT MAX SINCE 1 day ago"
  }

  tags = ["purpose:drift-test", "owner:pnandula"]
}

# ── Rule 2: Latency threshold ─────────────────────────────────────────────────
# Drift scenario B: change description in the UI, then terraform plan
# should surface the description change and plan to restore it.
resource "newrelic_scorecard_rule" "latency_threshold" {
  name         = "pnandula-drift-rule-latency"
  description  = "p95 latency must be under 500ms"
  enabled      = true
  run_interval = 720

  nrql_engine {
    accounts = [var.account_id]
    query    = "SELECT if(percentile(duration, 95) < 0.5, 1, 0) AS 'score' FROM Transaction FACET entity.guid LIMIT MAX SINCE 1 hour ago"
  }

  tags = ["purpose:drift-test"]
}

# ── Scorecard ─────────────────────────────────────────────────────────────────
# Drift scenario C: add a tag in the UI (e.g. "team:platform"), then
# terraform plan should detect the extra tag and plan to remove it.
resource "newrelic_scorecard" "engineering" {
  name        = "pnandula-drift-scorecard-engineering"
  description = "Engineering quality scorecard for drift testing"
  tags        = ["env:test", "owner:pnandula"]

  progress_levels {
    id             = "red"
    name           = "Needs Work"
    description    = "Score below 60%"
    hex_color_code = "#FF4444"
  }
  progress_levels {
    id             = "amber"
    name           = "Improving"
    description    = "Score between 60-80%"
    hex_color_code = "#FFAA00"
  }
  progress_levels {
    id             = "green"
    name           = "Healthy"
    description    = "Score above 80%"
    hex_color_code = "#00CC44"
  }

  rule_ids = [
    newrelic_scorecard_rule.alert_coverage.id,
    newrelic_scorecard_rule.latency_threshold.id,
  ]
}

# ── Outputs (handy for locating resources in the UI) ──────────────────────────
output "scorecard_guid" {
  value       = newrelic_scorecard.engineering.id
  description = "GUID of the scorecard — use this to find it in the NR UI"
}

output "rule_alert_coverage_guid" {
  value       = newrelic_scorecard_rule.alert_coverage.id
  description = "GUID of the alert-coverage rule"
}

output "rule_latency_guid" {
  value       = newrelic_scorecard_rule.latency_threshold.id
  description = "GUID of the latency rule"
}
