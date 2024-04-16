mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_domain" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_domain.fancy.aliases[0] == "fancy-domain"
    error_message = "wrong first alias in opslevel_domain.aliases"
  }

  assert {
    condition     = opslevel_domain.fancy.name == "Example"
    error_message = "wrong name for opslevel_domain"
  }

  assert {
    condition     = opslevel_domain.fancy.id != null && opslevel_domain.fancy.id != ""
    error_message = "opslevel_domain id should not be empty"
  }

  assert {
    condition     = opslevel_domain.fancy.owner == "Developers"
    error_message = "wrong owner of opslevel_domain resource"
  }
}

