variable "file_extensions" {
  type        = set(string)
  description = "Restrict the search to files of given extensions. Extensions should contain only letters and numbers. For example: [\"py\", \"rb\"]."
}

variable "file_contents_predicate" {
  type = object({
    type  = string
    value = optional(string)
  })
  description = "A condition that should be satisfied."
}
