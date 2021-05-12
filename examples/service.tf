resource "opslevel_service" "test" {
  name = "spaghetti and meatballs"

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

data "opslevel_service" "test" {
  filter {
    field = "alias"
    value = opslevel_service.test.aliases.0
  }
}

data "opslevel_service" "django" {
  filter {
    field = "framework"
    value = "django"
  }
}

data "opslevel_service" "zkms" {
  filter {
    field = "tag"
    value = "zkms:true"
  }
}

output "test_service_id" {
  value = opslevel_service.test.id
}

output "found_services" {
  value = data.opslevel_service.test
}

output "django_services" {
  value = data.opslevel_service.django
}

output "zkms_services" {
  value = data.opslevel_service.zkms
}