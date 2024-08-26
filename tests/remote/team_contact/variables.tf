variable "name" {
  type        = string
  description = "The name shown in the UI for the contact."
}

variable "team" {
  type        = string
  description = "The id or alias of the team the contact belongs to."
}

variable "type" {
  type        = string
  description = "The method of contact."
}

variable "value" {
  type        = string
  description = "The contact value."
}
