variables {
  resource_name = "opslevel_service_relationship"

  # required fields
  service = null # sourced from module
  system  = null # sourced from module
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_services.all,
      data.opslevel_systems.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_service_relationship_create_with_ids" {

  variables {
    service = run.from_data_module.first_service.id
    system  = run.from_data_module.first_system.id
  }

  module {
    source = "./opslevel_modules/modules/service/relationship"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_relationship.this.id),
      can(opslevel_service_relationship.this.service),
      can(opslevel_service_relationship.this.system),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_service_relationship.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_relationship.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_relationship.this.service,
    )
  }

  assert {
    condition = opslevel_service_relationship.this.system == var.system
    error_message = format(
      "expected '%v' but got '%v'",
      var.system,
      opslevel_service_relationship.this.system,
    )
  }

}

run "resource_service_create_minimal" {

  variables {
    name = "Test minimal service"
  }

  module {
    source = "./opslevel_modules/modules/service"
  }
}

run "resource_service_relationship_update_does_force_recreate" {

  variables {
    service = run.resource_service_create_minimal.this.id
    system  = run.from_data_module.first_system.id
  }

  module {
    source = "./opslevel_modules/modules/service/relationship"
  }

  assert {
    condition = run.resource_service_relationship_create_with_ids.this.id != opslevel_service_relationship.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_service_relationship_create_with_ids.this.id,
      opslevel_service_relationship.this.id,
    )
  }

  assert {
    condition = opslevel_service_relationship.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_relationship.this.service,
    )
  }

  assert {
    condition = opslevel_service_relationship.this.system == var.system
    error_message = format(
      "expected '%v' but got '%v'",
      var.system,
      opslevel_service_relationship.this.system,
    )
  }
}

run "resource_system_create_minimal" {

  variables {
    name = "Test minimal system"
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/system"
  }
}

run "resource_service_relationship_update_system" {

  variables {
    service = run.resource_service_create_minimal.this.id
    system  = run.resource_system_create_minimal.this.id
  }

  module {
    source = "./opslevel_modules/modules/service/relationship"
  }

  assert {
    condition = opslevel_service_relationship.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_relationship.this.service,
    )
  }

  assert {
    condition = opslevel_service_relationship.this.system == var.system
    error_message = format(
      "expected '%v' but got '%v'",
      var.system,
      opslevel_service_relationship.this.system,
    )
  }
}
