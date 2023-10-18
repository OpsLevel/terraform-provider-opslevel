data "opslevel_team" "devs" {
  alias = "developers"
}

resource "opslevel_team_contact" "slack" {
  team  = data.opslevel_team.devs.alias
  type  = "slack"
  name  = "Slack"
  value = "#devs"
}

resource "opslevel_team_contact" "email" {
  team  = data.opslevel_team.devs.alias
  type  = "email"
  name  = "Email"
  value = "developers@example.com"
}

resource "opslevel_team_contact" "example" {
  team  = data.opslevel_team.devs.alias
  type  = "web"
  name  = "Gitlab"
  value = "https://gitlab.com/groups/example/-/issues"
}
