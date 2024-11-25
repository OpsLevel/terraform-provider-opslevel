variables {
  resource_name = "opslevel_service_dependency"

  # required fields
  depends_upon = null
  service      = null

  # optional fields
  note = "Testing service dependency resource in Terraform"
}

run "from_service_module" {
  command = plan

  module {
    source = "./data/service"
  }
}

run "resource_service_dependency_create_with_service_id" {

  variables {
    depends_upon = run.from_service_module.all.services[0].id
    service      = run.from_service_module.all.services[1].id
    note         = var.note
  }

  module {
    source = "./opslevel_modules/modules/service/dependency"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_dependency.this.depends_upon),
      can(opslevel_service_dependency.this.id),
      can(opslevel_service_dependency.this.note),
      can(opslevel_service_dependency.this.service),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_dependency.this.depends_upon == var.depends_upon
    error_message = format(
      "expected '%v' but got '%v'",
      var.depends_upon,
      opslevel_service_dependency.this.depends_upon,
    )
  }

  assert {
    condition     = startswith(opslevel_service_dependency.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_dependency.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_service_dependency.this.note,
    )
  }

  assert {
    condition = opslevel_service_dependency.this.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_dependency.this.service,
    )
  }

}

run "resource_service_dependency_update_does_force_recreate" {

  variables {
    depends_upon = run.from_service_module.all.services[0].id
    service      = run.from_service_module.all.services[1].id
    note         = null
  }

  assert {
    condition = run.resource_service_dependency_create_with_service_id.this.id != opslevel_service_dependency.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_service_dependency_create_with_service_id.this.id,
      opslevel_service_dependency.this.id,
    )
  }

  module {
    source = "./opslevel_modules/modules/service/dependency"
  }

  assert {
    condition     = opslevel_service_dependency.this.note == null
    error_message = var.error_expected_null_field
  }

}
