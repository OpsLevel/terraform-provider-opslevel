data "opslevel_group" "foo" {
  identifer = "foo"
}

resource "opslevel_team" "example" {
  name = "foo"
  manager_email = "john.doe@example.com"
  responsibilities = "Responsible for foo frontend and backend"
  aliases = ["bar", "baz"]
  group = data.opslevel_group.foo.alias
}

output "team" {
  value = opslevel_team.example.id
}
