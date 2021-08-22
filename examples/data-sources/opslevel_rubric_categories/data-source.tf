data "opslevel_rubric_categories" "all" {
}

output "found" {
  value = data.opslevel_rubric_categories.all.ids[0]
}