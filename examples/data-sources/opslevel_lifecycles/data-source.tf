data "opslevel_lifecycles" "all" {
}

output "found" {
  value = data.opslevel_lifecycles.all.aliases[0]
}