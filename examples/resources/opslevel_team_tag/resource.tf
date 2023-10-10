data "opslevel_team" "foo" {
  alias = "foo"
}

resource "opslevel_tag" "foo" {
  type       = "Team"
  identifier = data.opslevel_team.foo.id

  key   = "environment"
  value = "foo"
}

resource "opslevel_tag" "bar" {
  type       = "Team"
  identifier = data.opslevel_team.foo.id

  key   = "environment"
  value = "bar"
}