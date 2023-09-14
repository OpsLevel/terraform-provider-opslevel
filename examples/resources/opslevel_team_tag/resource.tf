data "opslevel_team" "foo" {
  alias = "foo"
}

resource "opslevel_team_tag" "foo_environment" {
  team = data.opslevel_team.foo.id

  key   = "type"
  value = "frontend"
}

resource "opslevel_team_tag" "bar_environment" {
  team_alias = "bar"

  key   = "type"
  value = "frontend"
}