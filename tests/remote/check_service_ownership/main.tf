resource "opslevel_check_service_ownership" "test" {
  contact_method         = var.contact_method
  require_contact_method = var.require_contact_method
  tag_key                = var.tag_key
  tag_predicate          = var.tag_predicate

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
