variable "service" {
	description = "The service to create alerts for"
	type = object({
		name                       = string
		threshold_duration         = number
		cpu_threshold              = number
		response_time_threshold    = number
		error_percentage_threshold = number
		throughput_threshold       = number
	})
}

variable "notification_channel_ids" {
	description = "The notification channel IDs to add to this policy"
	type        = list(string)
}
