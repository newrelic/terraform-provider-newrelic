## 1.19.0 (June 05, 2020)

## 1.18.1 (June 03, 2020)

BUG FIXES:
* fix(alerts): update the client for epoch serialization fix ([#610](https://github.com/terraform-providers/terraform-provider-newrelic/issues/610))
* fix(test): Workloads returns ordered list of scope account IDs, update test

## 1.18.0 (May 15, 2020)

BUG FIXES:

* fix(client): update the client for pagination URL fix ([#549](https://github.com/terraform-providers/terraform-provider-newrelic/pull/549))
* fix(alerts): infra alert condition zero value detection ([#556](https://github.com/terraform-providers/terraform-provider-newrelic/pull/556))
* fix(alerts): allow blank runbook URL to be sent ([#563](https://github.com/terraform-providers/terraform-provider-newrelic/issues/563))

IMPROVEMENTS:

* feat(dashboard): add grid_column_count to dashboard schema ([#528](https://github.com/terraform-providers/terraform-provider-newrelic/pull/528))
* feat(application): implement newrelic_application_settings resource ([#558](https://github.com/terraform-providers/terraform-provider-newrelic/pull/558))

## 1.17.1 (May 04, 2020)

BUG FIXES:

* fix(client): update the client for pagination URL fix

## 1.17.0 (May 01, 2020)

IMPROVEMENTS:

* feat(dashboard): add grid_column_count to dashboard schema ([#535](https://github.com/terraform-providers/terraform-provider-newrelic/pull/535))

## 1.16.0 (March 24, 2020)

IMPROVEMENTS:

* feat(workloads): add a New Relic One workloads resource ([#474](https://github.com/terraform-providers/terraform-provider-newrelic/pull/474))
* docs(terraform-provider-newrelic): update getting started guide ([#473](https://github.com/terraform-providers/terraform-provider-newrelic/pull/473))
* docs(terraform-provider-newrelic): use correct default synthetics_api_url in config docs, remove inaccessible alert condition type ([#482](https://github.com/terraform-providers/terraform-provider-newrelic/pull/482))

## 1.15.1 (March 18, 2020)

BUG FIXES:

* fix(newrelic_alert_condition): fix import for terms threshold ([#469](https://github.com/terraform-providers/terraform-provider-newrelic/pull/469))

IMPROVEMENTS:

* docs(newrelic_alert_condition): remove `servers_metric` deprecated condition type ([#447](https://github.com/terraform-providers/terraform-provider-newrelic/pull/447))
* docs(terraform-provider-newrelic): add example modules ([#450](https://github.com/terraform-providers/terraform-provider-newrelic/pull/450))
* docs(terraform-provider-newrelic): add description to schema fields and enforce consistency ([#458](https://github.com/terraform-providers/terraform-provider-newrelic/pull/458))

## 1.15.0 (March 04, 2020)

BUG FIXES:

* fix(nrql_alert_condition): terms should be a TypeSet ([#421](https://github.com/terraform-providers/terraform-provider-newrelic/pull/421))

IMPROVEMENTS:

* feat(application_labels): Implementation of CRUD operations for NewRelic App Labels ([#417](https://github.com/terraform-providers/terraform-provider-newrelic/pull/417))
* feat(alert_policy_channels): add ability to add multiple channels to a policy ([#365](https://github.com/terraform-providers/terraform-provider-newrelic/pull/365))
* docs(newrelic_alert_condition): list of potential metrics for newrelic_alert_condition ([#431](https://github.com/terraform-providers/terraform-provider-newrelic/pull/431))

## 1.14.0 (February 20, 2020)

BUG FIXES:

* fix(provider): deprecate and re-enable the use of infra_api_url ([#411](https://github.com/terraform-providers/terraform-provider-newrelic/pull/411))

IMPROVEMENTS:

* feat(alerts): add ability to add multiple channels to a policy ([#398](https://github.com/terraform-providers/terraform-provider-newrelic/pull/398))
* feat(synthetics): add labels resource ([#407](https://github.com/terraform-providers/terraform-provider-newrelic/pull/407))
* feat(synthetics): add secure credentials resource ([#409](https://github.com/terraform-providers/terraform-provider-newrelic/pull/409))

## 1.13.1 (February 12, 2020)

BUG FIXES:

* fix(alert_condition): mark condition_scope optional
* fix(alert_channels): allow complex headers & payloads with new attributes
* fix(alert_channel): validate payload also has payload_type specified
* fix(newrelic_alert_channel): Force new resource for all config fields

IMPROVEMENTS:

* docs(alert_channel): add payload_type details to docs

## 1.13.0 (February 06, 2020)

BUG FIXES:
* fix: allow string representations of JSON for alert channel webhook and payload
* fix: clear client responses between pages

IMPROVEMENTS:

* feat: add `violation_close_timer` attribute to `newrelic_alert_condition` resource
* feat: add optional trace-level request logging
* docs: add debugging information to documentation website


## 1.12.2 (January 25, 2020)

BUG FIXES:

* fix: Error unmarshaling `newrelic_alert_channel` configuration headers and payload after release v1.12.1 ([#323](https://github.com/terraform-providers/terraform-provider-newrelic/issues/323))

## 1.12.1 (January 22, 2020)

IMPROVEMENTS:

* refactor: rebase `newrelic_alert_policy` resource on newrelic-client-go ([#341](https://github.com/terraform-providers/terraform-provider-newrelic/pull/341))
* refactor: migrate alert conditions to newrelic-client-go ([#338](https://github.com/terraform-providers/terraform-provider-newrelic/pull/338))
* docs: update alert-channel examples ([#325](https://github.com/terraform-providers/terraform-provider-newrelic/pull/325))

BUG FIXES:

* fix: Error unmarshaling `newrelic_alert_channel` configuration headers after release v1.12.0 ([#323](https://github.com/terraform-providers/terraform-provider-newrelic/issues/323))

## 1.12.0 (January 16, 2020)

IMPROVEMENTS:
* feat: deprecate the `configuration` attribute for `newrelic_alert_channel` ([#307](https://github.com/terraform-providers/terraform-provider-newrelic/pull/307))

BUG FIXES:
* fix: include `application_breakdown` as a valid visualization ([#305](https://github.com/terraform-providers/terraform-provider-newrelic/pull/305))

## 1.11.0 (January 09, 2020)
* feat: introduce new official New Relic client for Synthetics resource operations ([#294](https://github.com/terraform-providers/terraform-provider-newrelic/pull/294))

## 1.10.0 (December 18, 2019)

IMPROVEMENTS:
* feat: add ability to import `newrelic_synthetics_monitor` ([#267](https://github.com/terraform-providers/terraform-provider-newrelic/pull/267))
* docs: multiple improvements for readability and consistency

BUG FIXES:
* fix: add attribute validation for infra condition types ([#277](https://github.com/terraform-providers/terraform-provider-newrelic/pull/277))
* fix: loosen validation for threshold duration ([#277](https://github.com/terraform-providers/terraform-provider-newrelic/pull/277))
* fix: make event a computed attribute ([#277](https://github.com/terraform-providers/terraform-provider-newrelic/pull/277))

## 1.9.0 (December 05, 2019)

IMPROVEMENTS:

* feat: add `newrelic_plugins_alert_condition` resource ([#234](https://github.com/terraform-providers/terraform-provider-newrelic/pull/234))
* feat: add `newrelic_insights_event` resource ([#246](https://github.com/terraform-providers/terraform-provider-newrelic/pull/246))
* feat: add ability to import `newrelic_alert_channel` resource ([#241](https://github.com/terraform-providers/terraform-provider-newrelic/pull/241))
* feat: add ability to import `newrelic_alert_policy_channel` resource ([#249](https://github.com/terraform-providers/terraform-provider-newrelic/pull/249))
* feat: add ability to import `newrelic_infra_alert_condition` resource ([#254](https://github.com/terraform-providers/terraform-provider-newrelic/pull/254))
* feat: add ability to import `newrelic_nrql_alert_condition` resource for all condition types ([#250](https://github.com/terraform-providers/terraform-provider-newrelic/pull/250))
* feat: add `violation_time_limit_seconds` attribute to `newrelic_nrql_alert_condition` resource ([#198](https://github.com/terraform-providers/terraform-provider-newrelic/pull/198))
* docs: various improvements

BUG FIXES:

* fix: speed up `newrelic_application` data source state refresh for accounts with many applications ([#263](https://github.com/terraform-providers/terraform-provider-newrelic/pull/263))
* fix: `newrelic_alert_policy` data source now matches on policy name more strictly ([#197](https://github.com/terraform-providers/terraform-provider-newrelic/pull/197))
* fix: `newrelic_alert_channel` data source now matches on channel name more strictly ([#197](https://github.com/terraform-providers/terraform-provider-newrelic/pull/197))

## 1.8.0 (November 22, 2019)

IMPROVEMENTS:

* feat: add import functionality for `newrelic_alert_policy` data source ([#233](https://github.com/terraform-providers/terraform-provider-newrelic/pull/233))
* feat: allow passing an http transport ([#228](https://github.com/terraform-providers/terraform-provider-newrelic/pull/228))
* docs: add Getting Started section ([#225](https://github.com/terraform-providers/terraform-provider-newrelic/pull/225))
* docs: update alert infra condition docs to include runbook url argument ([#211](https://github.com/terraform-providers/terraform-provider-newrelic/pull/211))
* docs: add info for `nrql_alert_condition` arguments `type`, `expected_groups`, and `ignore_overlap` ([#231](https://github.com/terraform-providers/terraform-provider-newrelic/pull/231))
* fix: get `make` working again ([#236](https://github.com/terraform-providers/terraform-provider-newrelic/pull/236))


## 1.7.0 (November 13, 2019)

BUG FIXES:

* fix: align alert condition duration constraints to NR's API constraints ([#201](https://github.com/terraform-providers/terraform-provider-newrelic/issues/201))
* fix: align alert policy validation with NR's API validation ([#199](https://github.com/terraform-providers/terraform-provider-newrelic/issues/199))

IMPROVEMENTS:
* Dashboard improvements ([#206](https://github.com/terraform-providers/terraform-provider-newrelic/pull/206))
    * Support for more dashboard widget types:
        * `metric_line_chart`
        * `markdown`
        * `gauge`
        * `billboard`
        * `billboard_comparison`
    * Plan-time validation for:
        * `icon`
        * `visualization`
    * More robust validation of widgets, based on visualization type
    * Allow up to 300 dashboard widgets per the [API documentation]
    * Refresh dashboard state properly when underlying resource has been deleted
    * Documentation improvements
* Adds ability to skip TLS verification from a remote agent and trust self-signed certs ([#196](https://github.com/terraform-providers/terraform-provider-newrelic/pull/196))

NOTES:
* Documentation updates ([#207](https://github.com/terraform-providers/terraform-provider-newrelic/pull/207), [#195](https://github.com/terraform-providers/terraform-provider-newrelic/pull/195))



## 1.6.0 (November 07, 2019)

BUG FIXES:

* Perpetual drift in alert conditions with multiple entities ([#137](https://github.com/terraform-providers/terraform-provider-newrelic/issues/137))

IMPROVEMENTS

* Add provider version to the User-Agent ([#189](https://github.com/terraform-providers/terraform-provider-newrelic/pull/189))
* Add support for outlier NRQL alert conditions ([#141](https://github.com/terraform-providers/terraform-provider-newrelic/pull/141))
* Update module paultyng/go-newrelic/v4 to v4.6.0 ([#187](https://github.com/terraform-providers/terraform-provider-newrelic/pull/187))

## 1.5.2 (October 23, 2019)

BUG FIXES:

* `newrelic_nrql_alert_condition` modifies `duration` validation to match rest API's requirements ([#169](https://github.com/terraform-providers/terraform-provider-newrelic/issues/169))

## 1.5.1 (July 11, 2019)

BUG FIXES:

* `newrelic_nrql_alert_condition` modifies `since_value` validation to match rest API's requirements ([#144](https://github.com/terraform-providers/terraform-provider-newrelic/issues/144))

## 1.5.0 (March 26, 2019)

FEATURES

* **New Resource** `newrelic_synthetics_monitor` and `newrelic_synthetics_monitor_script` ([#67](https://github.com/terraform-providers/terraform-provider-newrelic/issues/67))

IMPROVEMENTS

* Add Terraform 0.12 support ([#107](https://github.com/terraform-providers/terraform-provider-newrelic/issues/107))

## 1.4.0 (February 27, 2019)

IMPROVEMENTS:

* `newrelic_alert_condition` make enabled status configurable ([#70](https://github.com/terraform-providers/terraform-provider-newrelic/issues/70))
* `newrelic_alert_condition` add name length validation ([#79](https://github.com/terraform-providers/terraform-provider-newrelic/issues/79))

## 1.3.0 (February 07, 2019)

FEATURES:

* **New Data Source** `newrelic_alert_policy` ([#64](https://github.com/terraform-providers/terraform-provider-newrelic/issues/64))

BUG FIXES:
* `newrelic_alert_policy` should have update functionality ([#68](https://github.com/terraform-providers/terraform-provider-newrelic/pull/68))
* Fix documentation typos for `newrelic_nrql_alert_condition` ([#76](https://github.com/terraform-providers/terraform-provider-newrelic/issues/76))
* Fix diff problem with `newrelic_alert_condition.term` ([#63](https://github.com/terraform-providers/terraform-provider-newrelic/issues/63))

## 1.2.0 (November 02, 2018)

FEATURES:

* **New Data Source** `newrelic_alert_policy` ([#34](https://github.com/terraform-providers/terraform-provider-newrelic/issues/34))

IMPROVEMENTS:

* `newrelic_infra_alert_condition`: Add support for `integration_provider` ([#48](https://github.com/terraform-providers/terraform-provider-newrelic/issues/48))
* `newrelic_dashboard`: Add support for `filter` ([#46](https://github.com/terraform-providers/terraform-provider-newrelic/issues/46))

## 1.1.0 (October 16, 2018)

FEATURES:

* **New Resource** `newrelic_synthetics_alert_condition` ([#22](https://github.com/terraform-providers/terraform-provider-newrelic/pull/22))
* **New Data Source** `newrelic_synthetics_monitor` ([#22](https://github.com/terraform-providers/terraform-provider-newrelic/pull/22))

BUG FIXES:

* `newrelic_alert_policy` int64 bug ([#42](https://github.com/terraform-providers/terraform-provider-newrelic/pull/42))
* Missing doc links ([#49](https://github.com/terraform-providers/terraform-provider-newrelic/pull/49))

## 1.0.1 (June 06, 2018)

FEATURES:

* **New Resource** `newrelic_infra_alert_condition` ([#30](https://github.com/terraform-providers/terraform-provider-newrelic/pull/30))

## 1.0.0 (February 12, 2018)

FEATURES:

* **New Resource** `newrelic_dashboard` ([#26](https://github.com/terraform-providers/terraform-provider-newrelic/pull/26))
* **New Data Source** `newrelic_key_transaction` ([#21](https://github.com/terraform-providers/terraform-provider-newrelic/pull/21))

IMPROVEMENTS:

* resource/newrelic_alert_condition: Add support for `apm_jvm_metric` and instance scope alerts ([#24](https://github.com/terraform-providers/terraform-provider-newrelic/pull/24))

## 0.1.1 (August 02, 2017)

FEATURES:

* **New Resource:** `newrelic_nrql_alert_condition` ([#15](https://github.com/terraform-providers/terraform-provider-newrelic/issues/15))

IMPROVEMENTS:

* resource/newrelic_alert_condition: Allow zero threshold value for terms ([#13](https://github.com/terraform-providers/terraform-provider-newrelic/issues/13))

## 0.1.0 (June 21, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
