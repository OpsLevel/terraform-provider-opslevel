resource "opslevel_team_tag" "test" {
  key        = var.key
  value      = var.value
  team       = var.team
  team_alias = var.team_alias
}
