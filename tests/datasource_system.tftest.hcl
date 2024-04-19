mock_provider "opslevel" {
  alias  = "fake"
  source = "./mock_datasource"
}

run "datasource_system" {
  providers = {
    opslevel = opslevel.fake
  }

  assert {
    condition     = data.opslevel_system.mock_system.aliases == tolist(["my_system", "my_sys"])
    error_message = "wrong aliases in opslevel_system.aliases"
  }

  assert {
    condition     = data.opslevel_system.mock_system.description == "This is my new system that has a domain."
    error_message = "wrong description in opslevel_system.description"
  }

  assert {
    condition     = data.opslevel_system.mock_system.domain == "sys_domain"
    error_message = "wrong domain in opslevel_system.description"
  }

  assert {
    condition     = data.opslevel_system.mock_system.id != null && data.opslevel_system.mock_system.id != ""
    error_message = "opslevel_system id should not be empty"
  }

  assert {
    condition     = data.opslevel_system.mock_system.name == "My New System"
    error_message = "wrong name in opslevel_system.name"
  }

  assert {
    condition     = data.opslevel_system.mock_system.owner == "sys_owner"
    error_message = "opslevel_system owner should be sys_owner"
  }
}

# TODO: re-enable me
#run "datasource_systems_all" {
#  providers = {
#    opslevel = opslevel.fake
#  }
#
#  assert {
#    condition     = length(data.opslevel_systems.all.systems) == 2
#    error_message = "wrong number of objects in opslevel_systems"
#  }
#
#  assert {
#    condition     = data.opslevel_systems.all.systems[0].description == "This is my new system that has a domain." && data.opslevel_systems.all.systems[1].description == "This is my new system that has a domain (2)."
#    error_message = "wrong descriptions in opslevel_systems"
#  }
#}
