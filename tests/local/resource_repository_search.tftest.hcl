mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_repository_search" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_repository_search.example.name == "foo"
    error_message = "wrong value name for opslevel_check_repository_search.example"
  }

  assert {
    condition     = opslevel_check_repository_search.example.file_extensions == toset(["sbt", "py"])
    error_message = "wrong value for file_extensions in opslevel_check_repository_search.example"
  }

  assert {
    condition = opslevel_check_repository_search.example.file_contents_predicate == {
      type  = "contains"
      value = "postgres"
    }
    error_message = "wrong value for file_contents_predicate in opslevel_check_repository_search.example"
  }
}
