resource "opslevel_check_tag_defined" "test" {
  tag_key = var.tag_key
  tag_predicate = {
    type  = var.tag_predicate.type
    value = var.tag_predicate.value
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
