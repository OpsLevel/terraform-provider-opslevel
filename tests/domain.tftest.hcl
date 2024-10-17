variables {
  resource_name = "opslevel_domain"

  # required fields
  name = "TF Test Domain"

  # optional fields
  description = "Domain description"
  note        = "Domain note"
  owner       = null # sourced from module
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

run "resource_domain_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    description = var.description
    owner       = run.from_data_module.first_team.id
    note        = var.note
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/domain"
  }

  assert {
    condition = alltrue([
      can(opslevel_domain.this.aliases),
      can(opslevel_domain.this.description),
      can(opslevel_domain.this.id),
      can(opslevel_domain.this.name),
      can(opslevel_domain.this.note),
      can(opslevel_domain.this.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_domain.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_domain.this.description,
    )
  }

  assert {
    condition     = startswith(opslevel_domain.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_domain.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_domain.this.owner,
    )
  }

  assert {
    condition = opslevel_domain.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_domain.this.name,
    )
  }

  assert {
    condition = opslevel_domain.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_domain.this.note,
    )
  }

}

run "resource_domain_unset_optional_fields" {

  variables {
    # required fields from file scoped variables block
    description = null
    note        = null
    owner       = null
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/domain"
  }

  assert {
    condition     = opslevel_domain.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_domain.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_domain.this.note == null
    error_message = var.error_expected_null_field
  }

}

run "delete_domain_outside_of_terraform" {

  variables {
    resource_id   = run.resource_domain_create_with_all_fields.this.id
    resource_type = "domain"
  }

  module {
    source = "./provisioner"
  }
}

run "resource_domain_create_with_required_fields" {

  variables {
    # required fields from file scoped variables block
    description = null
    owner       = null
    note        = null
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/domain"
  }

  assert {
    condition = run.resource_domain_create_with_all_fields.this.id != opslevel_domain.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_domain_create_with_all_fields.this.id,
      opslevel_domain.this.id,
    )
  }

  assert {
    condition     = opslevel_domain.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_domain.this.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_domain.this.note == null
    error_message = var.error_expected_null_field
  }

}

run "resource_domain_set_all_fields" {

  variables {
    # other fields from file scoped variables block
    owner = run.from_data_module.first_team.id
  }

  module {
    source = "./opslevel_modules/modules/hierarchy/domain"
  }

  assert {
    condition = opslevel_domain.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_domain.this.description,
    )
  }

  assert {
    condition = opslevel_domain.this.owner == var.owner
    error_message = format(
      "expected '%v' but got '%v'",
      var.owner,
      opslevel_domain.this.owner,
    )
  }

  assert {
    condition = opslevel_domain.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_domain.this.name,
    )
  }

  assert {
    condition = opslevel_domain.this.note == var.note
    error_message = format(
      "expected '%v' but got '%v'",
      var.note,
      opslevel_domain.this.note,
    )
  }


}
