variables {
  name               = "TF Test Filter with tags predicate"
  predicate_key_data = "test_tag"
  predicate_value    = "fancy"
  tags_predicates = setproduct(
    ["tags"],
    concat([
      "contains",
      "does_not_contain",
      "does_not_match_regex",
      "ends_with",
      "matches_regex",
      "satisfies_version_constraint",
      "starts_with",
      ],
      var.predicate_types_equals,
      var.predicate_types_exists
    ),
  )
}

run "resource_filter_with_tags_predicate_contains" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.tags_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = var.predicate_key_data,
        value    = var.predicate_value
      }
      if contains(["does_not_contain", "contains"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].type == "does_not_contain"
    error_message = format(
      "expected predicate type 'does_not_contain' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_does_not_contain"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_contains"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_contains"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_contains"].predicate[0].type == "contains"
    error_message = format(
      "expected predicate type 'does_not_contain' got '%s'",
      opslevel_filter.all_predicates["tags_contains"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_contains"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_contains"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_contains"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_contains"].predicate[0].value
    )
  }

}

run "resource_filter_with_tags_predicate_equals" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.tags_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = var.predicate_key_data,
        value    = var.predicate_value
      }
      if contains(["does_not_equal", "equals"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].type == "does_not_equal"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_does_not_equal"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_equals"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_equals"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_equals"].predicate[0].type == "equals"
    error_message = format(
      "expected predicate type 'does_not_equal' got '%s'",
      opslevel_filter.all_predicates["tags_equals"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_equals"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_equals"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_equals"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_equals"].predicate[0].value
    )
  }

}

run "resource_filter_with_tags_predicate_exists" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.tags_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = var.predicate_key_data,
        value    = null
      }
      if contains(["does_not_exist", "exists"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].type == "does_not_exist"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].key_data
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["tags_does_not_exist"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_exists"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_exists"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_exists"].predicate[0].type == "exists"
    error_message = format(
      "expected predicate type 'does_not_exist' got '%s'",
      opslevel_filter.all_predicates["tags_exists"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_exists"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_exists"].predicate[0].key_data
    )
  }

  assert {
    condition     = opslevel_filter.all_predicates["tags_exists"].predicate[0].value == null
    error_message = var.error_expected_null_field
  }

}

run "resource_filter_with_tags_predicate_matches_regex" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.tags_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = var.predicate_key_data,
        value    = var.predicate_value
      }
      if contains(["does_not_match_regex", "matches_regex"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].type == "does_not_match_regex"
    error_message = format(
      "expected predicate type 'does_not_match_regex' got '%s'",
      opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_does_not_match_regex"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].type == "matches_regex"
    error_message = format(
      "expected predicate type 'does_not_match_regex' got '%s'",
      opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_matches_regex"].predicate[0].value
    )
  }

}

run "resource_filter_with_tags_predicate_starts_or_ends_with" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.tags_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = var.predicate_key_data,
        value    = var.predicate_value
      }
      if contains(["ends_with", "starts_with"], pair[1])
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_ends_with"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_ends_with"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_ends_with"].predicate[0].type == "ends_with"
    error_message = format(
      "expected predicate type 'ends_with' got '%s'",
      opslevel_filter.all_predicates["tags_ends_with"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_ends_with"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_ends_with"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_ends_with"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_ends_with"].predicate[0].value
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_starts_with"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_starts_with"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_starts_with"].predicate[0].type == "starts_with"
    error_message = format(
      "expected predicate type 'ends_with' got '%s'",
      opslevel_filter.all_predicates["tags_starts_with"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_starts_with"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_starts_with"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_starts_with"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_starts_with"].predicate[0].value
    )
  }

}

run "resource_filter_with_tags_predicate_satisfies_version_constraint" {

  variables {
    connective = "and"
    predicates = tomap({
      for pair in var.tags_predicates : "${pair[0]}_${pair[1]}" => {
        key      = pair[0],
        type     = pair[1],
        key_data = var.predicate_key_data,
        value    = var.predicate_value
      }
      if "satisfies_version_constraint" == pair[1]
    })
  }

  module {
    source = "./filter"
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].key == "tags"
    error_message = format(
      "expected predicate key 'tags' got '%s'",
      opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].key
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].type == "satisfies_version_constraint"
    error_message = format(
      "expected predicate type 'ends_with' got '%s'",
      opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].type
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].key_data == var.predicate_key_data
    error_message = format(
      "expected predicate key_data '%s' got '%s'",
      var.predicate_key_data,
      opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].key_data
    )
  }

  assert {
    condition = opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].value == var.predicate_value
    error_message = format(
      "expected predicate value '%s' got '%s'",
      var.predicate_value,
      opslevel_filter.all_predicates["tags_satisfies_version_constraint"].predicate[0].value
    )
  }

}
