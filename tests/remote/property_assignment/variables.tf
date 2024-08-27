variable "definition" {
  type        = string
  description = "The custom property definition's ID or alias."
}

variable "owner" {
  type        = string
  description = "The ID or alias of the entity (currently only supports service) that the property has been assigned to."
}

variable "value" {
  type        = string
  description = "The value of the custom property (must be a valid JSON value or null or object)."
  default     = null
}
