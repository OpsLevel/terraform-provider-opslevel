variables {
  resource_name = "opslevel_system"

  # required fields
  name = "TF Test System"

  # optional fields
  description = "System description"
  domain      = null # sourced from module
  note        = "System note"
  owner       = null # sourced from module
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_domains.all,
      data.opslevel_teams.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_system_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    domain = run.from_data_module.first_domain.id
    owner  = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/system"
  }

  assert {
    condition = alltrue([
      can(opslevel_system.this.aliases),
      can(opslevel_system.this.description),
      can(opslevel_system.this.domain),
      can(opslevel_system.this.id),
      can(opslevel_system.this.name),
      can(opslevel_system.this.note),
      can(opslevel_system.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_system.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_system.this.description,
    )
  }

  assert {
    condition = opslevel_system.this.domain == var.domain
    error_message = format(
      "expected '%v' but got '%v'",
      var.domain,
      opslevel_system.this.domain,
    )
  }

  assert {
    condition     = startswith(opslevel_system.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_system.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_system.this.name,
    )
  }

  assert {
    condition = opslevel_system.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_system.this.note,
    )
  }

  assert {
    condition = opslevel_system.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_system.this.owner,
    )
  }

}

run "resource_sytem_unset_optional_fields" {

  variables {
    # required fields from file scoped variables block
    description = null
    domain      = null
    note        = null
    owner       = null
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/system"
  }

  assert {
    condition     = opslevel_system.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.this.domain == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.this.note == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "delete_system_outside_of_terraform" {

  variables {
    command = "delete system ${run.resource_system_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_system_create_with_required_fields" {

  variables {
    # required fields from file scoped variables block
    description = null
    domain      = null
    note        = null
    owner       = null
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/system"
  }

  assert {
    condition = run.resource_system_create_with_all_fields.this.id != opslevel_system.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_system_create_with_all_fields.this.id,
      opslevel_system.this.id,
    )
  }

  assert {
    condition     = opslevel_system.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.this.domain == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.this.note == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.this.owner == null
    error_message = var.error_expected_null_field
  }

}

run "resource_system_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    domain = run.from_data_module.first_domain.id
    owner  = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/system"
  }

  assert {
    condition = opslevel_system.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_system.this.description,
    )
  }

  assert {
    condition = opslevel_system.this.domain == var.domain
    error_message = format(
      "expected '%v' but got '%v'",
      var.domain,
      opslevel_system.this.domain,
    )
  }

  assert {
    condition = opslevel_system.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_system.this.name,
    )
  }

  assert {
    condition = opslevel_system.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_system.this.note,
    )
  }

  assert {
    condition = opslevel_system.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_system.this.owner,
    )
  }

}
