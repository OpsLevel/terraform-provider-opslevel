variable "connective" {
  type        = string
  description = "logical operator to be used in conjunction with multiple filters"

  validation {
    condition = var.connective == null ? true : contains([
      "and",
      "or",
    ], var.connective)
    error_message = "expected connective_enum to be 'and' or 'or'"
  }
}

# WIP
variable "predicate_list" {
  type = list(object({
    case_insensitive = optional(bool)
    case_sensitive   = optional(bool)
    key              = string
    key_data         = optional(string)
    type             = string
    value            = optional(string)
  }))
  default = []
}

#variable "predicate" {
#  type = object({
#    case_insensitive = optional(bool)
#    case_sensitive   = optional(bool)
#    key              = string
#    key_data         = optional(string)
#    type             = string
#    value            = optional(string)
#  })
#  description = "The filter's predicate."
#
#  validation {
#    condition = alltrue(
#      contains([
#        "aliases",
#        "creation_source",
#        "domain_id",
#        "filter_id",
#        "framework",
#        "group_ids",
#        "language",
#        "lifecycle_index",
#        "name",
#        "owner_id",
#        "owner_ids",
#        "product",
#        "properties",
#        "repository_ids",
#        "system_id",
#        "tags",
#        "tier_index",
#      ], var.predicate.key),
#      contains([
#        "belongs_to",
#        "contains",
#        "does_not_contain",
#        "does_not_equal",
#        "does_not_exist",
#        "does_not_match",
#        "does_not_match_regex",
#        "ends_with",
#        "equals",
#        "exists",
#        "greater_than_or_equal_to",
#        "less_than_or_equal_to",
#        "matches",
#        "matches_regex",
#        "satisfies_jq_expression",
#        "satisfies_version_constraint",
#        "starts_with",
#      ], var.predicate.type),
#    )
#    error_message = "invalid predicate key or type given"
#  }
#}

variable "name" {
  type        = string
  description = "The filter's display name."
}
