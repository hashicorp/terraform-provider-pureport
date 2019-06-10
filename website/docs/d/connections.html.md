---
layout: "pureport"
page_title: "Pureport: pureport_connections"
sidebar_current: "docs-pureport-datasource-connections"
description: |-
  Provides details about existing Pureport connections.
---

# Data Source: pureport\_connections

## Example Usage

```hcl
data "pureport_accounts" "main" {
  name_regex = "My Account.*"
}

data "pureport_networks" "main" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
  name_regex = "Connections"
}

data "pureport_connections" "empty" {
  network_href = "${data.pureport_networks.main.networks.0.href}"
}
```

## Argument Reference

The following arguments are supported:

* `network_href` - (Required) The HREF for the Pureport network associated with the connections.

- - -

* `name_regex` - (Optional) A regex string to apply to the network list returned by Pureport. This
  filtering is done locally on what Pureport returns, and could have a performance impact if the
  result is large.

## Attributes

* `connections` - A list of Pureport connections.

    * `id` - The unique identifier for the Pureport network.

    * `href` - The unique path reference for the Pureport connection.

    * `name` - The name of this connection.

    * `description` - The description for this connection.

    * `type` - The type of connection.

    * `speed` - The speed of this connection.

    * `location_href` - The HREF for the Pureport location associated with this connection.

    * `state` - The current state of this connection.


The Pureport Guide, []()
