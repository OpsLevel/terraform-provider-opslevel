variable "allowed_in_config_files" {
  type        = bool
  description = "Whether or not the property is allowed to be set in opslevel.yml config files."
}

variable "description" {
  type        = string
  description = "The description of the property definition."
}

variable "name" {
  type        = string
  description = "The display name of the property definition."
}

variable "property_display_status" {
  type        = string
  description = "The display status of a custom property on service pages."

  validation {
    condition = var.property_display_status == null ? true : contains([
      "hidden",
      "visible",
    ], var.property_display_status)
    error_message = "expected property_display_status to be 'hidden' or 'visible'"
  }
}

variable "schema" {
  type        = string
  description = "The schema of the property definition."

  validation {
    condition     = var.schema == null ? true : can(jsondecode(var.schema))
    error_message = "expected schema to be valid JSON"
  }
}
