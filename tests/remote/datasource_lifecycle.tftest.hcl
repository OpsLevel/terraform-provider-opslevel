run "datasource_lifecycles_all" {

  variables {
    datasource_type = "opslevel_lifecycles"
  }

  assert {
    condition     = can(data.opslevel_lifecycles.all.lifecycles)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_lifecycles.all.lifecycles) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_lifecycle_first" {

  variables {
    datasource_type = "opslevel_lifecycle"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_lifecycle.first_lifecycle_by_id.alias),
      can(data.opslevel_lifecycle.first_lifecycle_by_id.id),
      can(data.opslevel_lifecycle.first_lifecycle_by_id.index),
      can(data.opslevel_lifecycle.first_lifecycle_by_id.name),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.alias == data.opslevel_lifecycles.all.lifecycles[0].alias
    error_message = replace(var.wrong_alias_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.id == data.opslevel_lifecycles.all.lifecycles[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.index == data.opslevel_lifecycles.all.lifecycles[0].index
    error_message = replace(var.wrong_index_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_name.name == data.opslevel_lifecycles.all.lifecycles[0].name
    error_message = replace(var.wrong_name_error, "TYPE", var.datasource_type)
  }

}
