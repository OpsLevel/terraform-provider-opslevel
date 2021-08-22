data "opslevel_filters" "all" {
}

output "found" {
  value = data.opslevel_filters.all.id[0]
}