variable "description" {
  type        = string
  description = "The description of the rubric level."
}

variable "index" {
  type        = number
  description = "An integer allowing this level to be inserted between others. Must be unique per rubric."
}

variable "name" {
  type        = string
  description = "The display name of the rubric level."
}
