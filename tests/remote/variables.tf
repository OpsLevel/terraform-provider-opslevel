# Example usage in assert block
# error_message = replace(var.datasource_empty_error_msg, "TYPE", var.datasource_type)
variable "empty_datasource_error" {
  type    = string
  default = "zero 'TYPE' found in 'TYPE' datasource"
}

variable "unexpected_datasource_fields_error" {
  type    = string
  default = "cannot reference all expected 'TYPE' datasource fields"
}

variable "wrong_alias_error" {
  type    = string
  default = "wrong alias for 'TYPE'"
}

variable "wrong_index_error" {
  type    = string
  default = "wrong index for 'TYPE'"
}

variable "wrong_id_error" {
  type    = string
  default = "wrong id for 'TYPE'"
}

variable "wrong_name_error" {
  type    = string
  default = "wrong name for 'TYPE'"
}
