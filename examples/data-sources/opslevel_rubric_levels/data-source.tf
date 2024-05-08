data "opslevel_rubric_levels" "all" {}

output "found" {
  value = data.opslevel_rubric_levels.all.ids[0]
}
