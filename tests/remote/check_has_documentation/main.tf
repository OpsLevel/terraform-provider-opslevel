resource "opslevel_check_has_documentation" "test" {
  document_type    = var.document_type
  document_subtype = var.document_subtype

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
