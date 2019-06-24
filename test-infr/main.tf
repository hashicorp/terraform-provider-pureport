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
