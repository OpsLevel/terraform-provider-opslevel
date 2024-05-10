run "datasource_filters_all" {

  variables {
    datasource_type = "opslevel_filters"
  }

  module {
    source = "./filter"
  }

  assert {
    condition     = can(data.opslevel_filters.all.filters)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_filters.all.filters) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

}

run "datasource_filter_first" {

  variables {
    datasource_type = "opslevel_filter"
  }

  module {
    source = "./filter"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_filter.first_filter_by_id.id),
      can(data.opslevel_filter.first_filter_by_id.name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_filter.first_filter_by_id.id == data.opslevel_filters.all.filters[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_filter.first_filter_by_name.name == data.opslevel_filters.all.filters[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.datasource_type)
  }

}
