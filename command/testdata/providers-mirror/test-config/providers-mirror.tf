terraform {
  # Some fake providers for our mirror command to chew on while talking to
  # the fake provider registry (for example.com) created in the test.
  required_providers {
    foo       = { source = "example.com/a/foo" }
    bar       = { source = "example.com/a/bar", version = "< 2" }
    terraform = { source = "terraform.io/builtin/terraform" }
  }
}
