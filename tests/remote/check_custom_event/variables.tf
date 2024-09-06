variable "integration" {
  type        = string
  description = "The integration id this check will use."
}

variable "message" {
  type        = string
  description = "The check result message template. It is compiled with Liquid and formatted in Markdown."
  default     = null
}

variable "pass_pending" {
  type        = bool
  description = "True if this check should pass by default. Otherwise the default 'pending' state counts as a failure."
}

variable "service_selector" {
  type        = string
  description = "A jq expression that will be ran against your payload. This will parse out the service identifier."
}

variable "success_condition" {
  type        = string
  description = "A jq expression that will be ran against your payload. A truthy value will result in the check passing."
}
