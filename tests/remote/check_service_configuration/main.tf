resource "opslevel_check_service_configuration" "test" {
  category  = var.category
  enable_on = var.enable_on
  enabled   = var.enabled
  filter    = var.filter
  level     = var.level
  name      = var.name
  notes     = var.notes
  owner     = var.owner
}
