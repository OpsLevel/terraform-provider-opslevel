mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_system" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_system.mock_system.aliases[0] == "mock-alias-one"
    error_message = "wrong first alias in opslevel_system.aliases"
  }

  assert {
    condition     = data.opslevel_system.mock_system.aliases[1] == "mock-alias-two"
    error_message = "wrong second alias in opslevel_system.aliases"
  }

  assert {
    condition     = length(data.opslevel_system.mock_system.aliases) == 2
    error_message = "wrong number of aliases in opslevel_system.aliases"
  }

  assert {
    condition     = data.opslevel_system.mock_system.description == "mock-system-description"
    error_message = "wrong description in opslevel_system.description"
  }

  assert {
    condition     = data.opslevel_system.mock_system.id != null && data.opslevel_system.mock_system.id != ""
    error_message = "opslevel_system id should not be empty"
  }

  assert {
    condition     = data.opslevel_system.mock_system.name == "mock-system-name"
    error_message = "wrong name in opslevel_system.name"
  }

  assert {
    condition     = data.opslevel_system.mock_system.owner == null && data.opslevel_system.mock_system.owner != ""
    error_message = "opslevel_system owner should be null"
  }
}

run "datasource_systems_all" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = length(data.opslevel_systems.all.systems) == 2
    error_message = "wrong number of owners in opslevel_systems"
  }

  assert {
    condition     = data.opslevel_systems.all.systems[1].description == ""
    error_message = "wrong description in second opslevel_system"
  }
}
