run "datasource_domains_all" {

  variables {
    datasource_type = "opslevel_domains"
  }

  assert {
    condition     = can(data.opslevel_domains.all.domains)
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = length(data.opslevel_domains.all.domains) > 0
    error_message = replace(var.empty_datasource_error, "TYPE", var.datasource_type)
  }

}

run "datasource_domain_first" {

  variables {
    datasource_type = "opslevel_domain"
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
    error_message = replace(var.unexpected_datasource_fields_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_domain.first_domain_by_alias.identifier == data.opslevel_domain.first_domain_by_alias.aliases[0]
    error_message = replace(var.wrong_alias_error, "TYPE", var.datasource_type)
  }

  assert {
    condition     = data.opslevel_domain.first_domain_by_id.identifier == data.opslevel_domain.first_domain_by_id.id
    error_message = replace(var.wrong_id_error, "TYPE", var.datasource_type)
  }

}
