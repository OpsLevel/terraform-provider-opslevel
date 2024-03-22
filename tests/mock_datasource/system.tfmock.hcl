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
