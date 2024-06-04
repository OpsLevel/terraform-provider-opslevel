variable "contact_method" {
  type        = string
  description = "The type of contact method that is required."

  validation {
    condition = var.contact_method == null ? true : contains([
      "any",
      "email",
      "github",
      "slack",
      "slack_handle",
      "web",
    ], var.contact_method)
    error_message = "invalid predicate type for contact_method"
  }
}

variable "require_contact_method" {
  type        = bool
  description = "True if a service's owner must have a contact method, False otherwise."
}

variable "tag_key" {
  type        = string
  description = "The tag key where the tag predicate should be applied."
}

variable "tag_predicate" {
  type = object({
    type  = string
    value = optional(string)
  })
  description = "A condition that should be satisfied."
}
