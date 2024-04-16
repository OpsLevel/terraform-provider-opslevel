mock_data "opslevel_team" {
  defaults = {
    alias        = "platform"
    parent_id    = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xMDI0Mg"
    parent_alias = "engineering"
    members      = [{ "email" : "person1@opslevel.com", "role" : "manager" }, { "email" : "person2@opslevel.com", "role" : "contributor" }]
    name         = "Platform"
  }
}

