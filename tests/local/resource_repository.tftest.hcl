mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_repository_with_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_repository.with_alias.identifier == "github.com:rocktavious/autopilot"
    error_message = "wrong identifier for opslevel_repository.with_alias"
  }

  assert {
    condition     = opslevel_repository.with_alias.owner == null
    error_message = "expected null owner for opslevel_repository.with_alias"
  }

}


run "resource_repository_with_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_repository.with_id.id)
    error_message = "id attribute missing from opslevel_repository.with_id"
  }

  assert {
    condition     = opslevel_repository.with_id.identifier == var.test_id
    error_message = "wrong identifier for opslevel_repository.with_id"
  }

  assert {
    condition     = opslevel_repository.with_id.owner == var.test_id
    error_message = "wrong owner for opslevel_repository.with_id"
  }

}

