variable "update_frequency" {
  type = object({
    starting_date = string
    time_scale    = string
    value         = number
  })
  description = "Defines the minimum frequency of the updates."
}

variable "update_requires_comment" {
  type        = bool
  description = "Whether the check requires a comment or not."
}
