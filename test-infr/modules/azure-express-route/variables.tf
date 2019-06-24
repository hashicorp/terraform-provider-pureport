
variable "peering_location" {
  type = map(string)

  default = {
    "westus2" = "Seattle"
    "eastus"  = "Washington DC"
    "eastus2" = "Washington DC"
  }
}

variable "resource_group_name" {
}
