output "all" {
  value = data.opslevel_filters.all
}

output "first" {
  value = data.opslevel_filters.all.filters[0]
}
