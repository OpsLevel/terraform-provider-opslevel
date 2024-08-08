variables {
  name            = "TF Test Filter with lifecycle_index predicate"
  predicate_value = "1"
  lifecycle_index_predicates = setproduct(
    ["lifecycle_index"],
    concat(
      ["less_than_or_equal_to", "greater_than_or_equal_to"],
      var.predicate_types_equals,
      var.predicate_types_exists
    ),
  )
}

run "resource_filter_with_lifecycle_index_predicate_equals" {

  variables {
    predicates = tomap({
      for pair in var.lifecycle_index_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = var.predicate_value
      }
      if contains(var.predicate_types_equals, pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].key == "lifecycle_index"
    error_message = format(
      "expected predicate key 'lifecycle_index' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].type == "does_not_equal"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["lifecycle_index_does_not_equal"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].key == "lifecycle_index"
    error_message = format(
      "expected predicate key 'lifecycle_index' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].type == "equals"
    error_message = format(
      "expected predicate type 'equals' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["lifecycle_index_equals"].predicate[0].value
    )
  }

}

run "resource_filter_with_lifecycle_index_predicate_exists" {

  variables {
    predicates = tomap({
      for pair in var.lifecycle_index_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = null
      }
      if contains(var.predicate_types_exists, pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_does_not_exist"].predicate[0].key == "lifecycle_index"
    error_message = format(
      "expected predicate key 'lifecycle_index' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_does_not_exist"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_exists"].predicate[0].key == "lifecycle_index"
    error_message = format(
      "expected predicate key 'lifecycle_index' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'exists' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_exists"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_exists"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}

run "resource_filter_with_lifecycle_index_predicate_greater_than_or_equal_to" {

  variables {
    predicates = tomap({
      for pair in var.lifecycle_index_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = var.predicate_value
      }
      if "greater_than_or_equal_to" == pair[1]
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].key == "lifecycle_index"
    error_message = format(
      "expected predicate key 'lifecycle_index' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].type == "greater_than_or_equal_to"
    error_message = format(
      "expected predicate type 'greater_than_or_equal_to' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["lifecycle_index_greater_than_or_equal_to"].predicate[0].value
    )
  }

}

run "resource_filter_with_lifecycle_index_predicate_less_than_or_equal_to" {

  variables {
    predicates = tomap({
      for pair in var.lifecycle_index_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = var.predicate_value
      }
      if "less_than_or_equal_to" == pair[1]
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].key == "lifecycle_index"
    error_message = format(
      "expected predicate key 'lifecycle_index' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].type == "less_than_or_equal_to"
    error_message = format(
      "expected predicate type 'less_than_or_equal_to' got '%s'",
      opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["lifecycle_index_less_than_or_equal_to"].predicate[0].value
    )
  }

}