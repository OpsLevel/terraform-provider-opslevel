mock_data "opslevel_integration" {
  defaults = {
    name = "My Mock Integration"
    # filter is not set here because its fields are not computed
    # id intentionally omitted - will be assigned a random string
  }
}
