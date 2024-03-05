mock_provider "opslevel" {
  alias = "fake"
  override_data {
    target = data.opslevel_domain.example
    values = {
      name = "mock-domain-name"
    }

  }
}

run "datasource_domain" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_domain.example.name == "mock-domain-name"
    error_message = "Wrong domain name"
  }
}
