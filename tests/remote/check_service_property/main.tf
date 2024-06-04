resource "opslevel_check_service_property" "test" {
  property = var.property
  predicate = {
    type  = var.predicate.type
    value = var.predicate.value
  }

  # -- check base fields --
  category  = var.category
  enable_on = var.enable_on
  enabled   = var.enabled
  filter    = var.filter
  level     = var.level
  name      = var.name
  notes     = var.notes
  owner     = var.owner
}
