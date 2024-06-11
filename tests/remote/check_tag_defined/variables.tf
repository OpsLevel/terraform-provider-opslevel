variable "tag_key" {
  type        = string
  description = "The tag key where the tag predicate should be applied."
}

variable "tag_predicate" {
  type = object({
    type  = string
    value = optional(string, "")
  })
  description = "A condition that should be satisfied."
}
