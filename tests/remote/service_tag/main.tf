resource "opslevel_service_tag" "test" {
  key           = var.key
  value         = var.value
  service       = var.service
  service_alias = var.service_alias
}
