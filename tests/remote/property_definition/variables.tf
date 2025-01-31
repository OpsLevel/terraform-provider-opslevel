variable "allowed_in_config_files" {
  type        = bool
  description = "Whether or not the property is allowed to be set in opslevel.yml config files."
}

variable "description" {
  type        = string
  description = "The description of the property definition."
  default     = null
}

variable "name" {
  type        = string
  description = "The display name of the property definition."
}

variable "property_display_status" {
  type        = string
  description = "The display status of a custom property on service pages."
}

variable "locked_status" {
  type        = string
  description = "Restricts what sources are able to assign values to this property."
}

variable "schema" {
  type        = string
  description = "The schema of the property definition."
}
