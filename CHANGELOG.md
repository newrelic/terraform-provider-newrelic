<a name="v2.19.1"></a>
## [v2.19.1] - 2021-02-23
### Bug Fixes
- **deps:** update module gotest.tools/gotestsum to v1.6.2
- **deps:** update module github.com/golangci/golangci-lint to v1.37.1
- **deps:** Update module newrelic/newrelic-client-go to v0.58.2
- **deps:** update module goreleaser/goreleaser to v0.157.0

<a name="v2.19.0"></a>
## [v2.19.0] - 2021-02-18
### Bug Fixes
- **one_dashboard:** Table Widget should have filter on them
- **one_dashboard:** Inherit nrql_query account_id from dashboard by default

### Documentation Updates
- update changelog
- update changelog

### Features
- **one_dashboard:** Add support for widget_histogram
- **one_dashboard:** Add support for widget_funnel
- **one_dashboard:** Add support for widget_bullet
- **one_dashboard:** Add widget_heatmap

<a name="v2.18.0"></a>
## [v2.18.0] - 2021-02-09
### Bug Fixes
- **alert_muting_rule:** update test expectation to match input
- **alert_muting_rule:** condition tag validation

### Features
- **alert_muting_rule:** add schedule support

<a name="v2.17.0"></a>
## [v2.17.0] - 2021-02-01
### Bug Fixes
- **nrql_alert_condition:** validate operator based on condition type

### Documentation Updates
- **one_dashboard:** add linked_entity_guids to newrelic_one_dashboard resource docs

### Features
- **one_dashboard:** add linked entities to widget schema

<a name="v2.16.0"></a>
## [v2.16.0] - 2021-01-29
### Documentation Updates
- fix broken links in api_access_key.html.markdown
- **nrql_alert_condition:** Amends threshold_duration constraints for NRQL alert conditions

<a name="v2.15.1"></a>
## [v2.15.1] - 2021-01-22
### Documentation Updates
- **one_dashboard:** remove unused entity reference in example

<a name="v2.15.0"></a>
## [v2.15.0] - 2021-01-14
### Documentation Updates
- update changelog
- **one_dashboard:** Add overview doc for one_dashboard resource

### Features
- **one_dashboard:** Testing out one_dashboard resource

<a name="v2.14.1"></a>
## [v2.14.1] - 2021-01-12
### Documentation Updates
- change personal API key to user api key
- update API key instructions for getting started guide

<a name="v2.14.0"></a>
## [v2.14.0] - 2020-12-09
### Bug Fixes
- **infra_alert_condition:** fix integration tests

### Documentation Updates
- update getting started guide with a link to EU graphiql
- **nrql_alert_condition:** include notes about upgrading from 1.x

### Features
- **nrql_alert_condition:** swap deprecation of violation_time_limit fields

<a name="v2.13.5"></a>
## [v2.13.5] - 2020-11-13
### Bug Fixes
- **nrql_alert_condition:** reverse attribute detection for migration

<a name="v2.13.4"></a>
## [v2.13.4] - 2020-11-11
### Bug Fixes
- **docs:** Alert Channels do not manage Policies
- **newrelic_entity:** include additional ID attr for browser apps

### Documentation Updates
- include note about API key access

<a name="v2.13.3"></a>
## [v2.13.3] - 2020-10-27
### Bug Fixes
- **nrql_alert_condition:** fix fill_option DiffSuppressFunc

### Documentation Updates
- **alert_condition:** document apm_jvm_metric

<a name="v2.13.2"></a>
## [v2.13.2] - 2020-10-26
### Documentation Updates
- **alert_policy_channel:** update example reference

<a name="v2.13.1"></a>
## [v2.13.1] - 2020-10-19
<a name="v2.13.0"></a>
## [v2.13.0] - 2020-10-16
### Documentation Updates
- **newrelic_synthetics_monitor_script:** Use file method instead of template_file data source

### Features
- **client:** update newrelic-client-go (retry on nerdgraph timeouts)

<a name="v2.12.1"></a>
## [v2.12.1] - 2020-10-15
### Bug Fixes
- **dashboard:** use state migration to fix 500 error when upgrading from v2.7.5 to v2.8 and beyond
- **nrql_alert_condition:** avoid drift using computed value

### Documentation Updates
- add instructions for New Relic One users to get an api key
- **dashboard:** update docs regarding cross-account widget config drift

<a name="v2.12.0"></a>
## [v2.12.0] - 2020-10-08
### Features
- **alerts:** allow a 30 day violation limit for nrql conditions

<a name="v2.11.1"></a>
## [v2.11.1] - 2020-10-07
### Documentation Updates
- add website documentation for nrql_alert aggregation_window

<a name="v2.11.0"></a>
## [v2.11.0] - 2020-10-06
### Features
- **aggregation_window:** add support for nrql signal aggregationWindow

<a name="v2.10.3"></a>
## [v2.10.3] - 2020-10-05
### Documentation Updates
- remove admin key from documentation

<a name="v2.10.2"></a>
## [v2.10.2] - 2020-10-02
### Bug Fixes
- **build:** update version.ProviderVersion via ldflags during release process

### Documentation Updates
- remove admin API key from docs and various other updates
- update changelog
- change slack integration documentation
- update process running example
- **dashboard:** fix some broken links
- **synthetics:** remove newrelic_synthetics_label resource

### Features
- **alerts:** deprecate plugins conditions and un-deprecate APM alert conditions
- **synthetics:** replace REST API calls with Nerdgraph calls

<a name="v2.9.0"></a>
## [v2.9.0] - 2020-10-01
### Documentation Updates
- update changelog

### Features
- **dashboard:** enable Personal API Key auth for dashboards and some sythentics resources

<a name="v2.8.0"></a>
## [v2.8.0] - 2020-09-30
### Documentation Updates
- update infra alert condition api key type
- update changelog
- update development instructions for new TF version
- DEPRECATION notice for newrelic_alert_condition
- update supported Go information and test config
- **README:** update provider configuration pin version examples
- **dashboard:** update docs with info regarding widget.account_id and cross-account widgets
- **dashboard:** add cross-account example

### Features
- **dashboard:** support cross-account widgets :)

<a name="v2.7.5"></a>
## [v2.7.5] - 2020-09-23
### Bug Fixes
- **entity:** add VIZ domain

### Documentation Updates
- update changelog
- **nrql_condition:** add clarity around choosing between new and old/deprecated attributes
- **nrql_condition:** clarify when value_function attr is 'required' vs 'not required'

<a name="v2.7.4"></a>
## [v2.7.4] - 2020-09-18
### Bug Fixes
- **nrql_alert_condition:** update validation for nrql conditions

<a name="v2.7.3"></a>
## [v2.7.3] - 2020-09-17
### Bug Fixes
- **alerts:** avoid bad index reference

<a name="v2.7.2"></a>
## [v2.7.2] - 2020-09-16
### Documentation Updates
- update changelog

<a name="v2.7.1"></a>
## [v2.7.1] - 2020-09-11
### Bug Fixes
- **nrql_alert_condition:** Fixed an issue with extrapolation (gap filling) settings

### Documentation Updates
- fix references to newrelic_entity data sources
- update authentication table
- replace uses of APM conditions with NRQL conditions
- update changelog

<a name="v2.7.0"></a>
## [v2.7.0] - 2020-09-04
### Documentation Updates
- update changelog

### Features
- **nrql_alert_condition:** Added support for expiration (loss of signal) and extrapolation (gap filling) settings

<a name="v2.6.1"></a>
## [v2.6.1] - 2020-09-03
### Bug Fixes
- **changelog:** ensure proper branch to base from
- **nrql_alert_condition:** add missing zeros to violation_time_limit_seconds to the new:old map

<a name="v2.6.0"></a>
## [v2.6.0] - 2020-08-24
### Bug Fixes
- **alert_channel:** avoid drift with config.auth_password
- **alert_channel:** avoid config drift with sensitive values not returned by the API
- **alerts:** ensure threshold_occurrences case fold comparison
- **changelog:** update changelog on release only, drop reviewer spec
- **nrql_alert_condition:** fix drift with threshold_occurrences - store lowercase in terraform state

### Documentation Updates
- **alert_channel:** add note to import section regarding handling of sensitive data
- **alert_muting_rule:** Added docs for alert muting rule.

### Features
- **alert_muting_rule:** Creating alert muting rule resource.
- **newrelic_api_access_key:** Implement new resource: newrelic_api_access_key

<a name="v2.5.1"></a>
## [v2.5.1] - 2020-08-17
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
<a name="v0.11.0"></a>
## [v0.11.0] - 2020-02-27
### Features
- **http:** allow personal API keys to be used for alerts and APM resources

<a name="v0.10.1"></a>
## [v0.10.1] - 2020-02-20
### Bug Fixes
- **entities:** tags filter needs to use type TagValue in graphql query
- **newrelic:** Add option to set ServiceName in Config

<a name="v0.10.0"></a>
## [v0.10.0] - 2020-02-19
### Features
- **ci:** add release make target
- **ci:** the beginnings of some release automation
- **synthetics:** add secure credentials resource
- **synthetics:** implement label monitor support

<a name="v0.9.0"></a>
## [v0.9.0] - 2020-02-05
### Bug Fixes
- allow string representations of JSON for alert channel webhook and payload
- **http:** Clear client responses between pages

### Features
- **alerts:** Implement multi-location synthetics conditions
- **http:** add trace logging with additional request info

<a name="v0.8.0"></a>
## [v0.8.0] - 2020-01-29
### Bug Fixes
- **alerts:** ensure multiple channels can be added via /alerts_policy_channel.json endpoint ([#114](https://github.com/newrelic/terraform-provider-newrelic/issues/114))

### Features
- **apm:** Add support application metric names and data

<a name="v0.7.1"></a>
## [v0.7.1] - 2020-01-24
### Bug Fixes
- **alerts:** handle more complex JSON structures in headers and/or payload
- **logging:** use global methods for the default logger rather than a logrus instance

<a name="v0.7.0"></a>
## [v0.7.0] - 2020-01-23
### Features
- **newrelic:** add ConfigOptions for logging
- **newrelic:** add the ability to configure base URLs per API

<a name="v0.6.0"></a>
## [v0.6.0] - 2020-01-22
### Features
- **alerts:** add GetSyntheticsCondition method ([#105](https://github.com/newrelic/terraform-provider-newrelic/issues/105))

<a name="v0.5.1"></a>
## [v0.5.1] - 2020-01-21
### Bug Fixes
- **alerts:** custom unmarshal of channel configuration Headers and Payload fields ([#102](https://github.com/newrelic/terraform-provider-newrelic/issues/102))

<a name="v0.5.0"></a>
## [v0.5.0] - 2020-01-16
### Documentation Updates
- **newrelic:** update API key configuration documentation

<a name="v0.4.0"></a>
## [v0.4.0] - 2020-01-15
### Bug Fixes
- retry HTTP requests on 429 status codes

### Features
- **entities:** add entities search and entity tagging

<a name="v0.3.0"></a>
## [v0.3.0] - 2020-01-13
### Bug Fixes
- make use of ErrorNotFound type for Get methods that are based on List methods
- add policy ID to alert condition

### Documentation Updates
- update example
- **build:** Update README for commit message format
- **changelog:** Add auto-generation of CHANGELOG from git comments via `make changelog`

### Features
- add top-level logging package for convenience
- add option for JSON logging and fail gracefully when log level cannot be parsed
- introduce logging
- update monitor scripts with return design pattern, update tests

<a name="v0.2.0"></a>
## [v0.2.0] - 2020-01-08
### Bug Fixes
- rename variables to fix redeclared error
- update unit tests to use new method sigs
- fix monitor ID type and GetMonitor URL
- http client needs to handle other 'success' response status codes such as 201
- add godoc as a dep, and a warning about GOPATH and godoc
- fix paging bug for v2 API
- **lint:** formatting fixes for linter

### Documentation Updates
- update readme example
- add alerts package docs
- temporarily checking in broken import paths in generated markdown docs
- add inline documentation
- add badges to README
- fill in missing inline documentation
- document some methods

### Features
- add DeletePluginCondition
- add CreatePluginCondition
- add UpdatePluginCondition
- add GetPluginCondition
- add ListPluginsConditions
- encode monitor script text
- add ability to use 'detailed' query param in ListPlugins method
- add GetPlugin
- add ListPlugins
- publicly expose error types
- finish components endpoints
- add Components
- add internal utils package, move IntArrayToString() util to new home
- add integration tests for key transactions
- add query param filters for ListKeyTransactions
- add GetKeyTransaction
- add ListKeyTransactions
- add DeleteLabel
- add CreateLabel
- add ListLabels, add GetLabel
- add DeleteDeployment
- add CreateDeployment
- add ListDeployments
- centralize apm test helpers
- add DeleteNrqlAlertCondition
- add UpdateNrqlAlertCondition
- add CreateNrqlAlertCondition
- add GetNrqlAlertCondition
- add ListNrqlAlertConditions
- add UpdateAlertPolicy
- add DeleteAlertCondition
- add CreateAlertCondition
- add GetAlertCondition
- add ListAlertConditions
- get infra condition integration tests passing
- add InfrastructureConditions
- add MonitorScripts
- add MonitorScript
- add DeleteAlertPolicyChannel, update unit tests, add integration test (might need to remove this)
- add alert policy channels
- add synthetics alert conditions
- add synthetics alert conditions
- add GetAlertChannel method
- add CreateAlertChannel, ListAlertChannels, DeleteAlertChannel
- add DeleteMonitor
- add UpdateMonitor
- add CreateMonitor
- add dashboards
- add DeleteAlertPolicy method
- add UpdateAlertPolicy method
- add CreateAlertPolicy method
- add GetAlertPolicy method
- add ListAlertPolicies method
- alerts package
- create remaining CRUD methods for application resource
- add new dependency-free client implementation
- add version.go per auto-versioning docs
- add ListAlertConditions for infrastructure
- add infra namespace
- add catchall newrelic package
- add New Relic environment enum
- maximize page size for ListMonitors
- add ListMonitors method for Synthetics monitors
- add application filtering for ListApplications
- get TestListApplications passing

<a name="v0.1.1"></a>
## [v0.1.1] - 2017-08-02
<a name="v0.1.0"></a>
## v0.1.0 - 2017-06-21
[Unreleased]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.19.1...HEAD
[v2.19.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.19.0...v2.19.1
[v2.19.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.18.0...v2.19.0
[v2.18.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.17.0...v2.18.0
[v2.17.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.16.0...v2.17.0
[v2.16.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.15.1...v2.16.0
[v2.15.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.15.0...v2.15.1
[v2.15.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.14.1...v2.15.0
[v2.14.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.14.0...v2.14.1
[v2.14.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.13.5...v2.14.0
[v2.13.5]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.13.4...v2.13.5
[v2.13.4]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.13.3...v2.13.4
[v2.13.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.13.2...v2.13.3
[v2.13.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.13.1...v2.13.2
[v2.13.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.13.0...v2.13.1
[v2.13.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.12.1...v2.13.0
[v2.12.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.12.0...v2.12.1
[v2.12.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.11.1...v2.12.0
[v2.11.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.11.0...v2.11.1
[v2.11.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.10.3...v2.11.0
[v2.10.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.10.2...v2.10.3
[v2.10.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.9.0...v2.10.2
[v2.9.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.8.0...v2.9.0
[v2.8.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.7.5...v2.8.0
[v2.7.5]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.7.4...v2.7.5
[v2.7.4]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.7.3...v2.7.4
[v2.7.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.7.2...v2.7.3
[v2.7.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.7.1...v2.7.2
[v2.7.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.7.0...v2.7.1
[v2.7.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.6.1...v2.7.0
[v2.6.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.6.0...v2.6.1
[v2.6.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.5.1...v2.6.0
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
[v1.0.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.11.0...v1.0.0
[v0.11.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.10.1...v0.11.0
[v0.10.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.10.0...v0.10.1
[v0.10.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.9.0...v0.10.0
[v0.9.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.8.0...v0.9.0
[v0.8.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.7.1...v0.8.0
[v0.7.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.7.0...v0.7.1
[v0.7.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.6.0...v0.7.0
[v0.6.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.5.1...v0.6.0
[v0.5.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.5.0...v0.5.1
[v0.5.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.1.1...v0.2.0
[v0.1.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v0.1.0...v0.1.1
