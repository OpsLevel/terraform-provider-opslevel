variable "client_email" {
  type        = string
  description = "The service account email OpsLevel uses to access the Google Cloud account."
}

variable "name" {
  type        = string
  description = "The name of the integration."
}

variable "ownership_tag_keys" {
  type        = list(string)
  description = "An Array of tag keys used to associate ownership from an integration. Max 5 (default = [\"owner\"])"
  default     = null
}

variable "ownership_tag_overrides" {
  type        = bool
  description = "Allow tags imported from Google Cloud to override ownership set in OpsLevel directly. (default = true)"
  default     = null
}

variable "private_key" {
  type        = string
  sensitive   = true
  description = "The private key for the service account that OpsLevel uses to access the Google Cloud account."
}
