---
layout: "pureport"
page_title: "Pureport: pureport_network"
sidebar_current: "docs-pureport-resource-network"
description: |-
  Manages a Pureport Network.
---

# Resource: pureport\_network

## Example Usage

```hcl

data "pureport_accounts" "main" {
  name_regex = "MyAccount"
}

resource "pureport_network" "main" {
  name = "MyNetwork"
  description = "My Custom Network"
  account_href = "${data.pureport_accounts.main.accounts.0.href}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name used for the Network.
* `account_href` - (Required) HREF for the Account associated with the Network.

- - -

* `description` - (Optional) The description for the Network.

## Attributes

* `href` - The HREF to reference this Network.

The Pureport Guide, []()
