mock_provider "opslevel" {
  alias = "fake"
  source = "."
}

run "use_real_provider" {
}

run "use_mocked_provider" {
  providers = {
    opslevel = opslevel.fake
  }
}
