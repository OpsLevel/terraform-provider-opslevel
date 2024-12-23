data "opslevel_service" "foo" {
  id = "xzc098zxc098zxc098"
}

data "opslevel_system" "bar" {
  identifier = "bar"
}

resource "opslevel_service_relationship" "example" {
  service = data.opslevel_service.foo.alias
  system  = data.opslevel_system.bar.alias
}
