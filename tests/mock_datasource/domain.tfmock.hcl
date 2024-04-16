mock_data "opslevel_domain" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    name        = "mock-domain-name"
    description = "mock-domain-description"
    aliases     = ["mock-alias-one", "mock-alias-two"]
    owner       = null
  }
}

mock_data "opslevel_domains" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    domains = [
      {
        aliases = [
          "mock-alias-one",
          "mock-alias-two"
        ]
        name        = "mock-domain-name"
        description = "mock-domain-description"
        owner       = "mock-owner"
      },
      {
        aliases     = []
        name        = "fake-domain-name"
        description = ""
        owner       = null
      },
    ]
  }
}
