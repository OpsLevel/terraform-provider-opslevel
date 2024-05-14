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

variable "name" {
  type        = string
  description = "The filter's display name."
}

