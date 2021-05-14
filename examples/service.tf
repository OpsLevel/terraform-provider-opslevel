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

output "test_service_id" {
  value = opslevel_service.test.id
}

data "opslevel_services" "test" {
  filter {
    field = "alias"
    value = opslevel_service.test.aliases.0
  }
}
output "found_services" {
  value = data.opslevel_services.test
}

data "opslevel_services" "django" {
  filter {
    field = "framework"
    value = "django"
  }
}
output "django_services" {
  value = data.opslevel_services.django
}

data "opslevel_services" "feature" {
  filter {
    field = "tag"
    value = "feature:true"
  }
}
output "feature_services" {
  value = data.opslevel_services.feature
}

data "opslevel_services" "tag_present" {
  filter {
    field = "tag"
    value = "important_tag:"
  }
}
output "tagged_services" {
  value = data.opslevel_services.tag_present
}
