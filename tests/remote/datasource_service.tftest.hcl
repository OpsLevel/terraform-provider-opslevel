run "datasource_services_all" {

  variables {
    datasource_type = "opslevel_services"
  }

  module {
    source = "./service"
  }

  assert {
    condition     = can(data.opslevel_services.all.services)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_services.all.services) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.datasource_type)
  }

  assert {
    condition = alltrue([
      can(data.opslevel_services.all.services[0].id),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

}

run "datasource_service_first" {

  variables {
    datasource_type = "opslevel_service"
  }

  module {
    source = "./service"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_service.first_service_by_id.alias),
      can(data.opslevel_service.first_service_by_id.aliases),
      can(data.opslevel_service.first_service_by_id.api_document_path),
      can(data.opslevel_service.first_service_by_id.description),
      can(data.opslevel_service.first_service_by_id.framework),
      can(data.opslevel_service.first_service_by_id.id),
      can(data.opslevel_service.first_service_by_id.language),
      can(data.opslevel_service.first_service_by_id.lifecycle_alias),
      can(data.opslevel_service.first_service_by_id.name),
      can(data.opslevel_service.first_service_by_id.owner),
      can(data.opslevel_service.first_service_by_id.owner_id),
      can(data.opslevel_service.first_service_by_id.preferred_api_document_source),
      can(data.opslevel_service.first_service_by_id.product),
      can(data.opslevel_service.first_service_by_id.properties),
      can(data.opslevel_service.first_service_by_id.repositories),
      can(data.opslevel_service.first_service_by_id.tags),
      can(data.opslevel_service.first_service_by_id.tier_alias),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.id == data.opslevel_services.all.services[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.datasource_type)
  }

}
