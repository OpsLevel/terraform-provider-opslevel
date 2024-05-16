variable "description" {
  type        = string
  description = "The description for the system."
}

variable "domain_id" {
  type        = string
  description = "The id of the parent domain this system is a child for."

  validation {
    condition     = var.domain_id == null ? true : startswith(var.domain_id, "Z2lkOi8v")
    error_message = "expected domain_id to be a valid ID starting with 'Z2lkOi8v'"
  }
}

variable "name" {
  type        = string
  description = "The name for the system."
}

variable "note" {
  type        = string
  description = "Additional information about the system."
}

variable "owner_id" {
  type        = string
  description = "The id of the team that owns the system."

  validation {
    condition     = var.owner_id == null ? true : startswith(var.owner_id, "Z2lkOi8v")
    error_message = "expected owner_id to be a valid ID starting with 'Z2lkOi8v'"
  }
}
