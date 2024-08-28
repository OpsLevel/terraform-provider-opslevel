variable "alias" {
  type        = string
  description = "The alias for this secret. Can only be set at create time."
}

variable "owner" {
  type        = string
  description = "The owner of this secret."
}

variable "value" {
  type        = string
  description = "A sensitive value"
  sensitive   = true
}
