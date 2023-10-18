data "opslevel_team" "parent" {
  alias = "platform"
}

resource "opslevel_team" "example" {
  name             = "foo"
  manager_email    = "john.doe@example.com"
  members          = ["john.doe@example.com", "jane.doe@example.com"]
  responsibilities = "Responsible for foo frontend and backend"
  aliases          = ["bar", "baz"]
  parent           = data.opslevel_team.parent.id
}

output "team" {
  value = opslevel_team.example.id
}
