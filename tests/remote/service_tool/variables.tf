variable "category" {
  type        = string
  description = "The category that the tool belongs to."
}

variable "environment" {
  type        = string
  description = "The environment that the tool belongs to."
  default     = null
}

variable "name" {
  type        = string
  description = "The display name of the tool."
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

variable "url" {
  type        = string
  description = "The URL of the tool."
}
