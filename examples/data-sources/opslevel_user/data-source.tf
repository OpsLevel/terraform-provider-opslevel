data "opslevel_user" "foo" {
  identifier = "foo@example.com"
}

output "foo" {
  value = data.opslevel_user.foo
}
