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

data "pureport_cloud_regions" "main" {
  name_regex = "Los.*"
}

data "pureport_locations" "main" {
  name_regex = "^Sea.*"
}

data "pureport_networks" "main" {
  account_id = "${data.pureport_accounts.main.accounts.0.id}"
  name_regex = "Bansh.*"
}

resource "pureport_google_cloud_connection" "main" {
  name  = "GoogleCloudTest"
  speed = "50"

  location_href = "${data.pureport_locations.main.locations.0.href}"

  network {
    id   = "${data.pureport_networks.main.networks.0.id}"
    href = "${data.pureport_networks.main.networks.0.href}"
  }

  primary_pairing_key = "${google_compute_interconnect_attachment.main.0.pairing_key}"
}
