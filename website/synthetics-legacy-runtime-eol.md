Hello!

As already communicated by New Relic, support for legacy Synthetics runtimes **will reach end-of-life (EOL) on October 22, 2024**. In addition, creating **new** monitors using the legacy runtime **will no longer be supported after June 30, 2024**. This affects Synthetic Monitors running on the legacy runtime, created using any of the following resources used to manage Synthetic Monitors:

- [newrelic_synthetics_monitor](https://registry.terraform.io/providers/newrelic/newrelic/3.36.1/docs/resources/synthetics_monitor)
- [newrelic_synthetics_script_monitor](https://registry.terraform.io/providers/newrelic/newrelic/3.36.1/docs/resources/synthetics_script_monitor)
- [newrelic_synthetics_step_monitor](https://registry.terraform.io/providers/newrelic/newrelic/3.36.1/docs/resources/setics_step_monitor)
- [newrelic_synthetics_cert_check_monitor](https://registry.terraform.io/providers/newrelic/newrelic/3.36.1/docs/resources/synthetics_cert_check_monitor)
- [newrelic_synthetics_broken_links_monitor](https://registry.terraform.io/providers/newrelic/newrelic/3.36.1/docs/resources/synthetics_broken_links_monitor)

In light of the above, kindly **upgrade your Synthetic Monitors to the new runtime** before the EOL date. Check out [this page](https://forum.newrelic.com/s/hubtopic/aAXPh0000001brxOAA/upcoming-endoflife-legacy-synthetics-runtimes-and-cpm) for more details on the EOL, actions needed, relevant resources, and more.

We will synchronize relevant changes in the New Relic Terraform Provider with the Synthetics team. Watch this space for more updates.

The article linked above comprises key details around the EOL and has been published by New Relic Synthetics. We shall also continue to share any useful information around the EOL via this thread, or via the documentation of the New Relic Terraform Provider, if we find such information can be useful to a larger audience.

Should you have any questions/suggestions, please feel free to let us know in this thread.

Thanks a lot
The Observability as Code team
