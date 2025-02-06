variables {
  resource_name = "opslevel_rubric_level"

  # required fields
  name = "TF Rubric Level"

  # optional fields
  description = "TF Rubric Level description"
  index       = null
}

run "resource_rubric_level_create_with_all_fields" {

  variables {
    # other fields from file scoped variables block
    description = var.description
  }

  module {
    source = "./opslevel_modules/modules/rubric_level"
  }

  assert {
    condition = alltrue([
      can(opslevel_rubric_level.this.description),
      can(opslevel_rubric_level.this.id),
      can(opslevel_rubric_level.this.index),
      can(opslevel_rubric_level.this.name),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_rubric_level.this.description == var.description
    error_message = format(
      "expected '%v' but got '%v'",
      var.description,
      opslevel_rubric_level.this.description,
    )
  }

  assert {
    condition     = startswith(opslevel_rubric_level.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition     = opslevel_rubric_level.this.index >= 0
    error_message = "wrong index for ${var.resource_name}"
  }

  assert {
    condition = opslevel_rubric_level.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_rubric_level.this.name,
    )
  }

}

run "resource_rubric_level_update_unset_optional_fields" {

  variables {
    # other fields from file scoped variables block
    description = null
    index       = run.resource_rubric_level_create_with_all_fields.this.index
  }

  module {
    source = "./opslevel_modules/modules/rubric_level/data"
  }

  assert {
    condition     = opslevel_rubric_level.this.description == null
    error_message = var.error_expected_null_field
  }

}

run "delete_rubric_level_outside_of_terraform" {

  variables {
    command = "delete level ${run.resource_rubric_level_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_rubric_level_create_with_required_fields" {

  variables {
    # other fields from file scoped variables block
    description = null
  }

  module {
    source = "./opslevel_modules/modules/rubric_level/data"
  }

  assert {
    condition = run.resource_rubric_level_create_with_all_fields.this.id != opslevel_rubric_level.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_rubric_level_create_with_all_fields.this.id,
      opslevel_rubric_level.this.id,
    )
  }

  assert {
    condition     = opslevel_rubric_level.this.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = run.resource_rubric_level_create_with_all_fields.this.index == opslevel_rubric_level.this.index
    error_message = format(
      "expected '%v' but got '%v'",
      run.resource_rubric_level_create_with_all_fields.this.index,
      opslevel_rubric_level.this.index,
    )
  }

  assert {
    condition     = opslevel_rubric_level.this.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.resource_name)
  }

}

run "resource_rubric_level_set_all_fields" {

  variables {
    # all fields being updated from file scoped variables block
  }

  module {
    source = "./opslevel_modules/modules/rubric_level/data"
  }

  assert {
    condition     = opslevel_rubric_level.this.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.resource_name)
  }

}
