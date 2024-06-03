variable "aliases" {
  type        = list(string)
  description = "A list of human-friendly, unique identifiers for the service."

  validation {
    condition     = var.aliases == null ? true : var.aliases == distinct(var.aliases)
    error_message = "expected aliases to be unique"
  }
}

variable "api_document_path" {
  type        = string
  description = "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path."

  validation {
    condition = var.api_document_path == null ? true : anytrue([
      endswith(var.api_document_path, ".json"),
      endswith(var.api_document_path, ".yaml"),
      endswith(var.api_document_path, ".yml"),
    ])
    error_message = "expected api_document_path to end with '.json', '.yaml', or '.yml'"
  }
}

variable "description" {
  type        = string
  description = "A brief description of the service."
}

variable "framework" {
  type        = string
  description = "The primary software development framework that the service uses."
}

variable "language" {
  type        = string
  description = "The primary programming language that the service is written in."
}

variable "lifecycle_alias" {
  type        = string
  description = "The lifecycle stage of the service."
}

variable "name" {
  type        = string
  description = "The display name of the service."
}

variable "owner" {
  type        = string
  description = "The team that owns the service. ID or Alias may be used."
}

variable "preferred_api_document_source" {
  type        = string
  description = "The API document source used to determine the displayed document. If null, defaults to PUSH."

  validation {
    condition = var.preferred_api_document_source == null ? true : contains([
      "PULL",
      "PUSH",
    ], var.preferred_api_document_source)
    error_message = "expected preferred_api_document_source to be 'PULL' or 'PUSH'"
  }
}

variable "product" {
  type        = string
  description = "A product is an application that your end user interacts with. Multiple services can work together to power a single product."
}

variable "tags" {
  type        = set(string)
  description = "A list of unique tags applied to the service."
}

variable "tier_alias" {
  type        = string
  description = "The software tier that the service belongs to."
}
