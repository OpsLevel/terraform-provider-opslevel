mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_repository" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_repository.mock_repo.id)
    error_message = "id attribute missing from opslevel_repository.mock_repo"
  }

  assert {
    condition     = opslevel_repository.mock_repo.identifier == "github.com:rocktavious/autopilot"
    error_message = "wrong identifier for opslevel_repository.mock_repo"
  }

  assert {
    condition     = opslevel_repository.mock_repo.owner == "developers"
    error_message = "wrong owner for opslevel_repository.mock_repo"
  }

}

