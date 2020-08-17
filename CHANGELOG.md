<a name="v2.5.1"></a>
## [v2.5.1] - 2020-08-14
### Bug Fixes
- cannot create resource "newrelic_infra_alert_condition" of type "infra_host_not_reporting"
- **infra:** avoid nil pointer reference
- **infra:** avoid nil pointer reference

<a name="v2.5.0"></a>
## [v2.5.0] - 2020-08-03
### Bug Fixes
- **alert_policy:** avoid drift due to account_id inheritance in resource and data source
- **nrql_alert_condition:** avoid drift due to account_id inheritance in NRQL alert condition

### Documentation Updates
- **synthetics_monitor_location:** Adding docs for Synthetics monitor location data source.

### Features
- **synthetics_monitor_location:** Add data source newrelic_synthetics_monitor_location.

<a name="v2.4.2"></a>
## [v2.4.2] - 2020-07-30
### Documentation Updates
- **dashboard:** Improve docs for limit and order_by widget attributes

<a name="v2.4.1"></a>
## [v2.4.1] - 2020-07-29
### Bug Fixes
- **alerts:** flatten condition scope properly for APM JVM metrics
- **newrelic_alert_condition:** allow instance scope for JVM app metrics

<a name="v2.4.0"></a>
## [v2.4.0] - 2020-07-28
### Bug Fixes
- **alerts:** Unify how alert policy selects an account_id
- **infra_alert_condition:** support zero-value thresholds for infra_alert_condition resource

### Documentation Updates
- **alert_policy:** update alert_policy import section, add  default to arg ref

### Features
- **infra_alert_condition:** add description attribute

<a name="v2.3.0"></a>
## [v2.3.0] - 2020-07-23
### Features
- add a newrelic_account data source

<a name="v2.2.1"></a>
## [v2.2.1] - 2020-07-10
### Bug Fixes
- replacement for deadlink linter
- replacement for deadlink linter
- **alert_condition:** remove conditional to fix drift when using 'user_defined' attributes

### Documentation Updates
- fix broken links
- fix broken links
- fix broken links
- fix broken links
- communicate that most but not all keys have prefixes
- **alerts:** update documentation for newrelic_nrql_alert_condition

<a name="v2.2.0"></a>
## [v2.2.0] - 2020-07-08
### Bug Fixes
- **docs:** extra whitespace below table
- **docs:** better table header rendering
- **nrql_alert_condition:** use better term operator

### Documentation Updates
- **alerts:** include account_id attribute for alert_policy

### Features
- **alerts:** new newrelic_alerts_location_failure_condition resource

<a name="v2.1.2"></a>
## [v2.1.2] - 2020-06-26
### Bug Fixes
- **alerts:** require at least one violation time limit attr
- **alerts:** improve nil handling for alert_channel

### Documentation Updates
- **provider:** additional v2 updates, migration guide updates
- **provider:** add getting started guide to the quick links
- **provider:** fix incorrect newrelic_application reference in some examples
- **provider:** add account_id to argument reference, move argument reference above the fold
- **provider:** add environment variables and schema attribute table
- **provider:** update getting started example to reflect v2 updates
- **readme:** update title, add link to latest documentation

<a name="v2.1.1"></a>
## [v2.1.1] - 2020-06-23
### Features
- update the release process to prepare for repo handoff

<a name="v2.1.0"></a>
## [v2.1.0] - 2020-06-22
### Documentation Updates
- include information on pinning a version
- include sidebar link for 2.x upgrade

### Features
- **eventstometrics:** add an events to metrics rule resource ([#690](https://github.com/newrelic/terraform-provider-newrelic/issues/690))

<a name="v2.0.0"></a>
## [v2.0.0] - 2020-06-18
### Bug Fixes
- Require condition_scope = `instance` for validation_close_timer
- Add validation to newrelic_alert_condtion condition_scope
- **alerts:** remove DiffSuppressFunc on TypeSet to avoid test drift
- **alerts:** infra alert condition zero value detection
- **alerts:** handle a nil reference with more grace
- **application_settings:** Remove delete, as it is not possible
- **deps:** Revert terraform sdk to 1.10.0
- **newrelic:** fix the failing integration tests ([#519](https://github.com/newrelic/terraform-provider-newrelic/issues/519))
- **nrql_alert_condition:** threshold_occurrences is case insensitive, attribute description updates

### Documentation Updates
- add callout to top of each v1.x doc page
- tidy up after review
- DEPRECATION notice for 1.x
- update index header with improved words
- update getting started guide to reference new material
- update README with new pointers
- add table for current endpoint in use per resource
- include documentation about upgrading the provider to 2.x
- update API key references to match desires
- include v1 index.html in sidebar
- prep for v2.x, isolate v1.x docs
- **alert_channel:** fix broken 'nested config' anchor link
- **alerts:** include caveat about NRQL alerts condition operator usage with outliers
- **alerts:** update wording to avoid implementation details
- **alerts:** include deprecation notice for "terms"
- **alerts:** update examples to reflect deprecation
- **getting started:** fix resource naming
- **nrql_alert_condition:** add outlier example, add new attributes, deprecate old attributes, update import section
- **nrql_alert_condition:** update docs to reflect version 2.0 changes
- **provider:** add region to provider docs, removing references to API base URLs
- **provider:** add provider configuration guide page
- **workloads:** fix api key attribute name ([#489](https://github.com/newrelic/terraform-provider-newrelic/issues/489))

### Features
- **alerts:** convert Alerts Policies to nerdgraph
- **application:** implement newrelic_application resource
- **dashboard:** add grid_column_count to dashboard schema
- **entity_tags:** add an entity tag resource ([#679](https://github.com/newrelic/terraform-provider-newrelic/issues/679))
- **nrql_alert_condition:** integrate nerdgraph for nrql alert conditions
- **provider:** add region to provider schema, handle API URLs based off region

<a name="v1.20.1"></a>
## [v1.20.1] - 2020-07-27
### Bug Fixes
- **infra_alert_condition:** [v1.x] support zero-value thresholds for infra_alert_condition resource

<a name="v1.20.0"></a>
## [v1.20.0] - 2020-07-23
<a name="v1.19.1"></a>
## [v1.19.1] - 2020-06-24
### Bug Fixes
- **changelog:** remove 1.18.1 from changelog, 1.19.0 is the replacement

### Features
- update the release process to prepare for repo handoff

<a name="v1.19.0"></a>
## [v1.19.0] - 2020-06-05
### Bug Fixes
- **test:** Workloads returns ordered list of scope account IDs, update test

### Documentation Updates
- **application_settings:** add application settings resource to sidebar ([#582](https://github.com/newrelic/terraform-provider-newrelic/issues/582))

<a name="v1.18.0"></a>
## [v1.18.0] - 2020-05-15
### Bug Fixes
- **alerts:** infra alert condition zero value detection

### Features
- **application:** implement newrelic_application resource ([#558](https://github.com/newrelic/terraform-provider-newrelic/issues/558))

<a name="v1.17.1"></a>
## [v1.17.1] - 2020-05-04
### Bug Fixes
- **client:** update the client for pagination URL fix

<a name="v1.17.0"></a>
## [v1.17.0] - 2020-05-01
### Features
- **dashboard:** add grid_column_count to dashboard schema

<a name="v1.16.0"></a>
## [v1.16.0] - 2020-03-24
### Documentation Updates
- use correct default synthetics_api_url in config docs, remove inaccessible alert condition type
- Update getting started guide

### Features
- **workloads:** add a workloads resource ([#474](https://github.com/newrelic/terraform-provider-newrelic/issues/474))

<a name="v1.15.1"></a>
## [v1.15.1] - 2020-03-18
### Bug Fixes
- import condition terms regardless of threshold format ([#469](https://github.com/newrelic/terraform-provider-newrelic/issues/469))

### Documentation Updates
- ensure consistency ([#458](https://github.com/newrelic/terraform-provider-newrelic/issues/458))
- **examples:** add a golden signal alerting module example ([#450](https://github.com/newrelic/terraform-provider-newrelic/issues/450))

<a name="v1.15.0"></a>
## [v1.15.0] - 2020-03-04
### Bug Fixes
- **application_label:** use correct type assertions for applications and servers attributes
- **nrql_alert_condition:** terms should be a TypeSet

### Documentation Updates
- **alert_policy_channel:** include sorting recommendation for channel_ids

### Features
- **alert_policy_channels:** add ability to add multiple channels to a policy

<a name="v1.14.0"></a>
## [v1.14.0] - 2020-02-20
### Bug Fixes
- **provider:** deprecate and re-enable the use of infra_api_url ([#411](https://github.com/newrelic/terraform-provider-newrelic/issues/411))

### Features
- **alert_policy:** add ability to add multiple channels to a policy ([#398](https://github.com/newrelic/terraform-provider-newrelic/issues/398))
- **synthetics:** add secure credentials resource ([#409](https://github.com/newrelic/terraform-provider-newrelic/issues/409))
- **synthetics:** add labels resource ([#407](https://github.com/newrelic/terraform-provider-newrelic/issues/407))

<a name="v1.13.1"></a>
## [v1.13.1] - 2020-02-12
### Bug Fixes
- **alert_channel:** validate payload also has payload_type specified
- **alert_channels:** allow complex headers & payloads with new attributes
- **alert_condition:** mark condition_scope optional
- **newrelic_alert_channel:** Force new resource for all config fields

### Documentation Updates
- **alert_channel:** add payload_type details to docs

<a name="v1.13.0"></a>
## [v1.13.0] - 2020-02-06
### Documentation Updates
- Make a note about community resources and support
- Make note about ignoring secrets

### Features
- replace provider backend with newrelic-client-go ([#358](https://github.com/newrelic/terraform-provider-newrelic/issues/358))
- **infra_alert_condition:** add violation_close_timer to newrelic_infra_alert_condition resource ([#370](https://github.com/newrelic/terraform-provider-newrelic/issues/370))

<a name="v1.12.2"></a>
## [v1.12.2] - 2020-01-25
### Bug Fixes
- **alert_channels:** handle more complex JSON structures in payload or headers ([#361](https://github.com/newrelic/terraform-provider-newrelic/issues/361))

<a name="v1.12.1"></a>
## [v1.12.1] - 2020-01-22
### Bug Fixes
- **newrelic-client-go:** Fix API Key passing to provider

### Documentation Updates
- update alert-channel examples

<a name="v1.12.0"></a>
## [v1.12.0] - 2020-01-16
### Bug Fixes
- **dashboards:** include application_breakdown as a valid visualization

### Documentation Updates
- **alerts:** update documentation for newrelic_alert_channel
- **dashboards:** include application_breakdown in docs

### Features
- **alerts:** deprecate alerts channel configuration and add config block

<a name="v1.11.0"></a>
## [v1.11.0] - 2020-01-09
### Documentation Updates
- update docs for consistency
- document the new synthetics_api_url variable

### Features
- release 1.11.0
- update CHANGELOG for v1.11.0

<a name="v1.10.0"></a>
## [v1.10.0] - 2019-12-18
### Bug Fixes
- make event a computed attribute
- loosen validation for threshold duration
- add attribute validation for infra condition types

### Documentation Updates
- update documentation for newrelic_infra_alert_condition
- update newrelic_synthetics_monitor docs
- add missing resources and data source to sidebar
- updates for consistency

### Features
- add ability to import resource_newrelic_synthetics_monitor, update acceptance tests and add coverage

<a name="v1.9.0"></a>
## [v1.9.0] - 2019-12-05
### Bug Fixes
- use name as filter in application lookup
- fix newrelic_infra_alert imports and backfill acc testing

### Documentation Updates
- update for clarity and consistency
- update nrql_alert_condition docs to reference violation_time_limit_seconds
- update docs for newrelic_nrql_alert_condition
- refresh the infra alert condition docs
- add docs for newrelic_plugin_component
- update docs for newrelic_alert_channel resource and data source
- fix formatting in dashboard docs

### Features
- allow importing of violation_time_limit_seconds, add validation, remove inline docs
- add ability to import nrql_alert_condition for types static and outlier
- update newrelic_synthetics_alert_condition  acceptance tests
- update newrelic_synthetics_monitor_script acceptance tests
- add a plugin component data source
- create importer for alert policy channels
- add ability to import newrelic_alert_channel data source

<a name="v1.8.0"></a>
## [v1.8.0] - 2019-11-22
### Bug Fixes
- appease golangci-lint when running make

### Documentation Updates
- add Getting Started section

### Features
- add import functionality for newrelic_alert_policy data source

<a name="v1.7.0"></a>
## [v1.7.0] - 2019-11-13
### Bug Fixes
- align alert condition duration constraints to NR's API constraints
- align alert policy validation with NR's API validation
- lint issue, update modules
- merge conflicts
- typos

<a name="v1.6.0"></a>
## [v1.6.0] - 2019-11-07
<a name="v1.5.2"></a>
## [v1.5.2] - 2019-10-23
<a name="v1.5.1"></a>
## [v1.5.1] - 2019-07-11
<a name="v1.5.0"></a>
## [v1.5.0] - 2019-03-26
<a name="v1.4.0"></a>
## [v1.4.0] - 2019-02-27
<a name="v1.3.0"></a>
## [v1.3.0] - 2019-02-07
<a name="v1.2.0"></a>
## [v1.2.0] - 2018-11-02
<a name="v1.1.0"></a>
## [v1.1.0] - 2018-10-16
<a name="v1.0.1"></a>
## [v1.0.1] - 2018-06-06
<a name="v1.0.0"></a>
## [v1.0.0] - 2018-02-12
<a name="v0.1.1"></a>
## [v0.1.1] - 2017-08-02
<a name="v0.1.0"></a>
## v0.1.0 - 2017-06-21
[Unreleased]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.5.1...HEAD
[v2.5.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.5.0...v2.5.1
[v2.5.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.4.2...v2.5.0
[v2.4.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.4.1...v2.4.2
[v2.4.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.4.0...v2.4.1
[v2.4.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.3.0...v2.4.0
[v2.3.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.2.1...v2.3.0
[v2.2.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.2.0...v2.2.1
[v2.2.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.1.2...v2.2.0
[v2.1.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.1.1...v2.1.2
[v2.1.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.1.0...v2.1.1
[v2.1.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.0.0...v2.1.0
[v2.0.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.20.1...v2.0.0
[v1.20.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.20.0...v1.20.1
[v1.20.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.19.1...v1.20.0
[v1.19.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.19.0...v1.19.1
[v1.19.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.18.0...v1.19.0
[v1.18.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.17.1...v1.18.0
[v1.17.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.17.0...v1.17.1
[v1.17.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.16.0...v1.17.0
[v1.16.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.15.1...v1.16.0
[v1.15.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.15.0...v1.15.1
[v1.15.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.14.0...v1.15.0
[v1.14.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.13.1...v1.14.0
[v1.13.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.13.0...v1.13.1
[v1.13.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.12.2...v1.13.0
[v1.12.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.12.1...v1.12.2
[v1.12.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.12.0...v1.12.1
[v1.12.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.11.0...v1.12.0
[v1.11.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.10.0...v1.11.0
[v1.10.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.9.0...v1.10.0
[v1.9.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.8.0...v1.9.0
[v1.8.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.7.0...v1.8.0
[v1.7.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.6.0...v1.7.0
[v1.6.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.5.2...v1.6.0
[v1.5.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.5.1...v1.5.2
[v1.5.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.5.0...v1.5.1
[v1.5.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.4.0...v1.5.0
[v1.4.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.3.0...v1.4.0
[v1.3.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.2.0...v1.3.0
[v1.2.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.1.0...v1.2.0
[v1.1.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.0.1...v1.1.0
[v1.0.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.1.1...v1.0.0
[v0.1.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.1.0...v0.1.1
