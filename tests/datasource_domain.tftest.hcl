mock_provider "opslevel" {
  alias = "fake"
  source = "./mock_datasource"
}

run "datasource_domain" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_domain.example.aliases[0] == "mock-alias-one"
    error_message = "wrong first alias in opslevel_domain.aliases"
  }

  assert {
    condition     = data.opslevel_domain.example.aliases[1] == "mock-alias-two"
    error_message = "wrong second alias in opslevel_domain.aliases"
  }

  assert {
    condition     = length(data.opslevel_domain.example.aliases) == 2
    error_message = "wrong number of aliases in opslevel_domain.aliases"
  }

  assert {
    condition     = data.opslevel_domain.example.description == "mock-domain-description"
    error_message = "wrong description in opslevel_domain.description"
  }

  assert {
    condition     = data.opslevel_domain.example.id != null && data.opslevel_domain.example.id != ""
    error_message = "opslevel_domain id should not be empty"
  }

  assert {
    condition     = data.opslevel_domain.example.name == "mock-domain-name"
    error_message = "wrong name in opslevel_domain.name"
  }

  assert {
    condition     = data.opslevel_domain.example.owner == null && data.opslevel_domain.example.owner != ""
    error_message = "opslevel_domain owner should be null"
  }
}
