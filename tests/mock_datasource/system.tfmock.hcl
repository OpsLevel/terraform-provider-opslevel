mock_data "opslevel_system" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    aliases     = ["mock_alias_one", "mock_alias_two"]
    description = "Mock system description"
    domain      = "mock_domain"
    name        = "Mock System Name"
    owner       = "system_owner"
  }
}

mock_data "opslevel_systems" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    systems = [
      {
        aliases     = ["mock_alias_two"]
        description = "Mock system description"
        domain      = "mock_domain"
        name        = "Mock System Name"
        owner       = "system_owner"
      },
      {
        aliases     = ["mock_alias_three", "mock_alias_four"]
        description = "Mock system description the second"
        domain      = "mock_domain"
        name        = "Mock System Name The Second"
        owner       = "system_owner"
      },
    ]
  }
}
