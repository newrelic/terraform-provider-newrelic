package main

type FileMappings map[string]FileMapping
type ProductMapping string

var productMappingMapKeysSorted []ProductMapping

type FileMapping struct {
	Test           bool   `yaml:"test"`
	ProductMapping string `yaml:"product_mapping"`
}

var ProductMappingTypes = struct {
	ALERTS               ProductMapping
	APM                  ProductMapping
	APIKS                ProductMapping
	AUTH                 ProductMapping
	CLOUD                ProductMapping
	DASHBOARDS           ProductMapping
	ENTITY               ProductMapping
	EVENTS               ProductMapping
	KeyTransactions      ProductMapping
	LoggingIntegrations  ProductMapping
	SYNTHETICS           ProductMapping
	WorkflowIntegrations ProductMapping
	WORKLOADS            ProductMapping
}{
	ALERTS:               "ALERTS",
	APIKS:                "APIKS",
	APM:                  "APM",
	AUTH:                 "AUTH",
	CLOUD:                "CLOUD",
	DASHBOARDS:           "DASHBOARDS",
	ENTITY:               "ENTITY",
	EVENTS:               "EVENTS",
	KeyTransactions:      "KEY_TRANSACTIONS",
	LoggingIntegrations:  "LOGGING_INTEGRATIONS",
	SYNTHETICS:           "SYNTHETICS",
	WorkflowIntegrations: "WORKFLOW_INTEGRATIONS",
	WORKLOADS:            "WORKLOADS",
}

var productMappings = map[ProductMapping][]string{
	ProductMappingTypes.ALERTS: {
		"alert",
	},
	ProductMappingTypes.APIKS: {
		"api_access",
	},
	ProductMappingTypes.APM: {
		"application",
	},
	ProductMappingTypes.AUTH: {
		"account_management",
		"authentication_domain",
		"group",
		"user",
	},
	ProductMappingTypes.CLOUD: {
		"link_account",
		"cloud",
		"cloud_account",
	},
	ProductMappingTypes.DASHBOARDS: {
		"dashboard",
	},
	ProductMappingTypes.ENTITY: {
		"entity",
	},
	ProductMappingTypes.EVENTS: {
		"event",
	},
	ProductMappingTypes.KeyTransactions: {
		"key_transaction",
	},
	ProductMappingTypes.LoggingIntegrations: {
		"data_partition",
		"drop_rule",
		"grok",
		"log_parsing_rule",
		"obfuscation",
	},
	ProductMappingTypes.SYNTHETICS: {
		"monitor_downtime",
		"synthetics",
	},
	ProductMappingTypes.WorkflowIntegrations: {
		"notification",
		"workflow",
	},
	ProductMappingTypes.WORKLOADS: {
		"service_level",
		"workload",
	},
}
