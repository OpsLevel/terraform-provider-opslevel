variables {
  name            = "TF Test Filter with language predicate"
  predicate_value = "fancy"
  language_predicates = setproduct(
    ["language"],
    concat([
      "contains",
      "does_not_contain",
      "does_not_match_regex",
      "ends_with",
      "matches_regex",
      "starts_with",
      ],
      var.predicate_types_equals,
      var.predicate_types_exists
    ),
  )
}

run "resource_filter_with_language_predicate_contains" {

  variables {
    predicates = tomap({
      for pair in var.language_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = var.predicate_value
      }
      if contains(["does_not_contain", "contains"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].type == "does_not_contain"
    error_message = format(
      "expected predicate type 'does_not_contain' got '%s'",
      opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_does_not_contain"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_contains"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_contains"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_contains"].predicate[0].type == "contains"
    error_message = format(
      "expected predicate type 'does_not_contain' got '%s'",
      opslevel_filter.all_predicates["language_contains"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_contains"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_contains"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_contains"].predicate[0].value
    )
  }

}

run "resource_filter_with_language_predicate_equals" {

  variables {
    predicates = tomap({
      for pair in var.language_predicates : "${pair[0]}_${pair[1]}" => {
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
    condition = opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].type == "does_not_equal"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_does_not_equal"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_equals"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_equals"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_equals"].predicate[0].type == "equals"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["language_equals"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_equals"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_equals"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_equals"].predicate[0].value
    )
  }

}

run "resource_filter_with_language_predicate_exists" {

  variables {
    predicates = tomap({
      for pair in var.language_predicates : "${pair[0]}_${pair[1]}" => {
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
    condition = opslevel_filter.all_predicates["language_does_not_exist"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["language_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_does_not_exist"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_exists"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["language_exists"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_exists"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}

run "resource_filter_with_language_predicate_matches_regex" {

  variables {
    predicates = tomap({
      for pair in var.language_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = var.predicate_value
      }
      if contains(["does_not_match_regex", "matches_regex"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].type == "does_not_match_regex"
    error_message = format(
      "expected predicate type 'does_not_match_regex' got '%s'",
      opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_does_not_match_regex"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_matches_regex"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_matches_regex"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_matches_regex"].predicate[0].type == "matches_regex"
    error_message = format(
      "expected predicate type 'does_not_match_regex' got '%s'",
      opslevel_filter.all_predicates["language_matches_regex"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_matches_regex"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_matches_regex"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_matches_regex"].predicate[0].value
    )
  }

}

run "resource_filter_with_language_predicate_starts_or_ends_with" {

  variables {
    predicates = tomap({
      for pair in var.language_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = null,
        value    = var.predicate_value
      }
      if contains(["ends_with", "starts_with"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["language_ends_with"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_ends_with"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_ends_with"].predicate[0].type == "ends_with"
    error_message = format(
      "expected predicate type 'ends_with' got '%s'",
      opslevel_filter.all_predicates["language_ends_with"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_ends_with"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_ends_with"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_ends_with"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_starts_with"].predicate[0].key == "language"
    error_message = format(
      "expected predicate key 'language' got '%s'",
      opslevel_filter.all_predicates["language_starts_with"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["language_starts_with"].predicate[0].type == "starts_with"
    error_message = format(
      "expected predicate type 'ends_with' got '%s'",
      opslevel_filter.all_predicates["language_starts_with"].predicate[0].type
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["language_starts_with"].predicate[0].key_data == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["language_starts_with"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["language_starts_with"].predicate[0].value
    )
  }

}
