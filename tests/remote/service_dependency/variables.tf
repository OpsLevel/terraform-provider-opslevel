variable "depends_upon" {
  type        = string
  description = "The ID or alias of the service that is depended upon."
}

variable "note" {
  type        = string
  description = "Notes for service dependency."
  default     = ""
}

variable "service" {
  type        = string
  description = "The ID or alias of the service with the dependency."
}
