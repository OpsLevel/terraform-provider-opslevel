mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_tier_mocked_fields" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_tier.mock_tier.alias == "mock-tier-alias"
    error_message = "wrong alias in opslevel_tier mock"
  }

  assert {
    condition     = data.opslevel_tier.mock_tier.id != ""
    error_message = "id in opslevel_tier mock was not set"
  }

  assert {
    condition     = data.opslevel_tier.mock_tier.index == 0
    error_message = "index in opslevel_tier should default to int type 0"
  }

  assert {
    condition     = data.opslevel_tier.mock_tier.name == "mock-tier-name"
    error_message = "wrong name in opslevel_tier"
  }
}

run "datasource_tier_filter_by_alias" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_tier.alias_filter.filter.field == "alias"
    error_message = "filter field for opslevel_tier.alias_filter should be alias"
  }

  assert {
    condition     = data.opslevel_tier.alias_filter.filter.value == "alias-value"
    error_message = "filter value for opslevel_tier.alias_filter should be alias-value"
  }

}

run "datasource_tier_filter_by_id" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_tier.id_filter.filter.field == "id"
    error_message = "filter field should be id"
  }

  assert {
    condition     = data.opslevel_tier.id_filter.filter.value == "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
    error_message = "filter value for opslevel_tier.id_filter should be Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"
  }

}

run "datasource_tier_filter_by_index" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_tier.index_filter.filter.field == "index"
    error_message = "filter field should be id"
  }

  assert {
    condition     = tonumber(data.opslevel_tier.index_filter.filter.value) == 123
    error_message = "filter value for opslevel_tier.index_filter should be 123"
  }

}
run "datasource_tier_filter_by_name" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_tier.name_filter.filter.field == "name"
    error_message = "filter field should be name"
  }

  assert {
    condition     = data.opslevel_tier.name_filter.filter.value == "name-value"
    error_message = "filter value for opslevel_tier.name_filter should be name-value"
  }

}
