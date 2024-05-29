variable "environment_predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.environment_predicate == null ? true : contains([
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
    ], var.environment_predicate.type)
    error_message = "invalid predicate type for environment_predicate.type"
  }
}

variable "tool_category" {
  type        = string
  description = "The category that the tool belongs to."

  validation {
    condition = var.tool_category == null ? true : contains([
      "admin",
      "api_documentation",
      "architecture_diagram",
      "backlog",
      "code",
      "continuous_integration",
      "deployment",
      "design_documentation",
      "errors",
      "feature_flag",
      "health_checks",
      "incidents",
      "issue_tracking",
      "logs",
      "metrics",
      "observability",
      "orchestrator",
      "other",
      "resiliency",
      "runbooks",
      "security_scans",
      "status_page",
      "wiki",
    ], var.tool_category)
    error_message = "expected level to be a valid ID starting with 'Z2lkOi8v'"
  }
}

variable "tool_name_predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.tool_name_predicate == null ? true : contains([
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
    ], var.tool_name_predicate.type)
    error_message = "invalid predicate type for tool_name_predicate.type"
  }
}

variable "tool_url_predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.tool_url_predicate == null ? true : contains([
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
    ], var.tool_url_predicate.type)
    error_message = "invalid predicate type for tool_url_predicate.type"
  }
}
