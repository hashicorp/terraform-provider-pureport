---
layout: "pureport"
page_title: "Pureport: pureport_networks"
sidebar_current: "docs-pureport-datasource-networks"
description: |-
  Provides details about existing Pureport networks.
---

# Data Source: pureport\_networks

## Example Usage

```hcl
data "pureport_accounts" "main" {
  name_regex = "My Account.*"
}

data "pureport_networks" "empty" {
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
}
```

## Argument Reference

The following arguments are supported:

* `account_href` - (Required) The HREF for the Pureport account associated with this network.

- - -

* `filter` - (Optional) A filter used to scope the list e.g. by tags.
  * `name` - (Required) The name of the filter. The valid values are defined in the [Pureport SDK Model](https://github.com/pureport/pureport-sdk-go/blob/develop/docs/client/Network.md).
  * `values` - (Required) The value of the filter. Currently only regex strings are supported.

## Attributes

* `networks` - A list of Pureport networks.

    * `id` - The unique identifier for the Pureport network.

    * `href` - The unique path reference for the Pureport network. This will be used by other resources to identify the locations in most cases.

    * `name` - The name of the network.

    * `description` - The description for the network.

    * `account_href` - The HREF for the Pureport account associated with this network.

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.


The Pureport Guide, []()
