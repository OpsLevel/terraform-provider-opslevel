resource "opslevel_property_definition" "newpd" {
  name        = "Friends Property"
  schema     =  jsonencode({"$ref":"#/$defs/MyProp","$defs":{"MyProp":{"properties":{"name":{"type":"string","title":"the new name","description":"The name of a friend","default":"alex","examples":["joe","lucy"]}},"additionalProperties":false,"type":"object","required":["name"]}}})
}
