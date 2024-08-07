variable "aliases" {
  type        = list(string)
  description = "A list of human-friendly, unique identifiers for the team."
  default     = null
}

variable "name" {
  type        = string
  description = "The team's display name."
}

variable "parent" {
  type        = string
  description = "The id or alias of the parent team."
  default     = null
}

variable "responsibilities" {
  type        = string
  description = "A description of what the team is responsible for."
  default     = null
}
