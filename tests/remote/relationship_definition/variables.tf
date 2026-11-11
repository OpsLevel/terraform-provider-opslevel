variable "alias" {
  type        = string
  description = "The unique identifier of the relationship definition."
}

variable "component_type_alias" {
  type        = string
  description = "The alias of the component type the relationship definition belongs to."
}

variable "component_type_name" {
  type        = string
  description = "The name of the component type the relationship definition belongs to."
}

variable "component_type_reference" {
  type        = string
  description = "How the relationship definition names its component type: by alias, by ID, or by its alias in upper case."

  validation {
    condition     = contains(["alias", "id", "upper_alias"], var.component_type_reference)
    error_message = "component_type_reference must be alias, id or upper_alias"
  }
}

variable "name" {
  type        = string
  description = "The display name of the relationship definition."
}
