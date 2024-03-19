mock_data "opslevel_filter" {
  defaults = {
    name = "mock-filter-name"
    # filter is not set here because its fields are not computed
    # id intentionally omitted - will be assigned a random string
  }
}
