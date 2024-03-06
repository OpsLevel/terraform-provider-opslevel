mock_data "opslevel_domain" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    name        = "mock-domain-name"
    description = "mock-domain-description"
    aliases     = ["mock-alias-one", "mock-alias-two"]
    owner       = null
  }
}
