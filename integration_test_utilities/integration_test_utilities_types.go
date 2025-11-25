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
	FLEET                ProductMapping
	KeyTransactions      ProductMapping
	LoggingIntegrations  ProductMapping
	NGEP                 ProductMapping
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
	FLEET:                "FLEET",
	KeyTransactions:      "KEY_TRANSACTIONS",
	LoggingIntegrations:  "LOGGING_INTEGRATIONS",
	NGEP:                 "NGEP",
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
		"cloud_aws",
		"cloud_azure",
		"cloud_gcp",
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
	ProductMappingTypes.FLEET: {
		"fleet",
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
	ProductMappingTypes.NGEP: {
		"pipeline_cloud_rule",
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
