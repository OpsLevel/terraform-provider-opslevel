variable "aliases" {
  type        = list(string)
  description = "A list of human-friendly, unique identifiers for the service."
  default     = []
}

variable "api_document_path" {
  type        = string
  description = "The relative path from which to fetch the API document. If null, the API document is fetched from the account's default path."
  default     = null
}

variable "description" {
  type        = string
  description = "A brief description of the service."
  default     = null
}

variable "framework" {
  type        = string
  description = "The primary software development framework that the service uses."
  default     = null
}

variable "language" {
  type        = string
  description = "The primary programming language that the service is written in."
  default     = null
}

variable "lifecycle_alias" {
  type        = string
  description = "The lifecycle stage of the service."
  default     = null
}

variable "name" {
  type        = string
  description = "The display name of the service."
}

variable "owner" {
  type        = string
  description = "The team that owns the service. ID or Alias may be used."
  default     = null
}

variable "parent" {
  type        = string
  description = "The id or alias of the parent system of this service"
  default     = null
}

variable "preferred_api_document_source" {
  type        = string
  description = "The API document source used to determine the displayed document. If null, defaults to PUSH."
  default     = null
}

variable "product" {
  type        = string
  description = "A product is an application that your end user interacts with. Multiple services can work together to power a single product."
  default     = null
}

variable "tags" {
  type        = set(string)
  description = "A list of unique tags applied to the service."
  default     = []
}

variable "tier_alias" {
  type        = string
  description = "The software tier that the service belongs to."
  default     = null
}
