variable "environment_predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}

variable "tool_category" {
  type        = string
  description = "The category that the tool belongs to."
}

variable "tool_name_predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}

variable "tool_url_predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}
