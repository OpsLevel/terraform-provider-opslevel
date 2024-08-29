resource "opslevel_tag" "test" {
  key                 = var.key
  resource_identifier = var.resource_identifier
  resource_type       = var.resource_type
  value               = var.value
}
