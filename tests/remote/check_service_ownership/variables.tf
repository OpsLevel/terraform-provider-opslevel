variable "contact_method" {
  type        = string
  description = "The type of contact method that is required."

  validation {
    condition = var.contact_method == null ? true : contains([
      "any",
      "email",
      "github",
      "slack",
      "slack_handle",
      "web",
    ], var.contact_method)
    error_message = "invalid predicate type for contact_method"
  }
}

variable "require_contact_method" {
  type        = bool
  description = "True if a service's owner must have a contact method, False otherwise."
}

variable "tag_key" {
  type        = string
  description = "The tag key where the tag predicate should be applied."
}

variable "tag_predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.tag_predicate == null ? true : contains([
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
    ], var.tag_predicate.type)
    error_message = "invalid predicate type for tag_predicate.type"
  }
}
