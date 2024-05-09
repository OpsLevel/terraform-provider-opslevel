run "datasource_integrations_all" {

  variables {
    datasource_type = "opslevel_integrations"
  }

  assert {
    condition     = can(data.opslevel_integrations.all.integrations)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_integrations.all.integrations) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_integration_first" {

  variables {
    datasource_type = "opslevel_integration"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_integration.first_integration_by_id.id),
      can(data.opslevel_integration.first_integration_by_id.name),
    ])
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_integration.first_integration_by_id.id == data.opslevel_integrations.all.integrations[0].id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_integration.first_integration_by_name.name == data.opslevel_integrations.all.integrations[0].name
    error_message = replace(var.wrong_name_error, "TYPE", var.datasource_type)
  }

}
