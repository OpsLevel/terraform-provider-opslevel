variable "description" {
  type        = string
  description = "The description of the Webhook Action."
}

variable "headers" {
  type        = map(string)
  description = "HTTP headers to be passed along with your webhook when triggered."
}

variable "method" {
  type        = string
  description = "The name for the system."

  validation {
    condition = var.method == null ? true : contains([
      "DELETE",
      "GET",
      "PATCH",
      "POST",
      "PUT",
    ], var.method)
    error_message = "expected method to be one of 'DELETE', 'GET', 'PATCH', 'POST', or 'PUT'"
  }
}

variable "name" {
  type        = string
  description = "The name of the Webhook Action."
}

variable "payload" {
  type        = string
  description = "Template that can be used to generate a webhook payload."
}

variable "url" {
  type        = string
  description = "The URL of the Webhook Action."
}

locals {
  method_options = [
    "DELETE",
    "GET",
    "PATCH",
    "POST",
    "PUT",
  ]
}
