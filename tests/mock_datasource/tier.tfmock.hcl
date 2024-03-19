mock_data "opslevel_tier" {
  defaults = {
    alias = "mock-tier-alias"
    name  = "mock-tier-name"
    # filter is not set here because its fields are not computed
    # id intentionally omitted - will be assigned a random string
    # index intentionally omitted - defaults to number type 0
  }
}

