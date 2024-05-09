mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_repository_grep" {
  providers = {
    opslevel = opslevel.fake
  }
  assert {
    condition     = opslevel_check_repository_grep.example.directory_search == false
    error_message = "wrong value for directory_search in opslevel_check_repository_grep.example"
  }

  assert {
    condition = opslevel_check_repository_grep.example.file_contents_predicate == {
      type  = "contains"
      value = "**/hello.go"
    }
    error_message = "wrong value for file_contents_predicate in opslevel_check_repository_grep.example"
  }

  assert {
    condition     = opslevel_check_repository_grep.example.filepaths == tolist(["/src", "/tests"])
    error_message = "wrong value for filepaths in opslevel_check_repository_grep.example"
  }
}