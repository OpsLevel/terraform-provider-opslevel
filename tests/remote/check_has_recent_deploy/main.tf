resource "opslevel_check_has_recent_deploy" "test" {
  days = var.days

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
