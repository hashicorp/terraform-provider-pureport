data "pureport_accounts" "terraform_acceptance_tests" {
  name_regex = "Terraform Acceptance Tests"
}

resource "pureport_network" "connections" {
  name         = "Connections"
  description  = "Data Source Testing"
  account_href = data.pureport_accounts.terraform_acceptance_tests.accounts[0].href
}

resource "pureport_network" "siouxsie" {
  name         = "Siouxsie & The Banshees"
  description  = "Test Network #2"
  account_href = data.pureport_accounts.terraform_acceptance_tests.accounts[0].href
}

resource "pureport_network" "the_clash" {
  name         = "The Clash"
  description  = "Test Network #1"
  account_href = data.pureport_accounts.terraform_acceptance_tests.accounts[0].href
}

resource "pureport_network" "flock_of_seagulls" {
  name         = "A Flock of Seagulls"
  description  = "Test Network DataSource"
  account_href = data.pureport_accounts.terraform_acceptance_tests.accounts[0].href
}

module "connections" {
  source                   = "./connections"
  connections_network_href = pureport_network.connections.href
  datasource_network_href  = pureport_network.flock_of_seagulls.href
  datasource_express_route = var.datasource_express_route
  google_compute_network   = var.google_compute_network
}
