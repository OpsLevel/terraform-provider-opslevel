mock_resource "opslevel_system" {
  defaults = {
    # id intentionally omitted - will be assigned a random string
    aliases      = ["fancy-system"]
    last_updated = "2022-02-24T13:50:07Z"
    name         = "Example"
    description  = "The whole app in one monolith"
    note         = "This is an example"
    owner        = "Developers"
  }
}
