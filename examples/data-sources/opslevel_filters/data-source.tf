data "opslevel_filters" "all" {}

output "all" {
  value = data.opslevel_filters.all.filters
}

output "filter_names" {
  value = sort(data.opslevel_filters.all.filters[*].name)
}

