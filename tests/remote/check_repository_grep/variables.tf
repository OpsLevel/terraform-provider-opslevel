variable "directory_search" {
  type        = bool
  description = "Whether the check looks for the existence of a directory instead of a file."
}

variable "file_contents_predicate" {
  type = object({
    type  = string
    value = string
  })
  description = "A condition that should be satisfied."

  validation {
    condition = var.file_contents_predicate == null ? true : contains([
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
    ], var.file_contents_predicate.type)
    error_message = "invalid predicate type for file_contents_predicate.type"
  }
}

variable "filepaths" {
  type        = list(string)
  description = "Restrict the search to certain file paths."
}
