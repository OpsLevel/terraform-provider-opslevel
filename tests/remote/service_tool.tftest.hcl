variables {
  resource_name = "opslevel_service_tool"

  # required fields
  category = "observability"
  name     = "TF Test Service Tool"
  url      = "https://example.com"

  # optional fields
  environment   = "production"
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

run "resource_service_tool_create_with_service_id" {

  variables {
    category      = var.category
    environment   = var.environment
    name          = var.name
    service       = run.from_service_module.first_service.id
    service_alias = null
    url           = var.url
  }

  module {
    source = "./service_tool"
  }

  assert {
    condition = alltrue([
      can(opslevel_service_tool.test.category),
      can(opslevel_service_tool.test.environment),
      can(opslevel_service_tool.test.id),
      can(opslevel_service_tool.test.name),
      can(opslevel_service_tool.test.service),
      can(opslevel_service_tool.test.service_alias),
      can(opslevel_service_tool.test.url),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_tool.test.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_service_tool.test.category,
    )
  }

  assert {
    condition = opslevel_service_tool.test.environment == var.environment
    error_message = format(
      "expected '%v' but got '%v'",
      var.environment,
      opslevel_service_tool.test.environment,
    )
  }

  assert {
    condition     = startswith(opslevel_service_tool.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_service_tool.test.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_service_tool.test.name,
    )
  }

  assert {
    condition = opslevel_service_tool.test.service == var.service
    error_message = format(
      "expected '%v' but got '%v'",
      var.service,
      opslevel_service_tool.test.service,
    )
  }

  assert {
    condition     = opslevel_service_tool.test.service_alias == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_tool.test.url == var.url
    error_message = format(
      "expected '%v' but got '%v'",
      var.url,
      opslevel_service_tool.test.url,
    )
  }

}

run "resource_service_tool_create_with_service_alias" {

  variables {
    category      = var.category
    environment   = var.environment
    name          = var.name
    service       = null
    service_alias = run.from_service_module.first_service.aliases[0]
    url           = var.url
  }

  module {
    source = "./service_tool"
  }

  assert {
    condition = opslevel_service_tool.test.category == var.category
    error_message = format(
      "expected '%v' but got '%v'",
      var.category,
      opslevel_service_tool.test.category,
    )
  }

  assert {
    condition = opslevel_service_tool.test.environment == var.environment
    error_message = format(
      "expected '%v' but got '%v'",
      var.environment,
      opslevel_service_tool.test.environment,
    )
  }

  assert {
    condition     = startswith(opslevel_service_tool.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition     = opslevel_service_tool.test.service == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_service_tool.test.service_alias == var.service_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.service_alias,
      opslevel_service_tool.test.service_alias,
    )
  }

  assert {
    condition = opslevel_service_tool.test.url == var.url
    error_message = format(
      "expected '%v' but got '%v'",
      var.url,
      opslevel_service_tool.test.url,
    )
  }

}
