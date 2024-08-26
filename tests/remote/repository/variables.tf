variable "identifier" {
  type        = string
  description = "The id or human-friendly, unique identifier for the repository."
}

variable "owner_id" {
  type        = string
  description = "The id of the owner of the repository."
  default     = null
}
