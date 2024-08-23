variables {
  resource_name = "opslevel_service_dependency"

  # required fields
  depends_upon = "<id or alias of service>"
  service      = null

  # optional fields
  note = "Testing service dependency resource in Terraform"
}

run "from_service_module" {
  command = plan

  variables {
    name = ""
  }

  module {
    source = "./service"
  }
}

run "resource_service_dependency_create_with_service_id" {

  variables {
    depends_upon = run.from_service_module.first_service.id
    service      = run.from_service_module.last_service.id
    note         = var.note
  }

  module {
    source = "./service_dependency"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_dependency.test.depends_upon),
      can(opslevel_service_dependency.test.id),
      can(opslevel_service_dependency.test.note),
      can(opslevel_service_dependency.test.service),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_dependency.test.depends_upon == var.depends_upon
    error_message = format(
      "expected '%v' but got '%v'",
      var.depends_upon,
      opslevel_service_dependency.test.depends_upon,
    )
  }

  assert {
    condition     = startswith(opslevel_service_dependency.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_dependency.test.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_service_dependency.test.note,
    )
  }

  assert {
    condition = opslevel_service_dependency.test.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_dependency.test.service,
    )
  }

}

run "resource_service_dependency_update_unset_optional_fields" {

  variables {
    depends_upon = run.from_service_module.first_service.id
    service      = run.from_service_module.last_service.id
    note         = null
  }

  module {
    source = "./service_dependency"
  }

  assert {
    condition     = opslevel_service_dependency.test.note == null
    error_message = var.error_expected_null_field
  }

}
