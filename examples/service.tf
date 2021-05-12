resource "opslevel_service" "test" {
  name    = "spaghetti"

  description = "test service"
  framework   = "rails"
  language    = "ruby"

  tags = {
    foo = "bar"
  }
}