mock_data "opslevel_tier" {
  defaults = {
    alias = "mock-tier-alias"
    # filter is not set here because its fields are not computed
    # id intentionally omitted - will be assigned a random string
    index = 321
    name  = "mock-tier-name"
  }
}

