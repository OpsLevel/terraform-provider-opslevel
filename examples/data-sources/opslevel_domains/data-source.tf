data "opslevel_domains" "all" {}

output "found" {
  value = data.opslevel_domains.all
}
