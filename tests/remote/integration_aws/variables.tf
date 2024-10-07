variable "external_id" {
  type        = string
  description = "The External ID defined in the trust relationship to ensure OpsLevel is the only third party assuming this role (See https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html for more details)."
}

variable "iam_role" {
  type        = string
  description = "The IAM role OpsLevel uses in order to access the AWS account."
}

variable "ownership_tag_keys" {
  type        = list(string)
  description = "An Array of tag keys used to associate ownership from an integration. Max 5"
  default     = null
}

variable "ownership_tag_overrides" {
  type        = bool
  description = "Allow tags imported from AWS to override ownership set in OpsLevel directly."
  default     = null
}

variable "name" {
  type        = string
  description = "The name of the integration."
}

variable "region_override" {
  type        = list(string)
  description = "Overrides the AWS region(s) that will be synchronized by this integration."
  default     = null
}
