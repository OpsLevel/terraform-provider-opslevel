# TODO: test predicate_list
variables {
  filter_one  = "opslevel_filter"
  filters_all = "opslevel_filters"

  # required fields
  name = "TF Test Filter"

  # optional fields
  connective = "and"
  # predicate_list = null
}

run "resource_filter_create_with_all_fields" {

  variables {
    connective = var.connective
    name       = var.name
    # predicate_list = var.predicate_list
  }

  module {
    source = "./filter"
  }

  assert {
    condition = alltrue([
      can(opslevel_filter.test.connective),
      can(opslevel_filter.test.id),
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

#run "resource_filter_update_unset_optional_fields" {
#
#  variables {
#    predicate_list = null
#  }
#
#  module {
#    source = "./filter"
#  }
#
#  assert {
#    condition     = opslevel_filter.test.predicate_list == null
#    error_message = var.error_expected_null_field
#  }
#
#}

run "resource_filter_update_set_all_fields" {

  variables {
    connective = var.connective == "and" ? "or" : "and"
    name       = "${var.name} updated"
    # predicate_list = var.predicate_list
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

  #assert {
  #  condition     = opslevel_filter.test.predicate_list == var.predicate_list
  #  error_message = "wrong predicate_list of opslevel_filter resource"
  #}

}

run "datasource_filters_list_all" {

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
