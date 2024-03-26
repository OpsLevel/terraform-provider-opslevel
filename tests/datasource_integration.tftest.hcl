mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_integration_filter_by_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_integration.id_filter.id != ""
    error_message = "id in opslevel_integration mock was not set"
  }

  assert {
    condition     = data.opslevel_integration.id_filter.name == "My Mock Integration"
    error_message = "wrong name in opslevel_integration"
  }

  assert {
    condition     = data.opslevel_integration.id_filter.filter.field == "id"
    error_message = "filter field should be id"
  }

  assert {
    condition     = data.opslevel_integration.id_filter.filter.value == "Z2lkOi8vb3BzbGV2ZWwvSW50ZWdyYXRpb25zOjpTbGFja0ludGVncmF0aW9uLzI2"
    error_message = "filter value for opslevel_integration.id_filter should be Z2lkOi8vb3BzbGV2ZWwvSW50ZWdyYXRpb25zOjpTbGFja0ludGVncmF0aW9uLzI2"
  }

}

run "datasource_integration_filter_by_name" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_integration.name_filter.filter.field == "name"
    error_message = "filter field should be name"
  }

  assert {
    condition     = data.opslevel_integration.name_filter.filter.value == "My Integration I Got By Filtering Name"
    error_message = "filter value for opslevel_integration.name_filter should be 'My Integration I Got By Filtering Name'"
  }

}

