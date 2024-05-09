data "opslevel_scorecards" "all" {}

output "all" {
  value = data.opslevel_scorecards.all.scorecards
}

output "scorecard_names" {
  value = sort(data.opslevel_scorecards.all.scorecards[*].name)
}
