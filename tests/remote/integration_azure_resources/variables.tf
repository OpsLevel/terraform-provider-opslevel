variable "client_id" {
  type        = string
  description = "The client id OpsLevel uses to access the Azure account."
}

variable "client_secret" {
  type        = string
  sensitive   = true
  description = "The client secret OpsLevel uses to access the Azure account."
}

variable "name" {
  type        = string
  description = "The name of the integration."
}

variable "ownership_tag_keys" {
  type        = list(string)
  description = "An Array of tag keys used to associate ownership from an integration. Max 5"
  default     = null
}

variable "ownership_tag_overrides" {
  type        = bool
  description = "Allow tags imported from AWS to override ownership set in OpsLevel directly."
}

variable "subscription_id" {
  type        = string
  description = "The subscription OpsLevel uses to access the Azure account."
}

variable "tenant_id" {
  type        = string
  description = "The tenant OpsLevel uses to access the Azure account."
}

