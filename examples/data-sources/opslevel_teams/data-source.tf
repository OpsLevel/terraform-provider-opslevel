data "opslevel_teams" "all" {
}

data "opslevel_teams" "leet" {
  filter {
    field = "manager-email"
    value = "0p5l3v3l@example.com"
  }
}

output "found" {
  value = data.opslevel_teams.all.ids[3]
}