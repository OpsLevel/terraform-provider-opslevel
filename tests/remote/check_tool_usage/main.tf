resource "opslevel_check_tool_usage" "test" {
  tool_category = var.tool_category
  environment_predicate = {
    type  = var.environment_predicate.type
    value = var.environment_predicate.value
  }
  tool_name_predicate = {
    type  = var.tool_name_predicate.type
    value = var.tool_name_predicate.value
  }
  tool_url_predicate = {
    type  = var.tool_url_predicate.type
    value = var.tool_url_predicate.value
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
