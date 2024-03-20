mock_data "opslevel_system" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    name        = "mock-system-name"
    description = "mock-system-description"
    aliases     = ["mock-alias-one", "mock-alias-two"]
    owner       = null
  }
}

mock_data "opslevel_systems" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    systems = [
      {
        aliases = [
          "mock-alias-one",
          "mock-alias-two"
        ]
        name        = "mock-system-name"
        description = "mock-system-description"
        owner       = "mock-owner"
      },
      {
        aliases     = []
        name        = "fake-system-name"
        description = ""
        owner       = null
      },
    ]
  }
}
