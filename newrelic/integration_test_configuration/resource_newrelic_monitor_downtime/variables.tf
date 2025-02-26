variable "NEW_RELIC_ACCOUNT_ID" {
  type = number
}

variable "NEW_RELIC_API_KEY" {
  type      = string
  sensitive = true
}

variable "NEW_RELIC_REGION" {
  type    = string
  default = "US"
}

variable "name" {
  type = string
}

variable "mode" {
  type = string
}

variable "monitor_guids" {
  type = list(string)
}

variable "start_time" {
  type = string
}

variable "end_time" {
  type = string
}

variable "time_zone" {
  type = string
}

variable "include_end_repeat" {
  type = bool
  default = false
}

variable "end_repeat_on_date" {
  type = string
  default = ""
}

variable "end_repeat_on_repeat" {
  type = number
  default = -1
}

variable "maintenance_days" {
  type = list(string)
  default = []
}

variable "include_frequency" {
  type = bool
  default = false
}

variable "frequency_days_of_month" {
  type = list(number)
  default = []
}

variable "days_of_week_ordinal_day_of_month" {
  type = string
  default = ""
}

variable "days_of_week_week_day" {
  type = string
  default = ""
}
