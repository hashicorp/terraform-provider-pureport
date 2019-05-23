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

* `name_regex` - (Optional) A regex string to apply to the account list returned by Pureport. This
  filtering is done locally on what Pureport returns, and could have a performance impact if the
  result is large.

## Attributes

The Pureport Account resource exports the following attributes:

* `accounts` - The found list of accounts.

    * `id` - The unique identifier for the Pureport account.

    * `href` - The unique path reference to the Pureport account. This will be used by other resources to identify the account in most cases.

    * `name` - The name on the account.

    * `description` - The description of the account.

The Pureport Guide, []()
