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
  filter {
    name = "Name"
    values = ["MyAccount"]
  }
}

resource "pureport_network" "main" {
  name = "MyNetwork"
  description = "My Custom Network"
  account_href = data.pureport_accounts.main.accounts.0.href

  tags = {
    Environment = "production"
    Owner       = "Scott Pilgrim"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name used for the Network.
* `account_href` - (Required) HREF for the Account associated with the Network.

- - -

* `description` - (Optional) The description for the Network.
* `tags` - (Optional) A dictionary of user defined key/value pairs to associate with this resource.

## Attributes

* `href` - The HREF to reference this Network.

The Pureport Guide, []()
