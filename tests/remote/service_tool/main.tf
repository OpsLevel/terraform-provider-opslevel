resource "opslevel_service_tool" "test" {
  category      = var.category
  environment   = var.environment
  name          = var.name
  service       = var.service
  service_alias = var.service_alias
  url           = var.url
}
