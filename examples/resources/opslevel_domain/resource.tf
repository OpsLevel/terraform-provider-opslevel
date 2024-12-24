data "opslevel_team" "devs" {
  alias = "developers"
}

resource "opslevel_domain" "example" {
  name        = "Example"
  description = "The whole app in one monolith"
  owner       = data.opslevel_team.devs.id // or .aliases[0]
  note        = "This is an example"
}
