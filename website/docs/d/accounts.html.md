---
layout: "pureport"
page_title: "Pureport: pureport_account"
sidebar_current: "docs-pureport-datasource-account"
description: |-
  Provides details about an existing Pureport account.
---

# Data Source: pureport\_account

## Example Usage

```hcl
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to appy to the account list returned by Pureport. This
  filtering is done locally on what Pureport returns, and could have a performance impact if the
  result is large.

## Attributes

The Purport Account resource exports the following attributes:

* `accounts` - The found list of accounts.

    * `id` - The unique identifier for the Pureport account.

    * `href` - The unique path reference to the Pureport account. This will be used by other resources to identify the account in most cases.

    * `name` - The name on the account.

    * `description` - The description of the account.

The Pureport Guide, []()
