mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_filter_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_filter.mock_filter.name == "mock-filter-name"
    error_message = "wrong name in opslevel_filter mock"
  }

  assert {
    condition     = data.opslevel_filter.mock_filter.id != ""
    error_message = "id in opslevel_filter mock was not set"
  }

}

run "datasource_filter_filter_by_name" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_filter.name_filter.filter.field == "name"
    error_message = "filter field for opslevel_filter.name_filter should be name"
  }

  assert {
    condition     = data.opslevel_filter.name_filter.filter.value == "name-value"
    error_message = "filter value for opslevel_filter.name_filter should be name-value"
  }

}

run "datasource_filter_filter_by_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_filter.id_filter.filter.field == "id"
    error_message = "filter field should be id"
  }

  assert {
    condition     = data.opslevel_filter.id_filter.filter.value == "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
    error_message = "filter value for opslevel_filter.id_filter should be Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }

}
