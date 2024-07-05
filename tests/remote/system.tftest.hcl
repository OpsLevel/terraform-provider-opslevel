variables {
  system_one  = "opslevel_system"
  systems_all = "opslevel_systems"

  # required fields
  name = "TF Test System"

  # optional fields
  description = "System description"
  domain_id   = null
  note        = "System note"
  owner_id    = null
}

run "from_domain_get_domain_id" {
  command = plan

  variables {
    description = null
    name        = ""
    note        = null
    owner_id    = null
  }

  module {
    source = "./domain"
  }
}

run "from_team_get_owner_id" {
  command = plan

  variables {
    aliases          = null
    name             = ""
    parent           = null
    responsibilities = null
  }

  module {
    source = "./team"
  }
}

run "resource_system_create_with_all_fields" {

  variables {
    description = var.description
    domain_id   = run.from_domain_get_domain_id.first_domain.id
    name        = var.name
    note        = var.note
    owner_id    = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./system"
  }

  assert {
    condition = alltrue([
      can(opslevel_system.test.aliases),
      can(opslevel_system.test.description),
      can(opslevel_system.test.domain),
      can(opslevel_system.test.id),
      can(opslevel_system.test.name),
      can(opslevel_system.test.note),
      can(opslevel_system.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.system_one)
  }

  assert {
    condition     = opslevel_system.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.system_one)
  }

  assert {
    condition     = opslevel_system.test.domain == var.domain_id
    error_message = "wrong domain ID for opslevel_system resource"
  }

  assert {
    condition     = startswith(opslevel_system.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.system_one)
  }

  assert {
    condition     = opslevel_system.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.system_one)
  }

  assert {
    condition     = opslevel_system.test.note == var.note
    error_message = "wrong note for opslevel_system resource"
  }

  assert {
    condition     = opslevel_system.test.owner == var.owner_id
    error_message = "wrong owner ID for opslevel_system resource"
  }

}

run "resource_system_create_with_empty_optional_fields" {

  variables {
    description = ""
    name        = "New ${var.name} with empty fields"
    note        = ""
  }

  module {
    source = "./system"
  }

  assert {
    condition     = opslevel_system.test.description == ""
    error_message = var.error_expected_empty_string
  }

  assert {
    condition     = opslevel_system.test.note == ""
    error_message = var.error_expected_empty_string
  }

}

run "resource_system_update_unset_optional_fields" {

  variables {
    description = null
    domain_id   = null
    note        = null
    owner_id    = null
  }

  module {
    source = "./system"
  }

  assert {
    condition     = opslevel_system.test.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.test.domain == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.test.note == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_system.test.owner == null
    error_message = var.error_expected_null_field
  }

}

run "resource_system_update_set_all_fields" {

  variables {
    description = "${var.description} updated"
    domain_id   = run.from_domain_get_domain_id.first_domain.id
    name        = "${var.name} updated"
    note        = "${var.note} updated"
    owner_id    = run.from_team_get_owner_id.first_team.id
  }

  module {
    source = "./system"
  }

  assert {
    condition     = opslevel_system.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.system_one)
  }

  assert {
    condition     = opslevel_system.test.domain == var.domain_id
    error_message = "wrong domain ID for opslevel_system resource"
  }

  assert {
    condition     = opslevel_system.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.system_one)
  }

  assert {
    condition     = opslevel_system.test.note == var.note
    error_message = "wrong note for opslevel_system resource"
  }

  assert {
    condition     = opslevel_system.test.owner == var.owner_id
    error_message = "wrong owner ID for opslevel_system resource"
  }

}

run "datasource_systems_all" {

  module {
    source = "./system"
  }

  assert {
    condition     = can(data.opslevel_systems.all.systems)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.systems_all)
  }

  assert {
    condition     = length(data.opslevel_systems.all.systems) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.systems_all)
  }

}

run "datasource_system_first" {

  module {
    source = "./system"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_system.first_system_by_id.aliases),
      can(data.opslevel_system.first_system_by_id.description),
      can(data.opslevel_system.first_system_by_id.domain),
      can(data.opslevel_system.first_system_by_id.id),
      can(data.opslevel_system.first_system_by_id.identifier),
      can(data.opslevel_system.first_system_by_id.name),
      can(data.opslevel_system.first_system_by_id.owner),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.system_one)
  }

  assert {
    condition     = data.opslevel_system.first_system_by_id.id == data.opslevel_systems.all.systems[0].id
    error_message = replace(var.error_wrong_id, "TYPE", var.system_one)
  }

}
