variable "etl_definition" {
  type = object({
    extract_definition   = string
    transform_definition = string
  })
  description = "The YAML definitions for extracting and transforming data from the integration."
  default     = null
}

variable "name" {
  type        = string
  description = "The name of the integration."
}
