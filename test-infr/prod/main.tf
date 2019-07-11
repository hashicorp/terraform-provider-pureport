provider "azurerm" {
  version = "~> 1.30"
}

module "azure-infra" {
  source              = "../global/azure-express-route"
  resource_group_name = "terraform-acceptance-tests"
  env                 = "prod"
}

module "google-infra" {
  source = "../global/google-cloud-interconnect"
  env    = "prod"
}

module "pureport-infra" {
  source                   = "./pureport-infra"
  datasource_express_route = module.azure-infra.datasource_express_route
  google_compute_network   = module.google-infra.compute_network
}

