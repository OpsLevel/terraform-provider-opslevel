variable "document_subtype" {
  type        = string
  description = "The subtype of the document."

  validation {
    condition     = var.document_subtype == null ? true : contains(["openapi"], var.document_subtype)
    error_message = "expected document_type to be 'openapi'"
  }
}

variable "document_type" {
  type        = string
  description = "The type of the document."

  validation {
    condition     = var.document_type == null ? true : contains(["api", "tech"], var.document_type)
    error_message = "expected document_type to be 'api' or tech'"
  }
}
