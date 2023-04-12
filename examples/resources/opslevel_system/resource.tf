data "opslevel_team" "devs" {
  alias = "developers"
}

resource "opslevel_domain" "example" {
  name = "Example"
  description = "The whole app in one monolith"
  owner = data.opslevel_team.devs.id
  note = "This is an example"
}

resource "opslevel_system" "example" {
  name = "Example"
  description = "The async processing system for the monolith"
  owner = data.opslevel_team.devs.id
  domain = opslevel_domain.example.id // or .aliases[0]
  note = "This is another example"
}
