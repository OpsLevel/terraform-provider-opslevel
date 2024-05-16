variable "description" {
  type        = string
  description = "The description of the domain."
}

variable "name" {
  type        = string
  description = "The name of the domain."
}

variable "note" {
  type        = string
  description = "Additional information about the domain."
}

variable "owner_id" {
  type        = string
  description = "The id of the team that owns the domain."

  validation {
    condition     = var.owner_id == null ? true : startswith(var.owner_id, "Z2lkOi8v")
    error_message = "expected owner_id to be a valid ID starting with 'Z2lkOi8v'"
  }
}
