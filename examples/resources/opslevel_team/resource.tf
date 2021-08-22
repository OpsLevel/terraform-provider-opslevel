resource "opslevel_team" "example" {
  name = "foo"
  manager_email = "john.doe@example.com"
  responsibilities = "Responsible for foo frontend and backend"
}

output "team" {
  value = opslevel_team.example.id
}
