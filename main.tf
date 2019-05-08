data "google_compute_network" "default" {
  name = "default"
}

resource "google_compute_router" "main" {
  name    = "terraform-acc-router-${count.index + 1}"
  network = "${data.google_compute_network.default.name}"

  bgp {
    asn = "16550"
  }

  count = 2
}

resource "google_compute_interconnect_attachment" "main" {
  name                     = "terraform-acc-interconnect-${count.index + 1}"
  router                   = "${element(google_compute_router.main.*.self_link, count.index)}"
  type                     = "PARTNER"
  edge_availability_domain = "AVAILABILITY_DOMAIN_${count.index + 1}"

  count = 2
}

data "pureport_accounts" "main" {
  name_regex = "Terraform"
}

data "pureport_locations" "seattle" {
  name_regex = "^Sea.*"
}

data "pureport_networks" "main" {
  account_id = "${data.pureport_accounts.main.accounts.0.id}"
  name_regex = "Bansh.*"
}

data "pureport_cloud_regions" "oregon" {
  name_regex = "Oregon"
}

resource "pureport_google_cloud_connection" "main" {
  name  = "GoogleCloudTest"
  speed = "50"

  location_href = "${data.pureport_locations.seattle.locations.0.href}"

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }

  primary_pairing_key = "${google_compute_interconnect_attachment.main.0.pairing_key}"
}

resource "pureport_aws_connection" "main" {
  name              = "AwsDirectConnectTest"
  speed             = "100"
  high_availability = true

  location_href = "${data.pureport_locations.seattle.locations.0.href}"

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }

  aws_region     = "${data.pureport_cloud_regions.oregon.regions.0.identifier}"
  aws_account_id = "123456789012"
}

resource "pureport_azure_connection" "main" {
  name              = "AzureExpressRouteTest"
  description       = "Some random description"
  speed             = "100"
  high_availability = true

  location_href = "${data.pureport_locations.seattle.locations.0.href}"

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }

  service_key = "3166c9a8-1275-4e7b-bad2-0dc6db0c6e02"
}

resource "pureport_dummy_connection" "main" {
  name              = "DummyTest"
  speed             = "100"
  high_availability = true

  location_href = "${data.pureport_locations.seattle.locations.0.href}"

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }
}

resource "pureport_site_vpn_connection" "main" {
  name              = "SiteVPNTest"
  speed             = "100"
  high_availability = true

  location_href = "${data.pureport_locations.seattle.locations.0.href}"

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }

  ike_version = "V2"

  routing_type = "ROUTE_BASED_BGP"
  customer_asn = 30000

  primary_customer_router_ip   = "123.123.123.123"
  secondary_customer_router_ip = "124.124.124.124"
}

resource "pureport_network" "main" {
  name        = "NetworkTest"
  description = "Network Terraform Test"
  account_id  = "${data.pureport_accounts.main.accounts.0.id}"
}
