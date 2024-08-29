variable "aliases" {
  type        = set(string)
  description = "The aliases for the infrastructure resource."
  default     = null
}

variable "data" {
  type        = string
  description = "The data of the infrastructure resource in JSON format."
}

variable "owner" {
  type        = string
  description = "The id of the team that owns the infrastructure resource. Does not support aliases!"
}

variable "provider_data" {
  type = object({
    account = string
    name    = optional(string)
    type    = optional(string)
    url     = optional(string)
  })
  description = "The provider specific data for the infrastructure resource."
  default     = null
}

variable "schema" {
  type        = string
  description = "The schema of the infrastructure resource that determines its data specification."
}
