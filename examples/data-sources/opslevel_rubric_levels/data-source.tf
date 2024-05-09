data "opslevel_rubric_levels" "all" {}

output "all" {
  value = data.opslevel_rubric_levels.all.rubric_levels
}

output "level_names" {
  value = sort(data.opslevel_rubric_levels.all.rubric_levels[*].name)
}
