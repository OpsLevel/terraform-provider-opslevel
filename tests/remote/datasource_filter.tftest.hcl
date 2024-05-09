run "datasource_filters_all" {

  assert {
    condition     = length(data.opslevel_filters.all.filters) > 0
    error_message = "zero filters found in data.opslevel_filters"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_filters.all.filters[0].id),
      can(data.opslevel_filters.all.filters[0].name),
    ])
    error_message = "cannot set all expected filter datasource fields"
  }

}

run "datasource_filter_first" {

  assert {
    condition     = data.opslevel_filter.first_filter_by_id.id == data.opslevel_filters.all.filters[0].id
    error_message = "wrong ID on opslevel_filter"
  }

  assert {
    condition     = data.opslevel_filter.first_filter_by_name.name == data.opslevel_filters.all.filters[0].name
    error_message = "wrong name on opslevel_filter"
  }

}
