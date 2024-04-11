mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_resource"
}

run "resource_system_big" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_system.big.description == "It's a big system"
    error_message = "wrong description in opslevel_system.big"
  }

  assert {
    condition     = opslevel_system.big.domain == var.test_id
    error_message = "wrong domain id in opslevel_system.big"
  }

  assert {
    condition     = opslevel_system.big.name == "Big System"
    error_message = "wrong name in opslevel_system.big"
  }

  assert {
    condition     = opslevel_system.big.note == "Note on System"
    error_message = "wrong note in opslevel_system.big"
  }

  assert {
    condition     = opslevel_system.big.owner == var.test_id
    error_message = "wrong owner id in opslevel_system.big"
  }

}

run "resource_system_small" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = opslevel_system.small.aliases == tolist(["system-one", "system-two"])
    error_message = "wrong aliases in opslevel_system.small"
  }

  assert {
    condition     = opslevel_system.small.name == "Small System"
    error_message = "wrong name in opslevel_system.small"
  }

}
