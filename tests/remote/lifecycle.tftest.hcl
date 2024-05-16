variables {
  lifecycle_one  = "opslevel_lifecycle"
  lifecycles_all = "opslevel_lifecycles"
}

run "datasource_lifecycles_all" {

  module {
    source = "./lifecycle"
  }

  assert {
    condition     = can(data.opslevel_lifecycles.all.lifecycles)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.lifecycle_one)
  }

  assert {
    condition     = length(data.opslevel_lifecycles.all.lifecycles) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.lifecycle_one)
  }

}

run "datasource_lifecycle_first" {

  module {
    source = "./lifecycle"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_lifecycle.first_lifecycle_by_id.alias),
      can(data.opslevel_lifecycle.first_lifecycle_by_id.id),
      can(data.opslevel_lifecycle.first_lifecycle_by_id.index),
      can(data.opslevel_lifecycle.first_lifecycle_by_id.name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.lifecycle_one)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.alias == data.opslevel_lifecycles.all.lifecycles[0].alias
    error_message = replace(var.error_wrong_alias, "TYPE", var.lifecycle_one)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.id == data.opslevel_lifecycles.all.lifecycles[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.lifecycle_one)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_id.index == data.opslevel_lifecycles.all.lifecycles[0].index
    error_message = replace(var.error_wrong_index, "TYPE", var.lifecycle_one)
  }

  assert {
    condition     = data.opslevel_lifecycle.first_lifecycle_by_name.name == data.opslevel_lifecycles.all.lifecycles[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.lifecycle_one)
  }

}
