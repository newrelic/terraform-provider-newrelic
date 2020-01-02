package apm

// KeyTransaction represents information about a New Relic key transaction.
type KeyTransaction struct {
	ID              int                       `json:"id,omitempty"`
	Name            string                    `json:"name,omitempty"`
	TransactionName string                    `json:"transaction_name,omitempty"`
	HealthStatus    string                    `json:"health_status,omitempty"`
	LastReportedAt  string                    `json:"last_reported_at,omitempty"`
	Reporting       bool                      `json:"reporting"`
	Summary         ApplicationSummary        `json:"application_summary,omitempty"`
	EndUserSummary  ApplicationEndUserSummary `json:"end_user_summary,omitempty"`
	Links           KeyTransactionLinks       `json:"links,omitempty"`
}

// KeyTransactionLinks represents associations for a key transaction.
type KeyTransactionLinks struct {
	Application int `json:"application,omitempty"`
}
