variables {
  name               = "TF Test Filter with properties predicate"
  jq_expression      = ".[] | select(.name == fancy)"
  predicate_key      = "properties"
  predicate_key_data = "property_definition_id"
  predicate_type     = "satisfies_jq_expression"
}

run "resource_filter_with_properties_predicate_exists" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_exists : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = var.predicate_key_data,
        value    = null
      }
    ]
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key == var.predicate_key,
      opslevel_filter.this.predicate[1].key == var.predicate_key
    ])
    error_message = format(
      "expected predicate keys to all be '%v' got keys '%v'",
      var.predicate_key,
      [opslevel_filter.this.predicate[0].key, opslevel_filter.this.predicate[1].key]
    )
  }

  assert {
    condition = alltrue([
      contains(var.predicate_types_exists, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_exists, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_exists,
      tolist([opslevel_filter.this.predicate[0].type, opslevel_filter.this.predicate[1].type])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key_data == var.predicate_key_data,
      opslevel_filter.this.predicate[1].key_data == var.predicate_key_data
    ])
    error_message = var.error_expected_null_field
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == null,
      opslevel_filter.this.predicate[1].value == null
    ])
    error_message = var.error_expected_null_field
  }

}

run "resource_filter_with_properties_predicate_satisfies_jq_expression" {

  variables {
    predicates = [
      {
        key      = var.predicate_key
        type     = var.predicate_type,
        key_data = var.predicate_key_data,
        value    = var.jq_expression
      }
    ]
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }

  assert {
    condition = opslevel_filter.this.predicate[0].key == var.predicate_key
    error_message = format(
      "expected predicate keys to all be '%v' got keys '%v'",
      var.predicate_key,
      opslevel_filter.this.predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.this.predicate[0].type == var.predicate_type
    error_message = format(
      "expected predicate types to be '%v' got '%v'",
      var.predicate_type,
      opslevel_filter.this.predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.this.predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate types to be '%v' got '%v'",
      var.predicate_key_data,
      opslevel_filter.this.predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.this.predicate[0].value == var.jq_expression
    error_message = format(
      "expected predicate types to be '%v' got '%v'",
      var.jq_expression,
      opslevel_filter.this.predicate[0].value
    )
  }

}
