mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_rubric_category" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_rubric_category.mock_category.id != null && opslevel_rubric_category.mock_category.id != ""
    error_message = "opslevel_rubric_category id should not be empty"
  }

  assert {
    condition     = opslevel_rubric_category.mock_category.name == "Mock Category"
    error_message = "wrong name for opslevel_rubric_category.mock_category"
  }

}
