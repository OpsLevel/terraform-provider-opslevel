variable "extract_definition" {
  type        = string
  description = "The YAML definition for extracting data from inbound payloads."
  default     = null
}

variable "name" {
  type        = string
  description = "The name of the integration."
}

variable "transform_definition" {
  type        = string
  description = "The YAML definition for transforming extracted data to OpsLevel resources."
  default     = null
}
