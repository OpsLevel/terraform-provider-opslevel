variable "property" {
  type        = string
  description = "The property of the service that the check will verify."

  validation {
    condition = var.property == null ? true : contains([
      "custom_property",
      "description",
      "framework",
      "language",
      "lifecycle_index",
      "name",
      "note",
      "product",
      "system",
      "tier_index",
    ], var.property)
    error_message = "invalid property type"
  }
}

variable "predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.predicate == null ? true : contains([
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
    ], var.predicate.type)
    error_message = "invalid predicate type for predicate.type"
  }
}
