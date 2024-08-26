variables {
  integration_one  = "opslevel_integration"
  integrations_all = "opslevel_integrations"
}

run "resource_integration_create_with_all_fields" {

  variables {
    external_id             = var.external_id
  }

  module {
    source = "./integration"
  }
}

run "datasource_integrations_all" {

  module {
    source = "./integration"
  }

  assert {
    condition     = can(data.opslevel_integrations.all.integrations)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.integrations_all)
  }

  assert {
    condition     = length(data.opslevel_integrations.all.integrations) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.integrations_all)
  }

}

run "datasource_integration_first" {

  module {
    source = "./integration"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_integration.first_integration_by_id.id),
      can(data.opslevel_integration.first_integration_by_id.name),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.integration_one)
  }

  assert {
    condition     = data.opslevel_integration.first_integration_by_id.id == data.opslevel_integrations.all.integrations[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.integration_one)
  }

  assert {
    condition     = data.opslevel_integration.first_integration_by_name.name == data.opslevel_integrations.all.integrations[0].name
    error_message = replace(var.error_wrong_name, "TYPE", var.integration_one)
  }

}
