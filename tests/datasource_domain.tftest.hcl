mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_domain" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_domain.mock_domain.aliases[0] == "mock-alias-one"
    error_message = "wrong first alias in opslevel_domain.aliases"
  }

  assert {
    condition     = data.opslevel_domain.mock_domain.aliases[1] == "mock-alias-two"
    error_message = "wrong second alias in opslevel_domain.aliases"
  }

  assert {
    condition     = length(data.opslevel_domain.mock_domain.aliases) == 2
    error_message = "wrong number of aliases in opslevel_domain.aliases"
  }

  assert {
    condition     = data.opslevel_domain.mock_domain.description == "mock-domain-description"
    error_message = "wrong description in opslevel_domain.description"
  }

  assert {
    condition     = data.opslevel_domain.mock_domain.id != null && data.opslevel_domain.mock_domain.id != ""
    error_message = "opslevel_domain id should not be empty"
  }

  assert {
    condition     = data.opslevel_domain.mock_domain.name == "mock-domain-name"
    error_message = "wrong name in opslevel_domain.name"
  }

  assert {
    condition     = data.opslevel_domain.mock_domain.owner == null && data.opslevel_domain.mock_domain.owner != ""
    error_message = "opslevel_domain owner should be null"
  }
}

run "datasource_domains_all" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = length(data.opslevel_domains.all.domains) == 2
    error_message = "wrong number of owners in opslevel_domains"
  }

  assert {
    condition     = data.opslevel_domains.all.domains[1].description == ""
    error_message = "wrong description in second opslevel_domain"
  }
}
