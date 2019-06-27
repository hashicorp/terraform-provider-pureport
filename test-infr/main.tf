provider "azurerm" {
  version = "~> 1.30"
}

module "azure-infra" {
  source              = "./modules/azure-express-route"
  resource_group_name = "terraform-acceptance-tests"
}

module "google-infra" {
  source = "./modules/google-cloud-interconnect"
}

module "pureport-infra" {
  source                   = "./modules/pureport-infra"
  datasource_express_route = module.azure-infra.datasource_express_route
  google_compute_network   = module.google-infra.compute_network
}

