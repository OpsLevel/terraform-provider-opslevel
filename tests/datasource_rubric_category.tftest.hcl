mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_rubric_category_filter_by_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_category.id_filter.id != ""
    error_message = "id in opslevel_rubric_category mock was not set"
  }

  assert {
    condition     = data.opslevel_rubric_category.id_filter.name == "mock-rubric-category-name"
    error_message = "wrong name in opslevel_rubric_category"
  }

  assert {
    condition     = data.opslevel_rubric_category.id_filter.filter.field == "id"
    error_message = "filter field should be id"
  }

  assert {
    condition     = data.opslevel_rubric_category.id_filter.filter.value == "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
    error_message = "filter value for opslevel_rubric_category.id_filter should be Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }

}

run "datasource_rubric_category_filter_by_name" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_category.name_filter.filter.field == "name"
    error_message = "filter field should be name"
  }

  assert {
    condition     = data.opslevel_rubric_category.name_filter.filter.value == "name-value"
    error_message = "filter value for opslevel_rubric_category.name_filter should be name-value"
  }

}

