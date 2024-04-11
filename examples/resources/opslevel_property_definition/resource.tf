resource "opslevel_property_definition" "color_picker" {
  name   = "Color Picker"
  schema = jsonencode({
    "type": "string",
    "enum": [
      "red",
      "green",
      "blue",
    ]
  })
  allowed_in_config_files = false
  property_display_status = "visible"
}
