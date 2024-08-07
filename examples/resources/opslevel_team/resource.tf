data "opslevel_team" "parent" {
  alias = "platform"
}

resource "opslevel_team" "example" {
  name             = "foo"
  responsibilities = "Responsible for foo frontend and backend"
  aliases          = ["foo", "bar", "baz"] # NOTE: if set, slugified value of "name" must be included
  parent           = data.opslevel_team.parent.id

  member {
    email = "john.doe@example.com"
    role  = "manager"
  }
  member {
    email = "jane.doe@example.com"
    role  = "contributor"
  }
}

output "team" {
  value = opslevel_team.example.id
}
