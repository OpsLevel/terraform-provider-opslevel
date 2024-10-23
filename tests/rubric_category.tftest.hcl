variables {
  resource_name = "opslevel_rubric_category"

  # required fields
  name = "TF Rubric Category"
}

run "resource_rubric_category_create_with_all_fields" {

  module {
    source = "./opslevel_modules/modules/rubric_category"
  }

  assert {
    condition = alltrue([
      can(opslevel_rubric_category.this.id),
      can(opslevel_rubric_category.this.name),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.resource_name)
  }

  assert {
    condition     = startswith(opslevel_rubric_category.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition     = opslevel_rubric_category.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_rubric_category.this.name,
    )
  }

}

run "delete_rubric_category_outside_of_terraform" {

  variables {
    command = "delete category ${run.resource_rubric_category_create_with_all_fields.this.id}"
  }

  module {
    source = "./cli"
  }
}

run "resource_rubric_category_recreate_when_not_found" {

  module {
    source = "./opslevel_modules/modules/rubric_category"
  }

  assert {
    condition = run.resource_rubric_category_create_with_all_fields.this.id != opslevel_rubric_category.this.id
    error_message = format(
      "expected old id '%v' to be different from new id '%v'",
      run.resource_rubric_category_create_with_all_fields.this.id,
      opslevel_rubric_category.this.id,
    )
  }

  assert {
    condition     = opslevel_rubric_category.this.name == var.name
    error_message = format(
      "expected '%v' but got '%v'",
      var.name,
      opslevel_rubric_category.this.name,
    )
  }

}
