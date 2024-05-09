mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_check_tag_defined" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_check_tag_defined.example.tag_key == "environment"
    error_message = "wrong value for update_requires_comment in opslevel_check_tag_defined.example"
  }

  assert {
    condition = opslevel_check_tag_defined.example.tag_predicate == {
      type  = "contains"
      value = "dev"
    }
    error_message = "wrong update_frequency in opslevel_check_tag_defined.example"
  }
}