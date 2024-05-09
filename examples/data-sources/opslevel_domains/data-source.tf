data "opslevel_domains" "all" {}

output "all" {
  value = data.opslevel_domains.all.domains
}

output "domain_names" {
  value = sort(data.opslevel_domains.all.domains[*].name)
}
