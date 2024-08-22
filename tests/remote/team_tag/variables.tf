variable "key" {
  type        = string
  description = "The tag's key."
}

variable "team" {
  type        = string
  description = "The id of the team that this will be added to."
  default     = null
}

variable "team_alias" {
  type        = string
  description = "The alias of the team that this will be added to."
  default     = null
}


variable "value" {
  type        = string
  description = "The tag's value."
}
