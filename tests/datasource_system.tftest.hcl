mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_system" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_system.mock_system.aliases[0] == "fancy_system" && data.opslevel_system.mock_system.aliases[1] == "fancy_sys"
    error_message = "wrong aliases in opslevel_system.aliases"
  }

  assert {
    condition     = data.opslevel_system.mock_system.description == "A Fancy API Client"
    error_message = "wrong description in opslevel_system.description"
  }

  assert {
    condition     = data.opslevel_system.mock_system.description == "fancy_domain"
    error_message = "wrong domain in opslevel_system.description"
  }

  assert {
    condition     = data.opslevel_system.mock_system.id != null && data.opslevel_system.mock_system.id != ""
    error_message = "opslevel_system id should not be empty"
  }

  assert {
    condition     = data.opslevel_system.mock_system.name == "Mock System Name"
    error_message = "wrong name in opslevel_system.name"
  }

  assert {
    condition     = data.opslevel_system.mock_system.owner == "system_owner"
    error_message = "opslevel_system owner should be system_owner"
  }
}

run "datasource_systems_all" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = length(data.opslevel_systems.all.systems) == 2
    error_message = "wrong length of opslevel_systems"
  }

  assert {
    condition     = data.opslevel_systems.all.systems[0].description == "Mock system description" && data.opslevel_systems.all.systems[1].description == "Mock system description the second"
    error_message = "wrong descriptions in second opslevel_systems"
  }
}
