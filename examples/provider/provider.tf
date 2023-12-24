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

data "msgraph_user" "user" {
	id = "0acc8010-50ea-4a54-bd71-ec485d425a74"
	//user_principal_name = "AdeleV@msgraphtf.onmicrosoft.com"
}

output "my_data" {
	value = data.msgraph_user.user
}

