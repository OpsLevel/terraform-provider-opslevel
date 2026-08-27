variable "campaign_name" {
  type = string
}

variable "owner_id" {
  type = string
}

variable "source_check_name" {
  type = string
}

variable "notes" {
  type    = string
  default = null
}

variable "enabled" {
  type    = bool
  default = false
}
