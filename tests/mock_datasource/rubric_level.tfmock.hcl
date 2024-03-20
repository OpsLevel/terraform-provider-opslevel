mock_data "opslevel_rubric_level" {
  defaults = {
    alias = "mock-rubric-level-alias"
    name  = "mock-rubric-level-name"
    # filter is not set here because its fields are not computed
    # id intentionally omitted - will be assigned a random string
    # index intentionally omitted - defaults to number type 0
  }
}
