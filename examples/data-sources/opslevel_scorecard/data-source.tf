data "opslevel_scorecard" "foo" {
  identifier = "foo"
}

data "opslevel_scorecard" "bar" {
  identifier = "Z2lkOi8vb3BzbGV2ZWwvU2VydmljZS84Njcw"
}

output "foo_category_ids" {
  value = flatten(data.opslevel_scorecard.foo.categories[*].id)
}

output "foo_categories" {
  value = data.opslevel_scorecard.foo.categories
}
