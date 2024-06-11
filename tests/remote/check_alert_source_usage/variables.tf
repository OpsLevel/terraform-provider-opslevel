variable "alert_name_predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}

variable "alert_type" {
  type        = string
  description = "The type of the alert source."

  validation {
    condition = var.alert_type == null ? true : contains([
      "pagerduty",
      "datadog",
      "opsgenie",
      "new_relic",
    ], var.alert_type)
    error_message = "expected level to be a valid ID starting with 'Z2lkOi8v'"
  }
}
