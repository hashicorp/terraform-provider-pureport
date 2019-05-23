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

* `name_regex` - (Optional) A regex string to apply to the network list returned by Pureport. This
  filtering is done locally on what Pureport returns, and could have a performance impact if the
  result is large.

## Attributes

* `networks` - A list of Pureport networks.

    * `id` - The unique identifier for the Pureport network.

    * `href` - The unique path reference for the Pureport network. This will be used by other resources to identify the locations in most cases.

    * `name` - The name of the network.

    * `description` - The description for the network.

    * `account_href` - The HREF for the Pureport account associated with this network.


The Pureport Guide, []()
