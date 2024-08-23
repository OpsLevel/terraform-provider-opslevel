variables {
  resource_name = "opslevel_service_tag"

  # required fields
  key   = "test-service-tag-key"
  value = "test-service-tag-key"

  # optional fields
  service       = null
  service_alias = null
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

run "resource_service_tag_create_with_service_id" {

  variables {
    key           = var.key
    service       = run.from_service_module.first_service.id
    service_alias = null
    value         = var.value
  }

  module {
    source = "./service_tag"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_tag.test.key),
      can(opslevel_service_tag.test.id),
      can(opslevel_service_tag.test.service),
      can(opslevel_service_tag.test.service_alias),
      can(opslevel_service_tag.test.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_service_tag.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_tag.test.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_service_tag.test.key,
    )
  }

  assert {
    condition = opslevel_service_tag.test.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_tag.test.service,
    )
  }

  assert {
    condition     = opslevel_service_tag.test.service_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_tag.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_service_tag.test.value,
    )
  }

}

run "resource_service_tag_create_with_service_alias" {

  variables {
    key           = var.key
    service       = null
    service_alias = run.from_service_module.first_service.aliases[0]
    value         = var.value
  }

  module {
    source = "./service_tag"
  }

  assert {
    condition = opslevel_service_tag.test.key == var.key
    error_message = format(
      "expected '%v' but got '%v'",
      var.key,
      opslevel_service_tag.test.key,
    )
  }

  assert {
    condition     = opslevel_service_tag.test.service == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_tag.test.service_alias == var.service_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.service_alias,
      opslevel_service_tag.test.service_alias,
    )
  }

  assert {
    condition = opslevel_service_tag.test.value == var.value
    error_message = format(
      "expected '%v' but got '%v'",
      var.value,
      opslevel_service_tag.test.value,
    )
  }

}
