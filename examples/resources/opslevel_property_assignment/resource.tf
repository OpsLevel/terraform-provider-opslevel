resource "opslevel_property_assignment" "example" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODA"
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value      = jsonencode("green")
}

resource "opslevel_property_assignment" "example_2" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi85OA"
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value      = jsonencode(true)
}

resource "opslevel_property_assignment" "example_3" {
  definition = "Z2lkOi8vb3BzbGV2ZWwvUHJvcGVydGllczo6RGVmaW5pdGlvbi8xODI"
  owner      = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS85MTQyOQ"
  value      = jsonencode({ "container_id" : "1c6098d6-952a-4062-9293-1dc06e991118", "container_name" : "gcr.io/containername" })
}
