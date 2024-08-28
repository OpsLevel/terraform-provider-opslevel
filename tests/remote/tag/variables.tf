variable "key" {
  type        = string
  description = "The key of the tag."
}

variable "resource_identifier" {
  type        = string
  description = "The id or human-friendly, unique identifier of the resource this tag belongs to."
}

variable "resource_type" {
  type        = string
  description = "The resource type that the tag applies to."
}

variable "value" {
  type        = string
  description = "The value of the tag."
}
