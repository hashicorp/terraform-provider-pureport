---
layout: "pureport"
page_title: "Provider: Pureport Cloud Platform"
sidebar_current: "docs-pureport-index"
description: |-
   The Pureport provider is used to configure your Pureport Cloud Platform infrastructure
---

# Pureport Cloud Platform Provider

## Example Usage

```hcl
# Configure the Linode provider
provider "pureport" {
  api_key = "$SOME_KEY"
  api_secret = "$SOME_SECRET"
}

resource "pureport_account" "foobar" {
  # ...
}
```

## Configuration Reference

The following keys can be used to configure the provider.

* `api_key` - (Optional) The Pureport API Key.

* `api_secret` - (Optional) The Pureport API Secret. This is required when the `api_key` is specified.

* `api_url` - (Optional) The Pureport REST API URL. (default: https://api.pureport.com)

* `auth_profile` - (Optional) If you are using Pureport configuration files for authentication, you can use this to specified the profile that should be used to read the API Key and Secret.

The values above can also be configured via the Environment variables below:

* PUREPORT_API_KEY
* PUREPORT_API_SECRET
* PUREPORT_ENDPOINT
* PUREPORT_PROFILE

## Pureport Guides

## Debugging

You can use the standard Terraform `TF_LOG` levels to configure the debug logging output by this
provider.
