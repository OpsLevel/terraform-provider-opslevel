mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_scorecard_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_scorecard.big.affects_overall_service_levels == false
    error_message = "wrong affects_overall_service_levels in opslevel_scorecard.big"
  }

  assert {
    condition     = startswith(element(opslevel_scorecard.big.categories, 0), "Z2lkOi8v")
    error_message = "expected category id that starts with 'Z2lkOi8v'"
  }

  assert {
    condition     = opslevel_scorecard.big.description == "This is a big scorecard"
    error_message = "wrong description in opslevel_scorecard.big"
  }

  assert {
    condition     = opslevel_scorecard.big.filter_id == var.test_id
    error_message = "wrong filter_id opslevel_scorecard.big"
  }

  assert {
    condition     = can(opslevel_scorecard.big.id)
    error_message = "id attribute missing from filter in opslevel_scorecard.big"
  }

  assert {
    condition     = opslevel_scorecard.big.name == "Big Scorecard"
    error_message = "wrong name in opslevel_scorecard.big"
  }

  assert {
    condition     = opslevel_scorecard.big.owner_id == var.test_id
    error_message = "wrong owner_id in opslevel_scorecard.big"
  }

}

run "resource_scorecard_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_scorecard.small.affects_overall_service_levels == true
    error_message = "wrong affects_overall_service_levels in opslevel_scorecard.small"
  }

  assert {
    condition     = can(opslevel_scorecard.small.id)
    error_message = "id attribute missing from filter in opslevel_scorecard.small"
  }

  assert {
    condition     = opslevel_scorecard.small.owner_id == var.test_id
    error_message = "wrong owner_id in opslevel_scorecard.small"
  }

  assert {
    condition     = opslevel_scorecard.small.name == "Small Scorecard"
    error_message = "wrong name in opslevel_scorecard.small"
  }

}
