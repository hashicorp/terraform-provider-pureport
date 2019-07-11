---
layout: "pureport"
page_title: "Pureport: pureport_cloud_regions"
sidebar_current: "docs-pureport-datasource-cloud_regions"
description: |-
  Provides details about existing Pureport cloud regions.
---

# Data Source: pureport\_cloud\_regions

## Example Usage

```hcl
data "pureport_cloud_regions" "name_regex" {
  name_regex = "US East.*"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to apply to the list of cloud regions returned by Pureport. This
  filtering is done locally on what Pureport returns, and could have a performance impact if the
  result is large.

## Attributes

* `regions` - The found list of regions.

    * `id` - The unique identifier for the cloud region.

    * `name` - The display name for the cloud region.

    * `provider` - The cloud provider for the cloud region.

    * `identifier` - The identifier provided by the cloud provider for this region.

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.

The Pureport Guide, []()
