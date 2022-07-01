data "opslevel_groups" "all" {
}

output "found" {
  value = data.opslevel_groups.all
}