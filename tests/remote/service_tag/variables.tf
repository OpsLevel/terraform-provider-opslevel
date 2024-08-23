variable "key" {
  type        = string
  description = "The tag's key."
}

variable "service" {
  type        = string
  description = "The id of the service that this will be added to."
  default     = null
}

variable "service_alias" {
  type        = string
  description = "The alias of the service that this will be added to."
  default     = null
}

variable "value" {
  type        = string
  description = "The tag's value."
}
