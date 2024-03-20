mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_rubric_level_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_level.id_filter.alias == "mock-rubric-level-alias"
    error_message = "wrong alias in opslevel_rubric_level mock"
  }

  assert {
    condition     = data.opslevel_rubric_level.id_filter.id != ""
    error_message = "id in opslevel_rubric_level mock was not set"
  }

  assert {
    condition     = data.opslevel_rubric_level.id_filter.index == 321
    error_message = "wrong index in opslevel_rubric_level"
  }

  assert {
    condition     = data.opslevel_rubric_level.id_filter.name == "mock-rubric-level-name"
    error_message = "wrong name in opslevel_rubric_level"
  }
}

run "datasource_rubric_level_filter_by_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_level.alias_filter.filter.field == "alias"
    error_message = "filter field for opslevel_rubric_level.alias_filter should be alias"
  }

  assert {
    condition     = data.opslevel_rubric_level.alias_filter.filter.value == "alias-value"
    error_message = "filter value for opslevel_rubric_level.alias_filter should be alias-value"
  }

}

run "datasource_rubric_level_filter_by_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_level.id_filter.filter.field == "id"
    error_message = "filter field should be id"
  }

  assert {
    condition     = data.opslevel_rubric_level.id_filter.filter.value == "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
    error_message = "filter value for opslevel_rubric_level.id_filter should be Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }

}

run "datasource_rubric_level_filter_by_index" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_level.index_filter.filter.field == "index"
    error_message = "filter field should be id"
  }

  assert {
    condition     = tonumber(data.opslevel_rubric_level.index_filter.filter.value) == 123
    error_message = "filter value for opslevel_rubric_level.index_filter should be 123"
  }

}

run "datasource_rubric_level_filter_by_name" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_rubric_level.name_filter.filter.field == "name"
    error_message = "filter field should be name"
  }

  assert {
    condition     = data.opslevel_rubric_level.name_filter.filter.value == "name-value"
    error_message = "filter value for opslevel_rubric_level.name_filter should be name-value"
  }

}
