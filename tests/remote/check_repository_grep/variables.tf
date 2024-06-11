# NOTE: if directory_search is true, file_contents_predicate can only be "exists" or "does_not_exist"
variable "directory_search" {
  type        = bool
  description = "Whether the check looks for the existence of a directory instead of a file."
}

variable "file_contents_predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}

variable "filepaths" {
  type        = list(string)
  description = "Restrict the search to certain file paths."

  validation {
    condition     = length(var.filepaths) > 0
    error_message = "expected at least one file path"
  }
}
