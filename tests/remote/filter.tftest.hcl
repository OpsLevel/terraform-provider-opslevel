variables {
  filter_one  = "opslevel_filter"
  filters_all = "opslevel_filters"

  # opslevel_filter fields
  connective = "optional"
  name       = "required"
}

run "datasource_filters_all" {

  variables {
    connective = null
  }

  module {
    source = "./filter"
  }

  assert {
    condition     = can(data.opslevel_filters.all.filters)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.filters_all)
  }

  assert {
    condition     = length(data.opslevel_filters.all.filters) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.filters_all)
  }

}

run "datasource_filter_get_first" {

  variables {
    connective = null
  }

  module {
    source = "./filter"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_filter.first_filter_by_id.id),
      can(data.opslevel_filter.first_filter_by_id.name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.filter_one)
  }

  assert {
    condition     = data.opslevel_filter.first_filter_by_id.id == data.opslevel_filters.all.filters[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.filter_one)
  }

  assert {
    condition     = data.opslevel_filter.first_filter_by_name.name == data.opslevel_filters.all.filters[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.filter_one)
  }

}

run "resource_filter_create_with_all_fields" {

  variables {
    connective = "and"
    name       = "TF Test Filter"
  }

  module {
    source = "./filter"
  }

  assert {
    condition = alltrue([
      can(opslevel_filter.test.connective),
      can(opslevel_filter.test.id),
      can(opslevel_filter.test.last_updated),
      can(opslevel_filter.test.name),
      can(opslevel_filter.test.predicate),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.filter_one)
  }

  assert {
    condition     = opslevel_filter.test.connective == var.connective
    error_message = "wrong connective of opslevel_filter resource"
  }

  assert {
    condition     = opslevel_filter.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.filter_one)
  }

}

run "resource_filter_update_unset_optional_fields" {

  variables {
    connective = null
    name       = "TF Test Filter only required fields set"
  }

  module {
    source = "./filter"
  }

  assert {
    condition     = opslevel_filter.test.connective == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.filter_one)
  }

}

run "resource_filter_update_set_optional_fields" {

  variables {
    connective = null
    name       = "TF Test Filter only all fields set"
  }

  module {
    source = "./filter"
  }

  assert {
    condition     = opslevel_filter.test.connective == var.connective
    error_message = "wrong connective of opslevel_filter resource"
  }

  assert {
    condition     = opslevel_filter.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.filter_one)
  }

}
