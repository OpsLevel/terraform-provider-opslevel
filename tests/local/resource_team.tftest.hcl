mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_team_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition = alltrue([
      contains(opslevel_team.big.aliases, "big_team"),
      contains(opslevel_team.big.aliases, "the_big_team"),
    ])
    error_message = "wrong aliases in opslevel_team.big"
  }

  assert {
    condition     = opslevel_team.big.name == "The Big Team"
    error_message = "wrong display name in opslevel_team.big"
  }

  assert {
    condition     = opslevel_team.big.parent == "small_team"
    error_message = "wrong parent in opslevel_team.big"
  }

  assert {
    condition     = can(opslevel_team.big.id)
    error_message = "id attribute missing from team in opslevel_team.big"
  }

  assert {
    condition     = opslevel_team.big.responsibilities == "This is a big team"
    error_message = "wrong responsibilities in opslevel_team.big"
  }

  assert {
    condition     = opslevel_team.big.member[0].email == "alice@opslevel.com" && opslevel_team.big.member[0].role == "manager" && opslevel_team.big.member[1].email == "bob@opslevel.com" && opslevel_team.big.member[1].role == "contributor"
    error_message = "wrong members in opslevel_team.big"
  }
}

run "resource_team_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_team.small.aliases == null
    error_message = "wrong aliases in opslevel_team.small"
  }

  assert {
    condition     = opslevel_team.small.name == "Small Team"
    error_message = "wrong display name in opslevel_team.small"
  }

  assert {
    condition     = opslevel_team.small.parent == null
    error_message = "wrong parent in opslevel_team.small"
  }

  assert {
    condition     = can(opslevel_team.small.id)
    error_message = "id attribute missing from team in opslevel_team.small"
  }

  assert {
    condition     = opslevel_team.small.responsibilities == null
    error_message = "wrong responsibilities in opslevel_team.small"
  }

  assert {
    condition     = length(opslevel_team.small.member) == 0
    error_message = "wrong members in opslevel_team.small"
  }
}
