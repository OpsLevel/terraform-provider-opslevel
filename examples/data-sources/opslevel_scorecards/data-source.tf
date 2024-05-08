data "opslevel_scorecards" "all" {}

output "found" {
  value = data.opslevel_scorecards.all.ids[0]
}
