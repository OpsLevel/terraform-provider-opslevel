variable "test_id" {
  type        = string
  description = "id for testing"
  default     = "Z2lkOi8vb3BzbGV2ZWwvVGllci8yMTAw"

  validation {
    condition     = startswith(var.test_id, "Z2lkOi8v")
    error_message = "expected test_id to start with Z2lkOi8v"
  }
}

