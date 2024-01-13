data "opslevel_service" "foo" {
  id = "xzc098zxc098zxc098"
}

data "opslevel_service" "bar" {
  alias = "bar"
}

resource "opslevel_service_dependency" "example" {
  service    = data.opslevel_service.foo.alias
  depends_upon = data.opslevel_service.bar.alias
  note       = <<-EOT
    This is an example of notes on a service dependency
  EOT
}
