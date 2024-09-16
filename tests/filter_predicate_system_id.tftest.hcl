variables {
  name          = "TF Test Filter with system_id predicate"
  predicate_key = "system_id"
}

run "from_data_module" {
  command = plan
  plan_options {
    target = [
      data.opslevel_systems.all
    ]
  }

  module {
    source = "./data"
  }
}

run "resource_filter_with_system_id_predicate_equals" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_equals : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = null,
        value    = run.from_data_module.first_system.id
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
      contains(var.predicate_types_equals, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_equals, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_equals,
      tolist([opslevel_filter.this.predicate[0].type, opslevel_filter.this.predicate[1].type])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key_data == null,
      opslevel_filter.this.predicate[1].key_data == null
    ])
    error_message = var.error_expected_null_field
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == run.from_data_module.first_system.id,
      opslevel_filter.this.predicate[1].value == run.from_data_module.first_system.id
    ])
    error_message = format(
      "expected predicate values to all be '%v' got values '%v'",
      run.from_data_module.first_system.id,
      distinct([opslevel_filter.this.predicate[0].value, opslevel_filter.this.predicate[1].value])
    )
  }

}

run "resource_filter_with_system_id_predicate_exists" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_exists : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = null,
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
      opslevel_filter.this.predicate[0].key_data == null,
      opslevel_filter.this.predicate[1].key_data == null
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
