---
layout: "pureport"
page_title: "Pureport: pureport_accounts"
sidebar_current: "docs-pureport-datasource-accounts"
description: |-
  Provides details about existing Pureport accounts.
---

# Data Source: pureport\_account

## Example Usage

```hcl
data "pureport_accounts" "empty" {
}

data "pureport_accounts" "name_regex" {
  name_regex = "My Name.*"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A filter used to scope the list e.g. by tags.
  * `name` - (Required) The name of the filter. The valid values are defined in the
    [Pureport SDK Model](https://github.com/pureport/pureport-sdk-go/blob/develop/docs/client/Account.md).
    Nested values are supported. E.g.("Location.DisplayName")
  * `values` - (Required) The value of the filter. Currently only regex strings are supported.

## Attributes

The Pureport Account resource exports the following attributes:

* `accounts` - The found list of accounts.

    * `id` - The unique identifier for the Pureport account.

    * `href` - The unique path reference to the Pureport account. This will be used by other resources to identify the account in most cases.

    * `name` - The name on the account.

    * `description` - The description of the account.

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.

The Pureport Guide, []()
