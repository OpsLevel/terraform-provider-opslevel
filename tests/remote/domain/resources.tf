resource "opslevel_domain" "required_fields" {
  name = "Test - name only"
}

resource "opslevel_domain" "all_fields" {
  name        = "Test - all fields"
  description = "The whole app in one monolith"
  owner       = "Z2lkOi8vb3BzbGV2ZWwvVGVhbS8xNzQzNA"
  note        = "This is an example"
}
