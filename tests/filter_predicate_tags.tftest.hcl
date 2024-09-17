variables {
  name                                         = "TF Test Filter with tags predicate"
  predicate_key                                = "tags"
  predicate_key_data                           = "test_tag"
  predicate_types_satisfies_version_constraint = "satisfies_version_constraint"
  predicate_value                              = "fancy"
}

run "resource_filter_with_tags_predicate_contains" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_contains : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = var.predicate_key_data,
        value    = var.predicate_value
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
      contains(var.predicate_types_contains, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_contains, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_contains,
      tolist([opslevel_filter.this.predicate[0].type, opslevel_filter.this.predicate[1].type])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key_data == var.predicate_key_data,
      opslevel_filter.this.predicate[1].key_data == var.predicate_key_data
    ])
    error_message = format(
      "expected predicate key_data to be '%v' got '%v'",
      var.predicate_key_data,
      tolist([opslevel_filter.this.predicate[0].key_data, opslevel_filter.this.predicate[1].key_data])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == var.predicate_value,
      opslevel_filter.this.predicate[1].value == var.predicate_value
    ])
    error_message = format(
      "expected predicate value to be '%v' got '%v'",
      var.predicate_value,
      tolist([opslevel_filter.this.predicate[0].value, opslevel_filter.this.predicate[1].value])
    )
  }

}

run "resource_filter_with_tags_predicate_equals" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_equals : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = var.predicate_key_data,
        value    = var.predicate_value
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
      opslevel_filter.this.predicate[0].key_data == var.predicate_key_data,
      opslevel_filter.this.predicate[1].key_data == var.predicate_key_data
    ])
    error_message = format(
      "expected predicate key_data to be '%v' got '%v'",
      var.predicate_key_data,
      tolist([opslevel_filter.this.predicate[0].key_data, opslevel_filter.this.predicate[1].key_data])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == var.predicate_value,
      opslevel_filter.this.predicate[1].value == var.predicate_value
    ])
    error_message = format(
      "expected predicate value to be '%v' got '%v'",
      var.predicate_value,
      tolist([opslevel_filter.this.predicate[0].value, opslevel_filter.this.predicate[1].value])
    )
  }

}

run "resource_filter_with_tags_predicate_exists" {

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
    error_message = format(
      "expected predicate key_data to be '%v' got '%v'",
      var.predicate_key_data,
      tolist([opslevel_filter.this.predicate[0].key_data, opslevel_filter.this.predicate[1].key_data])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == null,
      opslevel_filter.this.predicate[1].value == null
    ])
    error_message = var.error_expected_null_field
  }

}

run "resource_filter_with_tags_predicate_matches_regex" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_matches_regex : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = var.predicate_key_data,
        value    = var.predicate_value
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
      contains(var.predicate_types_matches_regex, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_matches_regex, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_matches_regex,
      tolist([opslevel_filter.this.predicate[0].type, opslevel_filter.this.predicate[1].type])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key_data == var.predicate_key_data,
      opslevel_filter.this.predicate[1].key_data == var.predicate_key_data
    ])
    error_message = format(
      "expected predicate key_data to be '%v' got '%v'",
      var.predicate_key_data,
      tolist([opslevel_filter.this.predicate[0].key_data, opslevel_filter.this.predicate[1].key_data])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == var.predicate_value,
      opslevel_filter.this.predicate[1].value == var.predicate_value
    ])
    error_message = format(
      "expected predicate value to be '%v' got '%v'",
      var.predicate_value,
      tolist([opslevel_filter.this.predicate[0].value, opslevel_filter.this.predicate[1].value])
    )
  }

}

run "resource_filter_with_tags_predicate_starts_or_ends_with" {

  variables {
    predicates = [
      for predicate_type in var.predicate_types_ends_or_starts_with : {
        key      = var.predicate_key
        type     = predicate_type,
        key_data = var.predicate_key_data,
        value    = var.predicate_value
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
      contains(var.predicate_types_ends_or_starts_with, opslevel_filter.this.predicate[0].type),
      contains(var.predicate_types_ends_or_starts_with, opslevel_filter.this.predicate[1].type)
    ])
    error_message = format(
      "expected predicate types to be one of '%v' got '%v'",
      var.predicate_types_ends_or_starts_with,
      tolist([opslevel_filter.this.predicate[0].type, opslevel_filter.this.predicate[1].type])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].key_data == var.predicate_key_data,
      opslevel_filter.this.predicate[1].key_data == var.predicate_key_data
    ])
    error_message = format(
      "expected predicate key_data to be '%v' got '%v'",
      var.predicate_key_data,
      tolist([opslevel_filter.this.predicate[0].key_data, opslevel_filter.this.predicate[1].key_data])
    )
  }

  assert {
    condition = alltrue([
      opslevel_filter.this.predicate[0].value == var.predicate_value,
      opslevel_filter.this.predicate[1].value == var.predicate_value
    ])
    error_message = format(
      "expected predicate value to be '%v' got '%v'",
      var.predicate_value,
      tolist([opslevel_filter.this.predicate[0].value, opslevel_filter.this.predicate[1].value])
    )
  }

}

run "resource_filter_with_tags_predicate_satisfies_version_constraint" {

  variables {
    predicates = [
      {
        key      = var.predicate_key
        type     = var.predicate_types_satisfies_version_constraint,
        key_data = var.predicate_key_data,
        value    = var.predicate_value
      }
    ]
  }

  module {
    source = "./opslevel_modules/modules/filter"
  }

  assert {
    condition = opslevel_filter.this.predicate[0].key == var.predicate_key
    error_message = format(
      "expected predicate key to be '%v' got key '%v'",
      var.predicate_key,
      opslevel_filter.this.predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.this.predicate[0].type == var.predicate_types_satisfies_version_constraint
    error_message = format(
      "expected predicate type to be '%v' got '%v'",
      var.predicate_types_satisfies_version_constraint,
      opslevel_filter.this.predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.this.predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data to be '%v' got '%v'",
      var.predicate_key_data,
      opslevel_filter.this.predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.this.predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value to be '%v' got '%v'",
      var.predicate_value,
      opslevel_filter.this.predicate[0].value
    )
  }

}
