resource "opslevel_team_tag" "team_tag_1" {
  key   = "hello"
  value = "world"
  team  = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQxMg"
}

resource "opslevel_team_tag" "team_tag_2" {
  key        = "hello_with_alias"
  value      = "world_with_alias"
  team_alias = "team_foo"
}
