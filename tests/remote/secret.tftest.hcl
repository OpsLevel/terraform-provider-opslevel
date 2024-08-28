variables {
  resource_name = "opslevel_secret"

  # required fields
  alias = "tf_test_secret"
  owner = null
  value = "TFTestSecretValue"

  # optional fields - none
}


run "from_team_module" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "resource_secret_create" {

  variables {
    alias = var.alias
    owner = run.from_team_module.first_team.id
    value = var.value
  }

  module {
    source = "./secret"
  }

  assert {
    condition = alltrue([
      can(opslevel_secret.test.alias),
      can(opslevel_secret.test.created_at),
      can(opslevel_secret.test.id),
      can(opslevel_secret.test.owner),
      can(opslevel_secret.test.updated_at),
      can(opslevel_secret.test.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_secret.test.alias == var.alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.alias,
      opslevel_secret.test.alias,
    )
  }

  assert {
    condition     = opslevel_secret.test.created_at != null && opslevel_secret.test.updated_at != null
    error_message = "expected 'created_at' to be set"
  }

  assert {
    condition     = opslevel_secret.test.created_at == opslevel_secret.test.updated_at
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = startswith(opslevel_secret.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_secret.test.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_secret.test.owner,
    )
  }

  assert {
    condition     = opslevel_secret.test.value == var.value
    error_message = "expected different secret value, not printing sensitive value"
  }

}

run "resource_secret_update" {

  variables {
    alias = var.alias
    owner = run.from_team_module.first_team.id
    value = upper(var.value)
  }

  module {
    source = "./secret"
  }

  assert {
    condition     = opslevel_secret.test.created_at != opslevel_secret.test.updated_at
    error_message = "expected 'created_at' and 'updated_at' to be different"
  }

  assert {
    condition     = opslevel_secret.test.value == upper(var.value)
    error_message = "expected different secret value, not printing sensitive value"
  }

}
