resource "opslevel_group" "foo" {
  name        = "foo"
  description = "foo group"
  members     = ["foo@example.com"]
  teams       = ["foo-team"]
}

resource "opslevel_group" "bar" {
  name        = "bar"
  description = "bar group"
  members     = ["bar@example.com"]
  parent      = opslevel_group.foo.alias
  teams       = ["bar-team"]
}
