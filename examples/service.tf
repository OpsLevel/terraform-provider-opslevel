resource "opslevel_service" "test" {
  name    = "spaghetti and meatballs"

  description = "test service"
  framework   = "rails"
  language    = "ruby"

  aliases = [
    "meatballs"
  ]

  tags = {
    foo = "bar"
  }
}