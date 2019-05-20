---
layout: "pureport"
page_title: "Provider: Pureport Cloud Platform"
sidebar_current: "docs-pureport-provider-x"
description: |-
   The Pureport provider is used to configure your Pureport Cloud Platform infrastructure
---

# Pureport Cloud Platform Provider

## Example Usage

```hcl
# Configure the Linode provider
provider "pureport" {
  token = "$LINODE_TOKEN"
}

resource "pureport_account" "foobar" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

## Pureport Guides

## Debugging
