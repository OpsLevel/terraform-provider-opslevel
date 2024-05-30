variables {
  domain_one  = "opslevel_domain"
  domains_all = "opslevel_domains"

  # required fields
  name = "TF Test Domain"

  # optional fields
  description = "Domain description"
  note        = "Domain note"
  owner_id    = null
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

run "resource_domain_create_with_all_fields" {

  variables {
    description = var.description
    owner_id    = run.from_team_get_owner_id.first_team.id
    name        = var.name
    note        = var.note
  }

  module {
    source = "./domain"
  }

  assert {
    condition = alltrue([
      can(opslevel_domain.test.aliases),
      can(opslevel_domain.test.description),
      can(opslevel_domain.test.id),
      can(opslevel_domain.test.name),
      can(opslevel_domain.test.note),
      can(opslevel_domain.test.owner),
    ])
    error_message = replace(var.error_unexpected_resource_fields, "TYPE", var.domain_one)
  }

  assert {
    condition     = opslevel_domain.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.domain_one)
  }

  assert {
    condition     = startswith(opslevel_domain.test.id, var.id_prefix)
    error_message = replace(var.error_wrong_id, "TYPE", var.domain_one)
  }

  assert {
    condition     = opslevel_domain.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.domain_one)
  }

  assert {
    condition     = opslevel_domain.test.note == var.note
    error_message = "wrong note for opslevel_domain resource"
  }

  assert {
    condition     = opslevel_domain.test.owner == var.owner_id
    error_message = "wrong owner of opslevel_domain resource"
  }

}

run "resource_domain_update_unset_optional_fields" {

  variables {
    description = null
    owner_id    = null
    note        = null
  }

  module {
    source = "./domain"
  }

  assert {
    condition     = opslevel_domain.test.description == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_domain.test.owner == null
    error_message = var.error_expected_null_field
  }

  assert {
    condition     = opslevel_domain.test.note == null
    error_message = var.error_expected_null_field
  }

}

run "resource_domain_update_all_fields" {

  variables {
    description = "${var.description} updated"
    owner_id    = run.from_team_get_owner_id.first_team.id
    name        = "${var.name} updated"
    note        = "${var.note} updated"
  }

  module {
    source = "./domain"
  }

  assert {
    condition     = opslevel_domain.test.description == var.description
    error_message = replace(var.error_wrong_description, "TYPE", var.domain_one)
  }

  assert {
    condition     = opslevel_domain.test.name == var.name
    error_message = replace(var.error_wrong_name, "TYPE", var.domain_one)
  }

  assert {
    condition     = opslevel_domain.test.note == var.note
    error_message = "wrong note for ${var.domain_one}"
  }

  assert {
    condition     = opslevel_domain.test.owner == var.owner_id
    error_message = "wrong owner of opslevel_domain resource"
  }

}

run "datasource_domains_list_all" {

  variables {
    name = run.resource_domain_update_all_fields.test.name
  }

  module {
    source = "./domain"
  }

  assert {
    condition     = can(data.opslevel_domains.all.domains)
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.domains_all)
  }

  assert {
    condition     = length(data.opslevel_domains.all.domains) > 0
    error_message = replace(var.error_empty_datasource, "TYPE", var.domains_all)
  }

}

run "datasource_domain_get_first" {

  variables {
    name = run.resource_domain_update_all_fields.test.name
  }

  module {
    source = "./domain"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_domain.first_domain_by_id.aliases),
      can(data.opslevel_domain.first_domain_by_id.description),
      can(data.opslevel_domain.first_domain_by_id.id),
      can(data.opslevel_domain.first_domain_by_id.identifier),
      can(data.opslevel_domain.first_domain_by_id.name),
      can(data.opslevel_domain.first_domain_by_id.owner),
    ])
    error_message = replace(var.error_unexpected_datasource_fields, "TYPE", var.domain_one)
  }

  assert {
    condition     = data.opslevel_domain.first_domain_by_alias.identifier == data.opslevel_domain.first_domain_by_alias.aliases[0]
    error_message = replace(var.error_wrong_alias, "TYPE", var.domain_one)
  }

  assert {
    condition     = data.opslevel_domain.first_domain_by_id.identifier == data.opslevel_domain.first_domain_by_id.id
    error_message = replace(var.error_wrong_id, "TYPE", var.domain_one)
  }

}
