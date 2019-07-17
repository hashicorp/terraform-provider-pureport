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

* `filter` - (Optional) A filter used to scope the list e.g. by tags.
  * `name` - (Required) The name of the filter. The valid values are defined in the [Pureport SDK Model](https://github.com/pureport/pureport-sdk-go/blob/develop/docs/client/Connection.md).
  * `values` - (Required) The value of the filter. Currently only regex strings are supported.

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

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.


The Pureport Guide, []()
