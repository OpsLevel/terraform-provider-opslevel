resource "opslevel_check_tool_usage" "test" {
  tool_category         = var.tool_category
  environment_predicate = var.environment_predicate
  tool_name_predicate   = var.tool_name_predicate
  tool_url_predicate    = var.tool_url_predicate


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
