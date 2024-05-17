variable "category" {
  type        = string
  description = "The id of the category the check belongs to."

  validation {
    condition     = var.category == null ? true : startswith(var.category, "Z2lkOi8v")
    error_message = "expected category to be a valid ID starting with 'Z2lkOi8v'"
  }
}

variable "enable_on" {
  type        = bool
  description = "The date when the check will be automatically enabled."
}

variable "enabled" {
  type        = bool
  description = "Whether the check is enabled or not.  Do not use this field in tandem with 'enable_on'."
}

variable "filter" {
  type        = string
  description = "The id of the filter of the check."

  validation {
    condition     = var.filter == null ? true : startswith(var.filter, "Z2lkOi8v")
    error_message = "expected filter to be a valid ID starting with 'Z2lkOi8v'"
  }
}

variable "level" {
  type        = string
  description = "The id of the level the check belongs to."

  validation {
    condition     = var.level == null ? true : startswith(var.level, "Z2lkOi8v")
    error_message = "expected level to be a valid ID starting with 'Z2lkOi8v'"
  }
}

variable "name" {
  type        = string
  description = "The display name of the check."
}

variable "notes" {
  type        = string
  description = "Additional information to display to the service owner about the check."
}

variable "owner" {
  type        = string
  description = "The id of the team that owns the check."

  validation {
    condition     = var.owner == null ? true : startswith(var.owner, "Z2lkOi8v")
    error_message = "expected owner to be a valid ID starting with 'Z2lkOi8v'"
  }
}
