mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_team_contact" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition
  }
}