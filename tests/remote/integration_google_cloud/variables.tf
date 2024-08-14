variable "client_email" {
  type = string
}

variable "private_key" {
  type = string
}

variable "name" {
  type = string
}

variable "ownership_tag_keys" {
  type = list(string)
}

variable "ownership_tag_overrides" {
  type = bool
}
