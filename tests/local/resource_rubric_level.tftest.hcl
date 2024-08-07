mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_rubric_level_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = can(opslevel_rubric_level.small.id)
    error_message = "id attribute missing from filter in opslevel_rubric_level.small"
  }

  assert {
    condition     = opslevel_rubric_level.small.index == 2
    error_message = "wrong index for opslevel_rubric_level.small"
  }

  assert {
    condition     = opslevel_rubric_level.small.name == "small rubric level"
    error_message = "wrong name for opslevel_rubric_level.small"
  }

}

run "resource_rubric_level_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_rubric_level.big.description == "big rubric description"
    error_message = "wrong description for opslevel_rubric_level.big"
  }

  assert {
    condition     = opslevel_rubric_level.big.index == 5
    error_message = "wrong index for opslevel_rubric_level.big"
  }

  assert {
    condition     = opslevel_rubric_level.big.name == "big rubric level"
    error_message = "wrong name for opslevel_rubric_level.big"
  }

}
