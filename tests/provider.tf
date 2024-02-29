provider "opslevel" {
  api_url     = var.OPSLEVEL_API_URL
  api_token   = var.OPSLEVEL_API_TOKEN
  api_timeout = var.OPSLEVEL_API_TIMEOUT
}

variable "OPSLEVEL_API_URL" {
  type    = string
  default = "https://app.opslevel.com"
}

variable "OPSLEVEL_API_TOKEN" {
  type      = string
  sensitive = true
}

variable "OPSLEVEL_API_TIMEOUT" {
  type    = number
  default = 10
}