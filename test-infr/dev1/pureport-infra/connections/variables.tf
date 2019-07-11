variable "connections_network_href" {
  description = "The HREF for the Connections network that we will use for testing the pureport_connections data source."
}

variable "datasource_network_href" {
  description = "The HREF for the Pureport Network we should use for deploying connections that will be used during data source testing."
}

variable "datasource_express_route" {
  description = "The reference to the Azure Express Route Circuit data source we should use for deploying the Azure Connection."
}

variable "google_compute_network" {
  description = "The reference to the Google Compute Network data source we should use for deploying the Google Cloud Connection."
}
