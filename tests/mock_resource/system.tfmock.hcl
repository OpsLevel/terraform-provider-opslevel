mock_resource "opslevel_system" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    aliases      = ["fancy_system", "fancy_sys"]
    description  = "A Fancy API Client"
    domain       = "fancy_domain"
    last_updated = "2024-03-21T13:50:07Z"
    name         = "Fancy System"
    owner        = "developers"
  }
}
