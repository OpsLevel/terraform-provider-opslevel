resource "opslevel_property_assignment" "example" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODA"
  owner = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value = jsonencode("green")
}

resource "opslevel_property_assignment" "example_2" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi85OA"
  owner = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value = "true"
}
