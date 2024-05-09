data "opslevel_lifecycles" "all" {}

output "all" {
  value = data.opslevel_lifecycles.all.lifecycles
}

output "lifecycle_names" {
  value = sort(data.opslevel_lifecycles.all.lifecycles[*].name)
}
