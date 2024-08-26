resource "opslevel_service_repository" "test" {
  base_directory   = var.base_directory
  name             = var.name
  repository       = var.repository
  repository_alias = var.repository_alias
  service          = var.service
  service_alias    = var.service_alias
}
