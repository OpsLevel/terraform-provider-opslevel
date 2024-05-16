variable "affects_overall_service_levels" {
  type        = bool
  description = "Specifies whether the checks on this scorecard affect services' overall maturity level."
}

variable "description" {
  type        = string
  description = "The scorecard's description."
}

variable "filter_id" {
  type        = string
  description = "The scorecard's filter."
}

variable "name" {
  type        = string
  description = "The scorecard's name."
}

variable "owner_id" {
  type        = string
  description = "The scorecard's owner."

  validation {
    condition     = var.owner_id == null ? true : startswith(var.owner_id, "Z2lkOi8v")
    error_message = "expected owner_id to be a valid ID starting with 'Z2lkOi8v'"
  }
}
