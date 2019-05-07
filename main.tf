data "pureport_accounts" "main" {
  name_regex = "Terraform"
}

data "pureport_cloud_regions" "main" {
  name_regex = "Oregon"
}

data "pureport_locations" "main" {
  name_regex = "^Ral*"
}

data "pureport_networks" "main" {
  account_id = "${data.pureport_accounts.main.accounts.0.id}"
  name_regex = "Bansh.*"
}

resource "pureport_site_vpn_connection" "main" {
  name              = "SiteVPNTest"
  speed             = "100"
  high_availability = true

  location {
    id   = "${data.pureport_locations.main.locations.0.id}"
    href = "${data.pureport_locations.main.locations.0.href}"
  }

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }

  ike_version  = "V2"
  routing_type = "ROUTE_BASED_BGP"
  customer_asn = 30000

  primary_customer_router_ip   = "123.123.123.123"
  secondary_customer_router_ip = "124.124.124.124"
}
