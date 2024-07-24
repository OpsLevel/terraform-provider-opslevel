mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_filter_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_filter.small.name == "Blank Filter"
    error_message = "wrong name for opslevel_filter.small"
  }

}


run "resource_filter_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_filter.big.connective == var.connective_enum
    error_message = "wrong connective enum for opslevel_filter.big"
  }

  assert {
    condition     = opslevel_filter.big.name == "Big Filter"
    error_message = "wrong name for opslevel_filter.big"
  }

}

run "resource_filter_big_predicate_one" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_filter.big.predicate[0].case_sensitive == null
    error_message = "expected 'case_sensitive' to be null for opslevel_filter.big.predicate[0]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[0].key == var.predicate_key_enum
    error_message = "invalid predicate key enum for opslevel_filter.big.predicate[0]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[0].key_data == null
    error_message = "expected 'key_data' to be null for opslevel_filter.big.predicate[0]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[0].type == var.predicate_type_enum
    error_message = "invalid predicate type enum for opslevel_filter.big.predicate[0]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[0].value == "1"
    error_message = "expected 'value' to be null for opslevel_filter.big.predicate[0]"
  }

}

run "resource_filter_big_predicate_two" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_filter.big.predicate[1].key == "lifecycle_index"
    error_message = "wrong predicate 'key' for opslevel_filter.big.predicate[1]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[1].key_data == "big_predicate"
    error_message = "wrong 'key_data' for opslevel_filter.big.predicate[1]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[1].type == "greater_than_or_equal_to"
    error_message = "wrong predicate 'type' for opslevel_filter.big.predicate[1]"
  }

  assert {
    condition     = opslevel_filter.big.predicate[1].value == "1"
    error_message = "wrong predicate 'value' for opslevel_filter.big.predicate[1]"
  }

}
