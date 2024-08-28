variable "missing_package_result" {
  type        = string
  description = "The check result if the package isn't being used by a service. (Optional.)"
  default     = null
}

variable "package_constraint" {
  type        = string
  description = "The package constraint the service is to be checked for. (Required.)"
}

variable "package_manager" {
  type        = string
  description = "The package manager (ecosystem) this package relates to. (Required.)"
}

variable "package_name" {
  type        = string
  description = "The name of the package to be checked. (Required.)"
}

variable "package_name_is_regex" {
  type        = bool
  description = "Whether or not the value in the package name field is a regular expression. (Optional.)"
  default     = null
}

variable "version_constraint_predicate" {
  type = object({
    type  = string
    value = optional(string)
  })
  description = "A condition that should be satisfied."
  default     = null
}
