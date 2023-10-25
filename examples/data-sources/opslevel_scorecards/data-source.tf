data "opslevel_scorecards" "all" {
}

output "found" {
  value = data.opslevel_scorecards.bar.all.ids[0]
}
