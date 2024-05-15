variable "predicate_key_enum" {
  type        = string
  description = "fields that can be used as part of a filter for services"
  default     = "tier_index"

  validation {
    condition = contains([
      "aliases",
      "creation_source",
      "domain_id",
      "filter_id",
      "framework",
      "group_ids",
      "language",
      "lifecycle_index",
      "name",
      "owner_id",
      "owner_ids",
      "product",
      "properties",
      "repository_ids",
      "system_id",
      "tags",
      "tier_index",
    ], var.predicate_key_enum)
    error_message = "unknown predicate_key given"
  }
}

variable "predicate_type_enum" {
  type        = string
  description = "operations that can be used on predicates"
  default     = "equals"

  validation {
    condition = contains([
      "belongs_to",
      "contains",
      "does_not_contain",
      "does_not_equal",
      "does_not_exist",
      "does_not_match",
      "does_not_match_regex",
      "ends_with",
      "equals",
      "exists",
      "greater_than_or_equal_to",
      "less_than_or_equal_to",
      "matches",
      "matches_regex",
      "satisfies_jq_expression",
      "satisfies_version_constraint",
      "starts_with",
    ], var.predicate_type_enum)
    error_message = "unknown predicate_type given"
  }
}
