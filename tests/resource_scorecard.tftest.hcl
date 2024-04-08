mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_scorecard_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_scorecard.big.id)
    error_message = "id attribute missing from filter in opslevel_scorecard.small"
  }

  assert {
    condition     = opslevel_scorecard.big.name == "Mock Category"
    error_message = "wrong name for opslevel_scorecard.mock_category"
  }

  assert {
    condition     = can(opslevel_scorecard.big.id)
    error_message = "id attribute missing from filter in opslevel_scorecard.small"
  }

}

run "resource_scorecard_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_scorecard.small == "Mock Category"
    error_message = "wrong name for opslevel_scorecard.mock_category"
  }
}
