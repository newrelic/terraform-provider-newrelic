## 1.6.1 (Unreleased)
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
