mock_data "opslevel_system" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    aliases     = ["my_system", "my_sys"]
    description = "This is my new system that has a domain."
    domain      = "sys_domain"
    name        = "My New System"
    owner       = "sys_owner"
  }
}

mock_data "opslevel_systems" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    systems = [
      {
        aliases     = ["my_system", "my_sys"]
        description = "This is my new system that has a domain."
        domain      = "sys_domain"
        name        = "My New System"
        owner       = "sys_owner"
      },
      {
        aliases     = ["my_system_2", "my_sys_2"]
        description = "This is my new system that has a domain (2)."
        domain      = "sys_domain"
        name        = "My New System (2)"
        owner       = "sys_owner"
      },
    ]
  }
}
