data "opslevel_rubric_categories" "all" {}

output "all" {
  value = data.opslevel_rubric_categories.all.rubric_categories
}

output "category_names" {
  value = sort(data.opslevel_rubric_categories.all.rubric_categories[*].name)
}
