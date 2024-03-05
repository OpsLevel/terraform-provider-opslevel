mock_provider "opslevel" {
  alias = "fake"
}

run "use_mocked_provider" {
  providers = {
    opslevel = opslevel.fake
  }
}
