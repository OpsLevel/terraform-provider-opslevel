variable "base_directory" {
  type        = string
  description = "The directory in the repository containing opslevel.yml."
  default     = null
}

variable "name" {
  type        = string
  description = "The name displayed in the UI for the service repository."
  default     = null
}

variable "repository" {
  type        = string
  description = "The id of the repository that this will be added to."
  default     = null
}

variable "repository_alias" {
  type        = string
  description = "The alias of the repository that this will be added to."
  default     = null
}

variable "service" {
  type        = string
  description = "The id of the service that this will be added to."
  default     = null
}

variable "service_alias" {
  type        = string
  description = "The alias of the service that this will be added to."
  default     = null
}

