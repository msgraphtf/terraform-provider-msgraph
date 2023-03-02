terraform {
  required_providers {
    msgraph = {
      source = "99-lives/msgraph"
    }
  }
}

provider "msgraph" {
  # example configuration here
}

data "msgraph_users" "users" {}

output "my_users" {
	value = data.msgraph_users.users
}

