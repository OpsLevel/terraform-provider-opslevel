resource "opslevel_check_custom_event" "test" {
  integration       = var.integration
  message           = var.message
  pass_pending      = var.pass_pending
  service_selector  = var.service_selector
  success_condition = var.success_condition

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
