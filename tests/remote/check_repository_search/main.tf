resource "opslevel_check_repository_search" "test" {
  # file_contents_predicate = var.file_contents_predicate
  file_contents_predicate = {
    type  = "contains"
    value = "dev"
  }
  file_extensions = var.file_extensions

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
