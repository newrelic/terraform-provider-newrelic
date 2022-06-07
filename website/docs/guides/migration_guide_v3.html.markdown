---
layout: "newrelic"
page_title: "New Relic Terraform Provider v3.x Migration Guide"
sidebar_current: "docs-newrelic-provider-v3-migration-guide"
description: |-
  Use this guide to update the New Relic Terraform Provider from v2.x to v3.x
---

## Upgrade to v3.x of the New Relic Terraform Provider

Version 3.x of the provider uses a new underlying API for Synthetics. This results in some changes that will need to be made to existing Synthetics resource to keep them compatible with the new API.
