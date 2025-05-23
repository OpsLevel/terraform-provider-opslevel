data "opslevel_team" "devs" {
  alias = "developers"
}

resource "opslevel_repository" "repo" {
  identifier      = "github.com:rocktavious/autopilot"
  owner           = data.opslevel_team.devs.id
  sbom_generation = "opt_in"
}