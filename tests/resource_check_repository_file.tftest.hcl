mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_repository_file" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_repository_file.example.directory_search == false
    error_message = "wrong value for directory_search in opslevel_check_repository_file.example"
  }

  assert {
    condition = opslevel_check_repository_file.example.file_contents_predicate == {
      type  = "equals"
      value = "import shim"
    }
    error_message = "wrong value for file_contents_predicate in opslevel_check_repository_file.example"
  }

  assert {
    condition     = opslevel_check_repository_file.example.filepaths == tolist(["/src", "/tests"])
    error_message = "wrong value for filepaths in opslevel_check_repository_file.example"
  }

  assert {
    condition     = opslevel_check_repository_file.example.use_absolute_root == false
    error_message = "wrong value for use_absolute_root in opslevel_check_repository_file.example"
  }
}