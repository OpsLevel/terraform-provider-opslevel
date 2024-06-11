variable "property" {
  type        = string
  description = "The property of the service that the check will verify."

  validation {
    condition = var.property == null ? true : contains([
      "custom_property",
      "description",
      "framework",
      "language",
      "lifecycle_index",
      "name",
      "note",
      "product",
      "system",
      "tier_index",
    ], var.property)
    error_message = "invalid property type"
  }
}

variable "predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}
