variable "alert_name_predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.alert_name_predicate == null ? true : contains([
      "contains",
      "does_not_contain",
      "does_not_equal",
      "does_not_exist",
      "ends_with",
      "equals",
      "exists",
      "greater_than_or_equal_to",
      "less_than_or_equal_to",
      "starts_with",
      "satisfies_version_constraint",
      "matches_regex",
      "does_not_match_regex",
      "belongs_to",
      "matches",
      "does_not_match",
      "satisfies_jq_expression",
    ], var.alert_name_predicate.type)
    error_message = "invalid predicate type for alert_name_predicate.type"
  }
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
