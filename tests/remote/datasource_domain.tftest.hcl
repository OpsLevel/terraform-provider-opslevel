run "datasource_domains_all" {

  assert {
    condition     = length(data.opslevel_domains.all.domains) > 0
    error_message = "zero Domains found in data.opslevel_domains"
  }

  assert {
    condition = alltrue([
      can(data.opslevel_domains.all.domains[0].aliases),
      can(data.opslevel_domains.all.domains[0].description),
      can(data.opslevel_domains.all.domains[0].id),
      can(data.opslevel_domains.all.domains[0].name),
      can(data.opslevel_domains.all.domains[0].owner),
    ])
    error_message = "cannot set all expected Domain datasource fields"
  }

}

run "datasource_domain_first" {

  assert {
    condition     = data.opslevel_domain.first_domain_by_alias.identifier == data.opslevel_domain.first_domain_by_alias.aliases[0]
    error_message = "wrong alias on opslevel_domain"
  }

  assert {
    condition     = data.opslevel_domain.first_domain_by_id.identifier == data.opslevel_domain.first_domain_by_id.id
    error_message = "wrong ID on opslevel_domain"
  }

}
