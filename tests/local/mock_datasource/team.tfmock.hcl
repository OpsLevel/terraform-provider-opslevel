mock_data "opslevel_team" {
  defaults = {
    alias        = "platform"
    parent_id    = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xMDI0Mg"
    parent_alias = "engineering"
    members      = { "email" : "person1@opslevel.com", "role" : "manager" }
    name         = "Platform"
    properties = {
      "definition" = {
        aliases = [
          "mock-one",
          "mock-two",
          "mock-three",
        ]
        id = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8yODk"
      }
      "value" = "mock-property-definition",
    }
    tags = ["key1:value2", "key2:value2"]
  }
}

