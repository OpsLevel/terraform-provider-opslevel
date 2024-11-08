variables {
  resource_name = "opslevel_secret"

  # required fields
  alias = "tf_test_secret"
  owner = null # sourced from module
  value = "TFTestSecretValue"

  # optional fields - none
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_secret_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    owner = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/secret"
  }

  assert {
    condition = alltrue([
      can(opslevel_secret.this.alias),
      can(opslevel_secret.this.created_at),
      can(opslevel_secret.this.id),
      can(opslevel_secret.this.owner),
      can(opslevel_secret.this.updated_at),
      can(opslevel_secret.this.value),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_secret.this.alias == var.alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.alias,
      opslevel_secret.this.alias,
    )
  }

  assert {
    condition = opslevel_secret.this.created_at != null && opslevel_secret.this.updated_at != null
    error_message = format(
      "expected '%v' but to match '%v'",
      opslevel_secret.this.created_at,
      opslevel_secret.this.updated_at
    )
  }

  assert {
    condition = opslevel_secret.this.created_at == opslevel_secret.this.updated_at
    error_message = format(
      "expected '%v' but to match '%v'",
      opslevel_secret.this.created_at,
      opslevel_secret.this.updated_at
    )
  }

  assert {
    condition     = startswith(opslevel_secret.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_secret.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_secret.this.owner,
    )
  }

  assert {
    condition     = opslevel_secret.this.value == var.value
    error_message = "expected different secret value, not printing sensitive value"
  }

}

run "resource_secret_replaced_by_alias_update" {

  variables {
    # other fields from file scoped variables block
    alias = "tf_test_secret_updated"
    owner = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/secret"
  }

  assert {
    condition = opslevel_secret.this.alias == var.alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.alias,
      opslevel_secret.this.alias,
    )
  }

  assert {
    condition = run.resource_secret_create_with_all_fields.this.id != opslevel_secret.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_secret_create_with_all_fields.this.id,
      opslevel_secret.this.id,
    )
  }
}

run "delete_secret_outside_of_terraform" {

  variables {
    command = "delete secret ${run.resource_secret_replaced_by_alias_update.this.id}"
  }

  module {
    source = "./cli"
  }

}

run "resource_secret_recreated_with_original_data" {

  variables {
    owner = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/secret"
  }

  assert {
    condition = opslevel_secret.this.alias == run.resource_secret_create_with_all_fields.this.alias
    error_message = format(
      "expected '%v' but got '%v'",
      run.resource_secret_create_with_all_fields.this.alias,
      opslevel_secret.this.alias,
    )
  }

  assert {
    condition = opslevel_secret.this.owner == run.resource_secret_create_with_all_fields.this.owner
    error_message = format(
      "expected '%v' but got '%v'",
      run.resource_secret_create_with_all_fields.this.owner,
      opslevel_secret.this.owner,
    )
  }

  assert {
    condition     = opslevel_secret.this.value == run.resource_secret_create_with_all_fields.this.value
    error_message = "expected different secret value, not printing sensitive value"
  }

}
