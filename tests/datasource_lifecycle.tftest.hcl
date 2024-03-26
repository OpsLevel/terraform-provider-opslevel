mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_lifecycle_filter_by_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_lifecycle.id_filter.alias == "generally_available"
    error_message = "wrong alias in opslevel_lifecycle mock"
  }

  assert {
    condition     = data.opslevel_lifecycle.id_filter.id != ""
    error_message = "id in opslevel_lifecycle mock was not set"
  }

  assert {
    condition     = data.opslevel_lifecycle.id_filter.name == "Generally Available"
    error_message = "wrong name in opslevel_lifecycle"
  }

  assert {
    condition     = data.opslevel_lifecycle.id_filter.filter.field == "id"
    error_message = "filter field should be id"
  }

  assert {
    condition     = data.opslevel_lifecycle.id_filter.filter.value == "Z2lkOi8vb3BzbGV2ZWwvTGlmZWN5Y2xlLzQ"
    error_message = "filter value for opslevel_lifecycle.id_filter should be Z2lkOi8vb3BzbGV2ZWwvTGlmZWN5Y2xlLzQ"
  }

}

run "datasource_lifecycle_filter_by_name" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_lifecycle.name_filter.filter.field == "name"
    error_message = "filter field should be name"
  }

  assert {
    condition     = data.opslevel_lifecycle.name_filter.filter.value == "Generally Available"
    error_message = "filter value for opslevel_lifecycle.name_filter should be 'Generally Available'"
  }

}

run "datasource_lifecycle_filter_by_index" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_lifecycle.index_filter.filter.field == "index"
    error_message = "filter field should be index"
  }

  assert {
    condition     = data.opslevel_lifecycle.index_filter.filter.value == "123"
    error_message = "filter value for opslevel_lifecycle.name_filter should be '123'"
  }

}

run "datasource_lifecycle_filter_by_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_lifecycle.alias_filter.filter.field == "alias"
    error_message = "filter field should be alias"
  }

  assert {
    condition     = data.opslevel_lifecycle.alias_filter.filter.value == "generally_available"
    error_message = "filter value for opslevel_lifecycle.name_filter should be 'generally_available'"
  }

}