resource "opslevel_filter" "tier1" {
  name = "foo"
  predicate {
    key   = "tier_index"
    type  = "equals"
    value = "1"
  }
  connective = "and"
}

resource "opslevel_filter" "tier2_alpha" {
  name = "foo"
  predicate {
    key   = "tier_index"
    type  = "equals"
    value = "1"
  }
  predicate {
    key   = "lifecycle_index"
    type  = "equals"
    value = "1"
  }
  connective = "and"
}

resource "opslevel_filter" "tier3" {
  name = "foo"
  predicate {
    key      = "tags"
    type     = "does_not_exist"
    key_data = "tier3"
  }
}

resource "opslevel_filter" "case_sensitive" {
  name = "foo"
  predicate {
    key            = "tags"
    type           = "equals"
    key_data       = "my-custom-tag"
    value          = "hello-world"
    case_sensitive = true
  }
}
