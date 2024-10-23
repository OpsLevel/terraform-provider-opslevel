variable "command" {
  type        = string
  description = "OpsLevel CLI command to run - everything after 'opslevel'"
  default     = null
}

variable "resource_id" {
  type    = string
  default = null
}

variable "resource_type" {
  type    = string
  default = null
}
