<a name="v3.25.0"></a>
## [v3.25.0] - 2023-06-21
### Bug Fixes
- **dashboards:** update to 'text' field schema in widget_markdown ([#2400](https://github.com/newrelic/terraform-provider-newrelic/issues/2400))
- **dashboards:** code changes to enable dashboards with empty pages ([#2397](https://github.com/newrelic/terraform-provider-newrelic/issues/2397))

### Documentation Updates
- **cloud_guides:** beautification of cloud provider docs ([#2399](https://github.com/newrelic/terraform-provider-newrelic/issues/2399))
- **cloud_integrations_guide:** Fix missing a quote at cloud_integrations_guide.html.markdown ([#2401](https://github.com/newrelic/terraform-provider-newrelic/issues/2401))
- **guides:** updates to getting started guide to remove deprecated resources ([#2396](https://github.com/newrelic/terraform-provider-newrelic/issues/2396))

### Features
- **data_source_newrelic_notifications_destination:** Support name in destination data source ([#2395](https://github.com/newrelic/terraform-provider-newrelic/issues/2395))

<a name="v3.24.2"></a>
## [v3.24.2] - 2023-06-12
<a name="v3.24.1"></a>
## [v3.24.1] - 2023-06-08
### Bug Fixes
- **cloud_tests:** adds precheck to cloud tests to delink accounts before test run ([#2384](https://github.com/newrelic/terraform-provider-newrelic/issues/2384))
- **log_parsing:** fix to error handling in the log parsing resource ([#2390](https://github.com/newrelic/terraform-provider-newrelic/issues/2390))
- **newrelic_entity:** filtering entites returned based on account_id ([#2389](https://github.com/newrelic/terraform-provider-newrelic/issues/2389))

### Documentation Updates
- account_id param for newrelic_one_dashboard_json ([#2385](https://github.com/newrelic/terraform-provider-newrelic/issues/2385))
- **newrelic:** tiny typo corrections ([#2392](https://github.com/newrelic/terraform-provider-newrelic/issues/2392))

<a name="v3.24.0"></a>
## [v3.24.0] - 2023-05-30
### Bug Fixes
- **build:** Updated gofmt and gofmt-fix make rules

### Documentation Updates
- **script_monitor:** corrections to details on private, public locations ([#2372](https://github.com/newrelic/terraform-provider-newrelic/issues/2372))

### Features
- **newrelic_cloud_aws_integrations:** Add remaining aws polling integrations. ([#2383](https://github.com/newrelic/terraform-provider-newrelic/issues/2383))
- **newrelic_cloud_aws_integrations:** Add other polling integrations. ([#2374](https://github.com/newrelic/terraform-provider-newrelic/issues/2374))

<a name="v3.23.0"></a>
## [v3.23.0] - 2023-05-17
### Bug Fixes
- update missing fields, docs etc
- Missing fields from [#2289](https://github.com/newrelic/terraform-provider-newrelic/issues/2289)
- lint errors

### Documentation Updates
- update aws docs
- **cloud_integrations:** update to broken links ([#2365](https://github.com/newrelic/terraform-provider-newrelic/issues/2365))
- **examples:** updated aws example to include config stream and moved cloud examples to modules ([#2225](https://github.com/newrelic/terraform-provider-newrelic/issues/2225))
- **newrelic_cloud_aws_integrations:** update aws docs with new integrations.
- **newrelic_cloud_aws_integrations:** update aws docs with new integrations.
- **one_dashboard_json:** addition of a multipage dashboard example and other corrections ([#2347](https://github.com/newrelic/terraform-provider-newrelic/issues/2347))

### Features
- Add other polling integrations to newrelic_cloud_aws_integrations
- additional poll integrations (sqs, ebs, alb, elasticache) ([#2289](https://github.com/newrelic/terraform-provider-newrelic/issues/2289))

<a name="v3.22.0"></a>
## [v3.22.0] - 2023-05-08
### Documentation Updates
- **one_dashboard:** added documentation on how to use the CLI to convert dashboards to HCL ([#2339](https://github.com/newrelic/terraform-provider-newrelic/issues/2339))
- **service_level_alert_helper:** update to documentation ([#2361](https://github.com/newrelic/terraform-provider-newrelic/issues/2361))

### Features
- **data/account:** Allow reading the provider's `account_id`. ([#2314](https://github.com/newrelic/terraform-provider-newrelic/issues/2314))
- **newrelic_cloud_azure_integrations:** addition of the integration 'azureMonitor' ([#2338](https://github.com/newrelic/terraform-provider-newrelic/issues/2338))

<a name="v3.21.3"></a>
## [v3.21.3] - 2023-04-28
<a name="v3.21.2"></a>
## [v3.21.2] - 2023-04-28
<a name="v3.21.1"></a>
## [v3.21.1] - 2023-04-28
<a name="v3.21.0"></a>
## [v3.21.0] - 2023-04-28
### Bug Fixes
- **newrelic_synthetics_script_monitor:** Populate script argument changes upon modifications. ([#2335](https://github.com/newrelic/terraform-provider-newrelic/issues/2335))
- **nrql_alert_condition:** check for nil response on create
- **one_dashboard:** code changes to support y_axis_left_min=0 ([#2326](https://github.com/newrelic/terraform-provider-newrelic/issues/2326))
- **synthetics:** handle omitting runtime values for legacy runtime

### Documentation Updates
- update synthetic monitor default runtime notes. ([#2325](https://github.com/newrelic/terraform-provider-newrelic/issues/2325))
- **nrql_alert_condition:** Update nrql alert condition docs ([#2341](https://github.com/newrelic/terraform-provider-newrelic/issues/2341))

### Features
- **one_dashboard:** addition of the attribute zero to widget_line ([#2336](https://github.com/newrelic/terraform-provider-newrelic/issues/2336))
- **synthetics:** Synthetics monitor additional field to output period in minutes ([#2340](https://github.com/newrelic/terraform-provider-newrelic/issues/2340))
- **synthetics:** Synthetics monitor additional field to output period in minutes
- **synthetics:** add device emulation options to newrelic_synthetics_script_browser_monitor resource

<a name="v3.20.2"></a>
## [v3.20.2] - 2023-04-07
### Bug Fixes
- **newrelic_nrql_alert_condition:** add default 'violation_time_limit_seconds' ([#2319](https://github.com/newrelic/terraform-provider-newrelic/issues/2319))

### Documentation Updates
- add note on 'default_values' to one_dashboard ([#2321](https://github.com/newrelic/terraform-provider-newrelic/issues/2321))

<a name="v3.20.1"></a>
## [v3.20.1] - 2023-04-05
### Bug Fixes
- **newrelic_api_access_key:** add info to TF docs on API Access errors ([#2302](https://github.com/newrelic/terraform-provider-newrelic/issues/2302))

### Documentation Updates
- correct references to slack webhooks in newrelic_alert_channel ([#2315](https://github.com/newrelic/terraform-provider-newrelic/issues/2315))
- **guide:** update synthetics migration guide with correct resource reference to remove state ([#2318](https://github.com/newrelic/terraform-provider-newrelic/issues/2318))

<a name="v3.20.0"></a>
## [v3.20.0] - 2023-03-30
### Bug Fixes
- **newrelic_entity:** add helpers to escape single quotes in NRQL que… ([#2295](https://github.com/newrelic/terraform-provider-newrelic/issues/2295))

### Documentation Updates
- update newrelic_synthetics_cert_check_monitor example to use correct argument. ([#2311](https://github.com/newrelic/terraform-provider-newrelic/issues/2311))
- update docs with newrelic_account_management resource ([#2308](https://github.com/newrelic/terraform-provider-newrelic/issues/2308))

### Features
- **service_level:** add data source service_level_alert_helper ([#2298](https://github.com/newrelic/terraform-provider-newrelic/issues/2298))

<a name="v3.19.0"></a>
## [v3.19.0] - 2023-03-28
### Documentation Updates
- add examples to dashboard resource to demonstrate importing dashboards using GUID ([#2300](https://github.com/newrelic/terraform-provider-newrelic/issues/2300))

### Features
- **Alerts:** Expose entity_guid for legacy alert conditions ([#2301](https://github.com/newrelic/terraform-provider-newrelic/issues/2301))

<a name="v3.18.1"></a>
## [v3.18.1] - 2023-03-24
### Bug Fixes
- **NRQLDropRule:** Handle case where resource in state has been deleted

### Documentation Updates
- **entity:** update entity data source docs with additional entity types and examples ([#2297](https://github.com/newrelic/terraform-provider-newrelic/issues/2297))

<a name="v3.18.0"></a>
## [v3.18.0] - 2023-03-22
### Features
- Allow creating SL using CDF functions. ([#2293](https://github.com/newrelic/terraform-provider-newrelic/issues/2293))

<a name="v3.17.1"></a>
## [v3.17.1] - 2023-03-17
### Bug Fixes
- **account_management:** fixed the docs page link issue ([#2294](https://github.com/newrelic/terraform-provider-newrelic/issues/2294))

<a name="v3.17.0"></a>
## [v3.17.0] - 2023-03-16
### Bug Fixes
- **one_dashboard:** added raw configuration properties to one dashboard resource ([#2278](https://github.com/newrelic/terraform-provider-newrelic/issues/2278))

### Documentation Updates
- add additional example to one_dashboard_json to demonstrate setting thresholds ([#2292](https://github.com/newrelic/terraform-provider-newrelic/issues/2292))
- add/update resources and data sources in the NR terraform registry ([#2287](https://github.com/newrelic/terraform-provider-newrelic/issues/2287))

### Features
- added account managent resource, docs, examples and test cases ([#2275](https://github.com/newrelic/terraform-provider-newrelic/issues/2275))

<a name="v3.16.1"></a>
## [v3.16.1] - 2023-03-15
### Bug Fixes
- **cloud:** return err when error is present or payload is nil

<a name="v3.16.0"></a>
## [v3.16.0] - 2023-03-09
### Bug Fixes
- **cloud:** Reset id if resource not found

### Documentation Updates
- add deprecation messages to legacy alert resources and data sources
- **alerts:** adjustments to deprecation messages
- **alerts:** remove deprecated evaluation_offset from examples

### Features
- **agentapplication:** add browser agent application resource ([#2262](https://github.com/newrelic/terraform-provider-newrelic/issues/2262))
- **data_partition:** added data partition rules and tests

<a name="v3.15.0"></a>
## [v3.15.0] - 2023-02-27
### Bug Fixes
- **data_notification_destination:** fix naming
- **data_notification_destination:** add missing test - small fix
- **data_notification_destination:** add missing test
- **data_notification_destination:** fix test
- **data_notification_destination:** fix data source wip
- **data_notification_destination:** fix data source + try to fix tests
- **data_notification_destination:** fix data source - wip
- **one_dashboard:** fix for updating with filter_current_dashboard

### Documentation Updates
- moved list of valid values to correct argument ([#2251](https://github.com/newrelic/terraform-provider-newrelic/issues/2251))
- updated notice for v2 support and v3 details
- rename files to match resource names
- **index:** Note on default region

### Features
- **workflows:** expose workflow entity guid

<a name="v3.14.0"></a>
## [v3.14.0] - 2023-02-08
### Bug Fixes
- Muting rule inter condition operator ambiguous ([#2219](https://github.com/newrelic/terraform-provider-newrelic/issues/2219))
- **entities:** validations for domain and type
- **entity:** removed vaildations of domain and types

### Documentation Updates
- Removing value_function from the markdowns
- **alert:** Updated NRQL condition term duration max
- **entity:** Added additional examples
- **entity:** Added new entity type and domain

### Features
- Add evaluation delay to nrql alert condition
- **NotificationDestinationDataSpurce:** destination data source tests + lint
- **NotificationDestinationDataSpurce:** add notification destination data source

<a name="v3.13.0"></a>
## [v3.13.0] - 2023-01-19
### Bug Fixes
- **Workflow:** Linting issues
- **newrelic_one_dashboard:** fix for null pointer when using thresholds in widget_billboard

### Documentation Updates
- **Workflow:** Updated docs
- **dashboards:** improved deprecation message to be less scary

### Features
- **Workflows:** Merge conflict resolved

<a name="v3.12.0"></a>
## [v3.12.0] - 2023-01-10
### Bug Fixes
- **one_dashboard_json:** added retry mechanism to create/update to handle updateAt changes
- **private_locations:** returning the right err ([#2177](https://github.com/newrelic/terraform-provider-newrelic/issues/2177))

### Documentation Updates
- **log_parsing:** added docs for log parsing rule
- **log_parsing:** added docs for log parsing rule
- **log_parsing:** added docs for log parsing rule
- **log_parsing:** added docs for log parsing rule
- **log_parsing:** added docs for log parsing rule
- **log_parsing:** added docs for log parsing rule
- **notification_channel:** fix email notification channel example
- **synthetics:** remove unnecessary quotes around private location guid
- **test_grok:** Document for testgrok data source

### Features
- **log_parsing_rule:** added data source, resource, and tests - NR-53373
- **log_parsing_rule:** Added provider and tests for log parsing rule
- **log_parsing_rule:** Added provider and tests for log parsing rule
- **log_parsing_rule:** Added provider and tests for log parsing rule
- **log_parsing_rule:** Added provider and tests for log parsing rule
- **log_parsing_rule:** Added provider and tests for log parsing rule
- **provider:** set user agent service name via -ldflags
- **test_grok:**  added data source provider

<a name="v3.11.0"></a>
## [v3.11.0] - 2022-12-16
### Bug Fixes
- **workflows:** fix a bug that would prevent creation of disabled workflows

### Features
- Remove value_function from nrql alert condition resource

<a name="v3.10.0"></a>
## [v3.10.0] - 2022-12-16
### Bug Fixes
- **dashboards:** handle empty values returned from API
- **drop_rules:** verbose error message :bug:
- **newrelic_nrql_alert_condition:** Set entity_guid after creating NRQL alert condition
- **workflows:** stop silently removing channels on workflow updates/deletes
- **workloads:** Removed forcenew for entity guids and changed docs

### Documentation Updates
- **Workflow:** Changes in examples

### Features
- **secure_credential:** add account id in secure credential search

<a name="v3.9.0"></a>
## [v3.9.0] - 2022-12-06
### Bug Fixes
- linting2
- linting
- **notifications errors:** Added error details to response
- **synthetics:** unset private location ID and return nil if entity not found

### Documentation Updates
- **dashboard:** add documentation for variables
- **entity_tags:** add example of using a dynamic block to apply multiple tags to an entity

### Features
- **cloud:** add azure mysql flexible, postgres flexible and gcp alloydb integrations
- **one_dashboard:** add variables

<a name="v3.8.0"></a>
## [v3.8.0] - 2022-11-30
### Bug Fixes
- **destination:** unset destination ID and return nil if destination not found (prevents crash)
- **notifications:** added deleted test
- **notifications:** fix auth_basic bug
- **synthetics:** set verify_ssl and validation_string during read/import operations
- **synthetics:** set period and status during read/import operations
- **synthetics:** set period and status on update of newrelic_synthetics_script_monitor
- **synthetics:** set additional attribute values on import
- **workflow:** unset workflow ID and return nil if workflow not found (prevents crash)

### Documentation Updates
- Update 'violations' to 'incidents' in Alerts docs
- fix import command for script monitor
- **notifications:** Updated notifications channel / destination docs to be clearer about properties
- **obfusation_rule:** added rules docs and fixed expression docs
- **obfuscation_expression:** Added the data source document for the expression
- **obfuscation_rule:** minor changes
- **synthetics_cert_check_monitor:** fix url instead of domain in domain field

### Features
- **data_source_entity:** remove limit on tags
- **obfuscation_expression:** Added Obfuscation expression data source
- **obfuscation_rule:** Added Obfuscation rule and tests

<a name="v3.7.1"></a>
## [v3.7.1] - 2022-11-15
### Documentation Updates
- **obfuscation_expression:** change name and regex to required

<a name="v3.7.0"></a>
## [v3.7.0] - 2022-11-14
### Bug Fixes
- **workflows:** make it possible to remove all workflow enrichments

### Documentation Updates
- **obfuscation_expression:** minor changes
- **obfuscation_expression:** added docs for obfuscation expression

### Features
- add optional account_id parameter to private location data source
- **obfuscation_expression:** Added Obfuscation expression and tests
- **synthetics:** set defaults for runtime to new runtime and added instructions for old runtime

<a name="v3.6.1"></a>
## [v3.6.1] - 2022-10-27
### Bug Fixes
- dashboard import statement missing guid

### Documentation Updates
- improve the documentation for the muting rules handling parameter

<a name="v3.6.0"></a>
## [v3.6.0] - 2022-10-26
### Documentation Updates
- **workloads:** added note for rule
- **workloads:** Updated docs on status config

### Features
- add workload status_config

<a name="v3.5.2"></a>
## [v3.5.2] - 2022-10-21
### Bug Fixes
- **workflows:** fix compilation issues after client version update
- **workloads:** fix compilation issues after client version update

### Documentation Updates
- **alert_conditions:** added the docs for mobile_metric missing types

<a name="v3.5.1"></a>
## [v3.5.1] - 2022-10-19
### Bug Fixes
- **workflows:** correctly handle notification channels deleted outside TF

### Documentation Updates
- Update reference for data.newrelic_entity

<a name="v3.5.0"></a>
## [v3.5.0] - 2022-10-14
### Bug Fixes
- add check of default answer of the API for service level select
- correct attribute mapping in service level event query select
- **notifications:** fix monitoring property to destinations & channels
- **notifications:** add monitoring property to destinations & channels
- **workflows:** stop forcing workflow recreation when using the default account id

### Features
- bump of new-relic-client-go and add service level capability to use select field for event queries

<a name="v3.4.4"></a>
## [v3.4.4] - 2022-10-06
### Bug Fixes
- add account_id to secure credential data source schema

<a name="v3.4.3"></a>
## [v3.4.3] - 2022-10-05
### Bug Fixes
- up timeout to give entity chance to index

<a name="v3.4.2"></a>
## [v3.4.2] - 2022-10-04
### Documentation Updates
- change custom_headers to custom_header

<a name="v3.4.1"></a>
## [v3.4.1] - 2022-10-03
### Documentation Updates
- update monitor with private location examples

<a name="v3.4.0"></a>
## [v3.4.0] - 2022-10-03
### Bug Fixes
- stop setting enrichments to an empty array on state update
- stop fuzzy search on resource_newrelic_application_settings failing apply

### Documentation Updates
- fix issues with synthetics monitors docs
- Update migration_guide_v3.html.markdown
- **synthetics_monitor:** mentions NerdGraph is used.
- **synthetics_monitor:** mentions NerdGraph is used.

### Features
- **one_dashboard_json:** added newrelic_one_dashboard_json resource

<a name="v3.3.0"></a>
## [v3.3.0] - 2022-09-22
### Bug Fixes
- **docs:** Add guid to import command
- **docs:** Add safe migration pathway for monitors

### Features
- **one_dashboard:** added linked_entity_guids and filter_current_dashboard to heatmap

<a name="v3.2.1"></a>
## [v3.2.1] - 2022-09-15
### Bug Fixes
- **channels:** add check for read channels
- **channels:** lint
- **docs:** link to Azure cloud integration example
- **notifications:** add examples to docs
- **notifications:** made account_id optional [#1994](https://github.com/newrelic/terraform-provider-newrelic/issues/1994)
- **notifications:** add missing destination types
- **simple_monitor:** [#1984](https://github.com/newrelic/terraform-provider-newrelic/issues/1984) populating uri and custom headers on input struct
- **workflows:** add schema upgrade version
- **workflows:** add schema upgrade version - wip
- **workflows:** add schema upgrade version - wip
- **workflows:** add schema upgrade version
- **workflows:** docs small change
- **workflows:** change workflows schema

### Documentation Updates
- **newrelic_one_dashboard:** Fix incorrect permission value)

<a name="v3.2.0"></a>
## [v3.2.0] - 2022-09-08
### Bug Fixes
- made accountId optional in the schema
- **alert_policy:** added computed to account_id

### Documentation Updates
- added note in the synthetics private locations data source docs
- added note in the synthetics private locations data source docs

### Features
- add s3 and docDb to AWS integration

<a name="v3.1.0"></a>
## [v3.1.0] - 2022-09-02
### Bug Fixes
- **src:** Delete exactlyOne validation in scheme of Notification Destination

### Features
- **src:** Add ConflictsWith schema validation - Add this validation within the auth_token and auth_basic to validate uniqueness of each other

<a name="v3.0.4"></a>
## [v3.0.4] - 2022-09-01
### Bug Fixes
- **data_alert_channel:** added missing account_id
- **data_source_newrelic_alert_channel:** fixed support for multi accounts
- **newrelic_alert_policy:** added multi account support

<a name="v3.0.3"></a>
## [v3.0.3] - 2022-09-01
### Bug Fixes
- **newrelic_entity_tags:** increased timeout level

<a name="v3.0.2"></a>
## [v3.0.2] - 2022-08-31
### Bug Fixes
- **newrelic_synthetics_monitor:** convert old ID's into GUID
- **newrelic_synthetics_monitor:** deleted synthetics monitors are not detected
- **newrelic_synthetics_monitor:** handle already deleted checks
- **synthetics:** catch when synthetic check have been deleted from ui

<a name="v3.0.1"></a>
## [v3.0.1] - 2022-08-29
<a name="v3.0.0"></a>
## [v3.0.0] - 2022-08-26
### Bug Fixes
- **Notifications:** add update to destinations and channels + rename properties
- **Notifications:** pr review fixes: not needed remove computes and forcenew
- **channels:** renaming - wip
- **channels:** add tests
- **channels:** small renaming
- **channels:** fix lint
- **destinations:** general adjustments for destinations read functionality
- **destinations:** upgrade client version
- **destinations:** update docs
- **docs:** add notifications resources to index doc
- **newrelic_entity_tags:** extended timeout to see if it fixes not found errors
- **notification_channel:** ignore channel not found error during terraform destroy (deleting workflow deletes channel)
- **notifications:** add custom errors handling to notifications
- **notifications:** fix tests -wip
- **notifications:** fix tests
- **notifications:** fix tests + add jira
- **notifications:** cr changes
- **notifications:** add tests + update docs
- **notifications:** update docs
- **notifications:** add internal property
- **notifications:** upgrade newrelic-client-go
- **notifications:** lint fix
- **nrql_alert_condition:** do not show diff for streaming methods when not provided and defaults are used
- **nrql_alert_condition:** update docs and validation for 'expiration_duration'
- **workflows:** add custom errors handling to workflows
- **workflows:** small workflows test twick
- **workflows:** upgrade newrelic-client-go
- **workflows:** lint fix

### Documentation Updates
- add v3 migration guide
- boeken links monitor docs
- documentation for cert_broken_step monitors
- added the changes in website guide
- adding the website guides
- doc of simple & browser synthectics monitor
- **synthetics:** update docs for private location data source

### Features
- swap secure credential resource to GraphQL API
- add new synthetics resources
- **newrelic_synthetics_cert_check_monitor:** adding cert check monitor
- **synthetics:** added private location resource
- **synthetics:** add step monitor resource
- **synthetics:** [wip] add broken links monitor resource
- **synthetics:** add newrelic_synthetics_script_monitor resource
- **synthetics:** migrate monitor location data source to new GraphQL API
- **workflows:** add workflows resources - wip
- **workflows:** add workflows resources
- **workflows:** add workflows resources workinggg - wip
- **workflows:** add tests
- **workflows:** resolve conflict

### BREAKING CHANGE

new synthetics resources use GraphQL API schema

<a name="2.51.0"></a>
## [2.51.0] - 2022-07-22
<a name="v2.50.2"></a>
## [v2.50.2] - 2022-08-30
<a name="v2.50.1"></a>
## [v2.50.1] - 2022-08-30
<a name="v2.50.0"></a>
## [v2.50.0] - 2022-08-30
### Features
- updated dependencies

<a name="v2.49.1"></a>
## [v2.49.1] - 2022-07-22
### Bug Fixes
- **Notifications:** relocate notifications docs to the corrects folder

<a name="v2.49.0"></a>
## [v2.49.0] - 2022-07-19
### Bug Fixes
- **channels:** lint fix
- **channels:** add import test
- **channels:** add docs and fix small bug
- **destinations:** lint fix
- **destinations:** add note on doc and import test
- **destinations:** fix pr review suggetion
- **destinations:** add docs and fix small bug
- **destinations:** fix types according to the new go client api - still WIP
- **destinations:** fix tests
- **destinations:** fix tests

### Features
- **channels:** upgrade go-client version
- **channels:** add tests - wip
- **channels:** add tests - wip
- **channels:** add notifications channels provider
- **destinations:** upgrade go-client version
- **destinations:** upgrade go-client version
- **destinations:** fix tests
- **notifications:** add tests
- **notifications:** working notification destination
- **notifications:** add notifications destinations

<a name="v2.48.2"></a>
## [v2.48.2] - 2022-07-15
<a name="v2.48.1"></a>
## [v2.48.1] - 2022-06-30
### Bug Fixes
- force service level when changing account id

<a name="v2.48.0"></a>
## [v2.48.0] - 2022-06-21
### Bug Fixes
- add sleep to wait the SL to be indexed and avoid flaky test

### Features
- **docs:** add azure integrations guide and examples

<a name="v2.47.1"></a>
## [v2.47.1] - 2022-06-15
### Bug Fixes
- remove unnecessary read of service level entity after creation

<a name="v2.47.0"></a>
## [v2.47.0] - 2022-06-10
### Features
- **newrelic_one_dashboard:** added support for ignore_time_range

<a name="v2.46.2"></a>
## [v2.46.2] - 2022-06-06
<a name="v2.46.1"></a>
## [v2.46.1] - 2022-05-25
<a name="v2.46.0"></a>
## [v2.46.0] - 2022-05-23
### Features
- **docs:** add azure integrations guide and examples

<a name="v2.45.1"></a>
## [v2.45.1] - 2022-05-14
### Documentation Updates
- **gcp-example:** added GCP example to docs
- **provider:** additional information regarding version constraints
- **provider:** add info regarding setting/upgrading the provider version

<a name="v2.45.0"></a>
## [v2.45.0] - 2022-05-12
### Bug Fixes
- **alerts:** Allow negative threshold values for non-baseline NRQL conditions

### Documentation Updates
- **alerts:** Document new NRQL term threshold operators
- **alerts:** Remove minimum threshold requirement for non-baseline NRQL conditions

### Features
- added gcp integrations test tf script
- **alerts:** Adds 3 term threshold operators for NRQL conditions

<a name="v2.44.0"></a>
## [v2.44.0] - 2022-05-03
### Bug Fixes
- Description of entity_guid for NRQL conditions.

### Features
- Update NRQL Condition docs with tag management example
- Update docs with entity_guid attribute reference description
- Expose NRQL Condition entityGUID on conditions

<a name="v2.43.4"></a>
## [v2.43.4] - 2022-04-25
### Bug Fixes
- **cloud_azure:** correct firewalls integration on update

### Documentation Updates
- updated the compiler
- update terraform versions
- **cloud_azure:** remove duplicate event_hub in example

<a name="v2.43.3"></a>
## [v2.43.3] - 2022-04-22
### Bug Fixes
- **one_dashboard:** make limit for bullet widget required

<a name="v2.43.2"></a>
## [v2.43.2] - 2022-04-21
### Bug Fixes
- Add payload_string to channel data source

<a name="v2.43.1"></a>
## [v2.43.1] - 2022-04-20
### Documentation Updates
- **resource/one_dashboard:** Add filter_current_dashboard doc for widget_bar and widget_pie

<a name="v2.43.0"></a>
## [v2.43.0] - 2022-04-19
### Bug Fixes
- Change logic for reading violation time limits to fix imports

### Documentation Updates
- added gcp integrations documentation

### Features
- added test.go for gcp integrations
- added gcp cloud integrations to resource group map
- **cloud:** add azure integrations resource

<a name="v2.42.1"></a>
## [v2.42.1] - 2022-04-14
### Bug Fixes
- Tags resource implicit dependency
- Import payload_string when importing webhook channel
- remove 14 days deprecated option from service levels
- Update DiffSuppressFunc for aggregation delay/timer

<a name="v2.42.0"></a>
## [v2.42.0] - 2022-04-07
### Documentation Updates
- add aws integrations documentation
- **cloud-integrations:** added guide and example for AWS

### Features
- add aws integrations resource

<a name="v2.41.4"></a>
## [v2.41.4] - 2022-04-06
### Bug Fixes
- force to create new alert condition if type changes

### Documentation Updates
- Remove signal block from deprecation message

<a name="v2.41.3"></a>
## [v2.41.3] - 2022-04-05
### Documentation Updates
- remove beta level and add an example with tags for service levels

<a name="v2.41.2"></a>
## [v2.41.2] - 2022-03-25
### Bug Fixes
- add len check when creating cloud link accounts
- Allow 0 values for aggregation_delay
- Allow 0 values for aggregation_delay
- **cloud:** add import
- **cloud_azure_link_account:** rename client_secret_id to client_secret
- **docs:** removed two accountID warnings that were incorrect
- **docs:** removed two accountID warnings that were incorrect

### Documentation Updates
- **cloud_azure_link_account:** change client_secret_id to client_secret

<a name="v2.41.1"></a>
## [v2.41.1] - 2022-03-21
### Bug Fixes
- add len check when creating cloud link accounts

<a name="v2.41.0-beta.2"></a>
## [v2.41.0-beta.2] - 2022-03-17
### Bug Fixes
- handle setting linked_entity_guids
- **newrelic_drop_rule:** fix for [#1611](https://github.com/newrelic/terraform-provider-newrelic/issues/1611) added extra check on API return
- **newrelic_one_dashboard:** fixed 0 not getting pushed to API
- **nrql_alert_condition:** Crash with deprecated since_value and evaluation_offset fields
- **nrql_alert_condition:** Move condition validation to API
- **nrql_alert_condition:** remove computed flag for slide_by field
- **rql_alert_condition:** Add nil check when flattening slide_by
- **service_levels:** Removing an option of 14 days for trailing windows

### Documentation Updates
- minor changes
- added note to specify the resource and data source that uses rest api's
- remove beta info ahead of general release
- fix broken link to terraform docs
- document allowed values for expiration_duration
- **nrql_alert_condition:** removed legacy urls from doc
- **nrql_alert_condition:** removed legacy urls from doc
- **nrql_alert_condition:** removed legacy urls from doc
- **servicelevel:** SLO periods now include complete weeks

<a name="v2.41.0-beta.1"></a>
## [v2.41.0-beta.1] - 2022-03-15
### Bug Fixes
- handle filter_current_dashboard on update

### Documentation Updates
- minor changes in docs

### Features
- awsGov cloud integration
- awsGov cloud integration

<a name="v2.40.0"></a>
## [v2.40.0] - 2022-03-14
### Bug Fixes
- allow manipulation of SLIs with non existing related entity, return sli_guid, change service level example
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc
- fixed mistakes in the Doc

### Documentation Updates
- added gcp integration documentation
- added gcp integration documentation
- added gcp integration documentation
- added gcp integration documentation
- added gcp integration documentation
- minor changes
- added gcp integration documentation
- added gcp integration documentation
- minor changes
- added gcp integration documentation
- added gcp integration documentation
- added gcp integration documentation
- add cloud_aws_link_account resource to docs
- update version in docs
- **servicelevel:** SLO periods now include complete weeks

### Features
- azure integration
- added gcp cloud link account
- added gcp cloud link account
- azure integration
- added gcp cloud link account
- added gcp cloud link account
- added gcp cloud link account
- added gcp cloud link account
- **aws_link_account:** add AWS Link Account resource
- **provider:** added gcp resource to resource map

<a name="v2.39.2"></a>
## [v2.39.2] - 2022-03-08
### Features
- **docs:** added accountID to alert_channel and alert_policy_channel

<a name="v2.39.1"></a>
## [v2.39.1] - 2022-03-07
### Bug Fixes
- **newrelic_alert_policy_channel:** add account_id to terraform schema
- **newrelic_alert_policy_channel:** add XAccountID to context when creating or updating alert channels

<a name="v2.39.0"></a>
## [v2.39.0] - 2022-03-04
### Bug Fixes
- **dashboard:** handle filter_current_dashboard on read

### Features
- **Alerts:** deprecate nrql condition value_function, sum, single_value
- **servicelevel:** SLI should contain one and only one objective

<a name="v2.38.0"></a>
## [v2.38.0] - 2022-02-18
### Documentation Updates
- Fix broken V2 migration links
- add cloud_account data source docs
- fix broken link to data sources docs

### Features
- add cloud account data source

<a name="v2.37.0"></a>
## [v2.37.0] - 2022-02-13
### Bug Fixes
- **newrelic_drop_rule:** fix for [#1611](https://github.com/newrelic/terraform-provider-newrelic/issues/1611) added extra check on API return
- **nrql_alert_condition:** Move condition validation to API
- **nrql_alert_condition:** Crash with deprecated since_value and evaluation_offset fields

### Documentation Updates
- remove beta info ahead of general release
- fix broken link to terraform docs
- document allowed values for expiration_duration
- **servicelevel:** SLO periods now include complete weeks

<a name="v2.36.2"></a>
## [v2.36.2] - 2022-02-03
### Bug Fixes
- **nrql_alert_condition:** remove computed flag for slide_by field

<a name="v2.36.1"></a>
## [v2.36.1] - 2022-02-02
### Bug Fixes
- **rql_alert_condition:** Add nil check when flattening slide_by
- **service_levels:** Removing an option of 14 days for trailing windows

<a name="v2.36.0"></a>
## [v2.36.0] - 2022-02-01
### Documentation Updates
- update versions url
- fix URL to debugging terraform
- **servicelevel:** SLO periods now include complete weeks

### Features
- **nrql_alert_condition:** Add slide by support for alert conditions

<a name="v2.35.1"></a>
## [v2.35.1] - 2022-01-24
### Bug Fixes
- Remove condition name length validation
- **renovate:** removed old feature branch check

### Documentation Updates
- update provider version documentation link
- add example of workload using tags

<a name="v2.35.0"></a>
## [v2.35.0] - 2022-01-10
### Bug Fixes
- correct the HMAC calculation for synthetics
- revert previous changes, add DiffSuppressFunc for default values
- aggregation_method and aggregation_delay diff when not provided
- muting rule with repeat = null crashes plugin
- correct violation_time_limit if none is provided
- use violation_time_limit_seconds on condition import
- **servicelevel:** Force new resource if GUID changes
- **synthetics_multilocation_alert:** fix values for violation_time_limit_seconds

### Documentation Updates
- add vse_password documentation
- Fix link to Install Terraform

<a name="v2.34.1"></a>
## [v2.34.1] - 2021-12-10
### Bug Fixes
- manually handled the state change for filter_current_dashboard. Issue 1494
- **newrelic_one_dashboard:** cannot remove billboard threshold

### Features
- **newrelic_one_dashboard:** return nil value if critical/warning is not set for billboard

<a name="v2.34.0"></a>
## [v2.34.0] - 2021-12-07
### Features
- **monitor_script:** add vse_password for private monitor script locations

<a name="v2.33.0"></a>
## [v2.33.0] - 2021-12-01
### Documentation Updates
- **servicelevel:** Update Service Level docs

### Features
- **nrql_drop_rule:** Add 'drop_attributes_from_metric_aggregates' to drop rule actions.

<a name="v2.32.0"></a>
## [v2.32.0] - 2021-11-16
### Bug Fixes
- **one_dashboard:** surface error messages on dashboard update

<a name="v2.31.1"></a>
## [v2.31.1] - 2021-11-08
### Documentation Updates
- Fix synthetics_multilocation_alert_condition.markdown
- **servicelevel:** Fix example in the docs

### Features
- **alert_muting_rule:** Accept entity.guid and tags.NAME attrs

<a name="v2.30.2"></a>
## [v2.30.2] - 2021-10-26
### Bug Fixes
- update docs and validation for muting rule condition operator
- remove AtLeastOneOf to make violation_time_limit_seconds optional
- **resource_newrelic_entity_tags:** immutable tags are no longer returned

### Features
- improve error messaging for alert_muting_rule

<a name="v2.30.1"></a>
## [v2.30.1] - 2021-10-22
### Bug Fixes
- resolved an issue where multiple pages would link to the wrong page

### Features
- Added documentation for filter_current_dashboard

<a name="v2.30.0"></a>
## [v2.30.0] - 2021-10-14
### Documentation Updates
- update docs to explain baseline thresholds
- fixing grammar error / missing word

### Features
- added another check that filter_current_dashboard is set before collecting widget details
- Added validation for if linked_entity_guids is set. Reworked complexity of finding and setting linked page entity
- added filter_current_dashboard support and added a test to verify functionality

<a name="v2.29.0"></a>
## [v2.29.0] - 2021-10-12
### Documentation Updates
- change timeWindow to time_window

### Features
- **dashboard:** added widget_stacked_bar

<a name="v2.28.0"></a>
## [v2.28.0] - 2021-10-08
### Bug Fixes
- update expected error text
- **docs:** added new replacement pattern

### Documentation Updates
- Modifying referenced hyperlinks
- **Alerts:** Clean-up and clarify NRQL alert condition documentation
- **linked_entity_guids:** Documentation and some tests
- **servicelevel:** Add Service Level documentation

### Features
- **Alerts:** Enhance alerts nrql condition errors
- **linked_entity_guids:** Expose linked_entity_guids for dashboard_raw

<a name="v2.27.1"></a>
## [v2.27.1] - 2021-10-06
<a name="v2.27.0"></a>
## [v2.27.0] - 2021-10-05
### Features
- **Alerts:** Add streaming methods fields to nrql alert conditions

<a name="v2.26.0"></a>
## [v2.26.0] - 2021-10-01
### Features
- **servicelevel:** Add service level resource

<a name="v2.25.0"></a>
## [v2.25.0] - 2021-08-04
<a name="v2.24.1"></a>
## [v2.24.1] - 2021-07-21
### Bug Fixes
- **build:** fix compile-only and compile-all build targets to aid with local development
- **docs:** added baseline to list of options that require operator to be set to above

### Documentation Updates
- **InfraAlerts:** Update docs
- **MonitorScript:** Update docs with monitor script location
- **OneDashboardRaw:** Update docs

### Features
- removed binary file
- fix tests
- 0 violation TTL for Infra Conditions returns warning
- use client method context from within resources
- **MonitorScript:** Add monitor script locations
- **dashboard_raw:** add newrelic_one_dashboard_raw
- **docs:** added documentation links to the dashboard migration guide

<a name="v2.23.0"></a>
## [v2.23.0] - 2021-06-10
### Bug Fixes
- **AlertCondition:** Fix TestAccNewRelicAlertCondition_LongName unit test
- **newrelic_infra_alert_condition:** Prevent index out of range on expandIfraAlertThreshold ([#606](https://github.com/newrelic/terraform-provider-newrelic/issues/606))
- **plugins:** Run gofmt
- **resource_newrelic_alert_condition:** update unit test to match resource
- **resource_newrelic_alert_condition:** adapted name length to fit API definition
- **resource_newrelic_alert_condition:** incorrect description field

### Documentation Updates
- **plugins:** Update website with deprecation notices

<a name="v2.22.1"></a>
## [v2.22.1] - 2021-05-12
### Bug Fixes
- **docs:** bumped version used in docs to latest

### Features
- **docs:** added guide for newrelic_dashboard migration

<a name="v2.22.0"></a>
## [v2.22.0] - 2021-05-10
### Bug Fixes
- **infra_alert_condition:** Added missing documentation, fixes [#1280](https://github.com/newrelic/terraform-provider-newrelic/issues/1280)

### Features
- **resource_newrelic_dashboard:** Added deprecation notice

<a name="v2.21.2"></a>
## [v2.21.2] - 2021-05-06
### Features
- **newrelic_one_dashboard:** add JSON widgets

<a name="v2.21.1"></a>
## [v2.21.1] - 2021-04-16
<a name="v2.21.0"></a>
## [v2.21.0] - 2021-03-04
### Bug Fixes
- **alert_channel:** Ensure include_json_attachment is sent to API as `true` or `false` string
- **alert_policy_channel:** Ignore configured channel_id order  (convert to Set)

### Features
- Add NRQL Drop Rule support

<a name="v2.20.0"></a>
## [v2.20.0] - 2021-03-02
### Features
- **newrelic_entity:** Add ignore_case to name search for entity

<a name="v2.19.1"></a>
## [v2.19.1] - 2021-02-25
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

<a name="2.11.1"></a>
## [2.11.1] - 2020-10-07
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
- update development instructions for new TF version
- DEPRECATION notice for newrelic_alert_condition
- update supported Go information and test config
- update infra alert condition api key type
- update changelog
- **README:** update provider configuration pin version examples
- **dashboard:** add cross-account example
- **dashboard:** update docs with info regarding widget.account_id and cross-account widgets

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
- **provider:** add getting started guide to the quick links
- **provider:** fix incorrect newrelic_application reference in some examples
- **provider:** add account_id to argument reference, move argument reference above the fold
- **provider:** add environment variables and schema attribute table
- **provider:** update getting started example to reflect v2 updates
- **provider:** additional v2 updates, migration guide updates
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
- **alerts:** handle a nil reference with more grace
- **alerts:** infra alert condition zero value detection
- **application_settings:** Remove delete, as it is not possible
- **deps:** Revert terraform sdk to 1.10.0
- **newrelic:** fix the failing integration tests ([#519](https://github.com/newrelic/terraform-provider-newrelic/issues/519))
- **nrql_alert_condition:** threshold_occurrences is case insensitive, attribute description updates

### Documentation Updates
- update API key references to match desires
- prep for v2.x, isolate v1.x docs
- DEPRECATION notice for 1.x
- update index header with improved words
- include v1 index.html in sidebar
- update README with new pointers
- tidy up after review
- add table for current endpoint in use per resource
- update getting started guide to reference new material
- add callout to top of each v1.x doc page
- include documentation about upgrading the provider to 2.x
- **alert_channel:** fix broken 'nested config' anchor link
- **alerts:** update examples to reflect deprecation
- **alerts:** update wording to avoid implementation details
- **alerts:** include deprecation notice for "terms"
- **alerts:** include caveat about NRQL alerts condition operator usage with outliers
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
[Unreleased]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.25.0...HEAD
[v3.25.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.24.2...v3.25.0
[v3.24.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.24.1...v3.24.2
[v3.24.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.24.0...v3.24.1
[v3.24.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.23.0...v3.24.0
[v3.23.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.22.0...v3.23.0
[v3.22.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.21.3...v3.22.0
[v3.21.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.21.2...v3.21.3
[v3.21.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.21.1...v3.21.2
[v3.21.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.21.0...v3.21.1
[v3.21.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.20.2...v3.21.0
[v3.20.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.20.1...v3.20.2
[v3.20.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.20.0...v3.20.1
[v3.20.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.19.0...v3.20.0
[v3.19.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.18.1...v3.19.0
[v3.18.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.18.0...v3.18.1
[v3.18.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.17.1...v3.18.0
[v3.17.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.17.0...v3.17.1
[v3.17.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.16.1...v3.17.0
[v3.16.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.16.0...v3.16.1
[v3.16.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.15.0...v3.16.0
[v3.15.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.14.0...v3.15.0
[v3.14.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.13.0...v3.14.0
[v3.13.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.12.0...v3.13.0
[v3.12.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.11.0...v3.12.0
[v3.11.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.10.0...v3.11.0
[v3.10.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.9.0...v3.10.0
[v3.9.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.8.0...v3.9.0
[v3.8.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.7.1...v3.8.0
[v3.7.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.7.0...v3.7.1
[v3.7.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.6.1...v3.7.0
[v3.6.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.6.0...v3.6.1
[v3.6.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.5.2...v3.6.0
[v3.5.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.5.1...v3.5.2
[v3.5.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.5.0...v3.5.1
[v3.5.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.4.4...v3.5.0
[v3.4.4]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.4.3...v3.4.4
[v3.4.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.4.2...v3.4.3
[v3.4.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.4.1...v3.4.2
[v3.4.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.4.0...v3.4.1
[v3.4.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.3.0...v3.4.0
[v3.3.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.2.1...v3.3.0
[v3.2.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.2.0...v3.2.1
[v3.2.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.1.0...v3.2.0
[v3.1.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.0.4...v3.1.0
[v3.0.4]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.0.3...v3.0.4
[v3.0.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.0.2...v3.0.3
[v3.0.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.0.1...v3.0.2
[v3.0.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v3.0.0...v3.0.1
[v3.0.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/2.51.0...v3.0.0
[2.51.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.50.2...2.51.0
[v2.50.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.50.1...v2.50.2
[v2.50.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.50.0...v2.50.1
[v2.50.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.49.1...v2.50.0
[v2.49.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.49.0...v2.49.1
[v2.49.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.48.2...v2.49.0
[v2.48.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.48.1...v2.48.2
[v2.48.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.48.0...v2.48.1
[v2.48.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.47.1...v2.48.0
[v2.47.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.47.0...v2.47.1
[v2.47.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.46.2...v2.47.0
[v2.46.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.46.1...v2.46.2
[v2.46.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.46.0...v2.46.1
[v2.46.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.45.1...v2.46.0
[v2.45.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.45.0...v2.45.1
[v2.45.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.44.0...v2.45.0
[v2.44.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.43.4...v2.44.0
[v2.43.4]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.43.3...v2.43.4
[v2.43.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.43.2...v2.43.3
[v2.43.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.43.1...v2.43.2
[v2.43.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.43.0...v2.43.1
[v2.43.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.42.1...v2.43.0
[v2.42.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.42.0...v2.42.1
[v2.42.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.41.4...v2.42.0
[v2.41.4]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.41.3...v2.41.4
[v2.41.3]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.41.2...v2.41.3
[v2.41.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.41.1...v2.41.2
[v2.41.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.41.0-beta.2...v2.41.1
[v2.41.0-beta.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.41.0-beta.1...v2.41.0-beta.2
[v2.41.0-beta.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.40.0...v2.41.0-beta.1
[v2.40.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.39.2...v2.40.0
[v2.39.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.39.1...v2.39.2
[v2.39.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.39.0...v2.39.1
[v2.39.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.38.0...v2.39.0
[v2.38.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.37.0...v2.38.0
[v2.37.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.36.2...v2.37.0
[v2.36.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.36.1...v2.36.2
[v2.36.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.36.0...v2.36.1
[v2.36.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.35.1...v2.36.0
[v2.35.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.35.0...v2.35.1
[v2.35.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.34.1...v2.35.0
[v2.34.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.34.0...v2.34.1
[v2.34.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.33.0...v2.34.0
[v2.33.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.32.0...v2.33.0
[v2.32.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.31.1...v2.32.0
[v2.31.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.30.2...v2.31.1
[v2.30.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.30.1...v2.30.2
[v2.30.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.30.0...v2.30.1
[v2.30.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.29.0...v2.30.0
[v2.29.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.28.0...v2.29.0
[v2.28.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.27.1...v2.28.0
[v2.27.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.27.0...v2.27.1
[v2.27.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.26.0...v2.27.0
[v2.26.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.25.0...v2.26.0
[v2.25.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.24.1...v2.25.0
[v2.24.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.23.0...v2.24.1
[v2.23.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.22.1...v2.23.0
[v2.22.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.22.0...v2.22.1
[v2.22.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.21.2...v2.22.0
[v2.21.2]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.21.1...v2.21.2
[v2.21.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.21.0...v2.21.1
[v2.21.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.20.0...v2.21.0
[v2.20.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.19.1...v2.20.0
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
[v2.12.0]: https://github.com/newrelic/terraform-provider-newrelic/compare/2.11.1...v2.12.0
[2.11.1]: https://github.com/newrelic/terraform-provider-newrelic/compare/v2.11.1...2.11.1
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
