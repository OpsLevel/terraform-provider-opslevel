resource "opslevel_rubric_level" "example" {
  name = "foo"
  description = "foo level"
  index = 2
}

output "level" {
  value = opslevel_rubric_level.example.id
}
