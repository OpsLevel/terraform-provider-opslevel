mock_data "opslevel_rubric_level" {
  defaults = {
    alias = "mock-rubric-level-alias"
    # filter is not set here because its fields are not computed
    # id intentionally omitted - will be assigned a random string
    index = 321
    name  = "mock-rubric-level-name"
  }
}
