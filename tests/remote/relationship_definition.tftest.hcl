variables {
  resource_name = "opslevel_relationship_definition"

  # required fields
  name                 = "TF Test Relationship Definition"
  alias                = "tf_test_relationship_definition"
  component_type_name  = "TF Test Relationship Definition Component Type"
  component_type_alias = "tf_test_relationship_definition_component_type"
}

run "resource_relationship_definition_create_with_component_type_alias" {

  variables {
    component_type_reference = "alias"
  }

  module {
    source = "./relationship_definition"
  }

  assert {
    condition     = startswith(opslevel_relationship_definition.this.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.resource_name)
  }

  assert {
    condition = opslevel_relationship_definition.this.component_type == var.component_type_alias
    error_message = format(
      "expected '%v' but got '%v'",
      var.component_type_alias,
      opslevel_relationship_definition.this.component_type,
    )
  }

}

run "resource_relationship_definition_switch_to_component_type_id" {

  variables {
    component_type_reference = "id"
  }

  module {
    source = "./relationship_definition"
  }

  assert {
    condition = run.resource_relationship_definition_create_with_component_type_alias.this.id == opslevel_relationship_definition.this.id
    error_message = format(
      "expected id '%v' to be unchanged after naming the same component type by its ID but got '%v'",
      run.resource_relationship_definition_create_with_component_type_alias.this.id,
      opslevel_relationship_definition.this.id,
    )
  }

  assert {
    condition = opslevel_relationship_definition.this.component_type == opslevel_component_type.this.id
    error_message = format(
      "expected '%v' but got '%v'",
      opslevel_component_type.this.id,
      opslevel_relationship_definition.this.component_type,
    )
  }

}

run "resource_relationship_definition_switch_to_component_type_alias_in_another_case" {

  variables {
    component_type_reference = "upper_alias"
  }

  module {
    source = "./relationship_definition"
  }

  assert {
    condition = run.resource_relationship_definition_create_with_component_type_alias.this.id == opslevel_relationship_definition.this.id
    error_message = format(
      "expected id '%v' to be unchanged after naming the same component type by its alias in another case but got '%v'",
      run.resource_relationship_definition_create_with_component_type_alias.this.id,
      opslevel_relationship_definition.this.id,
    )
  }

  assert {
    condition = opslevel_relationship_definition.this.component_type == upper(var.component_type_alias)
    error_message = format(
      "expected '%v' but got '%v'",
      upper(var.component_type_alias),
      opslevel_relationship_definition.this.component_type,
    )
  }

}
