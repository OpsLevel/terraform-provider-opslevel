data "opslevel_team" "devs" {
    alias = "developers"
}

resource "opslevel_domain" "example" {
  name = "Example"
  description = "The whole app in one monolith"
  owner = data.opslevel_team.devs.alias // or .id
}

resource "opslevel_system" "example" {
  name = "Example"
  description = "The async processing system for the monolith"
  owner = data.opslevel_team.devs.alias // or .id
  domain = data.opslevel_domain.example.alias // or .id
}
