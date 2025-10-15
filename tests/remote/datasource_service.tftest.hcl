variables {
  api_token = ""
}

run "datasource_service_reads_fields" {
  command = plan

  variables {
    api_token = var.api_token
    name      = "Test Service for Data Source"
  }

  module {
    source = "./service"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.id != null
    error_message = "expected id to be set for opslevel_service data source"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.name != null
    error_message = "expected name to be set for opslevel_service data source"
  }

  assert {
    condition     = data.opslevel_service.first_service_by_id.description != null
    error_message = "expected description to be set for opslevel_service data source"
  }

  assert {
    condition     = can(data.opslevel_service.first_service_by_id.owner)
    error_message = "expected owner field to be accessible for opslevel_service data source"
  }

  assert {
    condition     = data.opslevel_service.last_service.id != null
    error_message = "expected id to be set for opslevel_service data source"
  }

  assert {
    condition     = data.opslevel_service.last_service.name != null
    error_message = "expected name to be set for opslevel_service data source"
  }
}

