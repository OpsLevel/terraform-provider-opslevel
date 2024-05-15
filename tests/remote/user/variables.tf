variable "email" {
  type        = string
  description = "The email address of the user."
}

variable "name" {
  type        = string
  description = "The name of the user."
}

variable "role" {
  type        = string
  description = "The access role of the user."
  validation {
    condition = var.role == null ? true : contains([
      "admin",
      "basic_user",
      "user",
    ], var.role)
    error_message = "expected role to be 'admin', 'basic_user', or 'user'"
  }
}
