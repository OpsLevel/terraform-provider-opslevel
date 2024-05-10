data "opslevel_domains" "all" {}

data "opslevel_domain" "first_domain_by_alias" {
  identifier = data.opslevel_domains.all.domains[0].aliases[0]
}

data "opslevel_domain" "first_domain_by_id" {
  identifier = data.opslevel_domains.all.domains[0].id
}
