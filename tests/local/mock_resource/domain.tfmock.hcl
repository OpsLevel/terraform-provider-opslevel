mock_resource "opslevel_domain" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    aliases     = ["fancy-domain"]
    name        = "Example"
    description = "The whole app in one monolith"
    note        = "This is an example"
    owner       = "Developers"
  }
}
