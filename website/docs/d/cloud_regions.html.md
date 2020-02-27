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
  filter {
    name = "Name"
    values = ["US East.*"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A filter used to scope the list e.g. by tags.
  * `name` - (Required) The name of the filter. The valid values are defined in the [Pureport SDK Model](https://github.com/pureport/pureport-sdk-go/blob/develop/docs/client/CloudRegion.md).
  * `values` - (Required) The value of the filter. Currently only regex strings are supported.

## Attributes

* `regions` - The found list of regions.

    * `id` - The unique identifier for the cloud region.

    * `name` - The display name for the cloud region.

    * `provider` - The cloud provider for the cloud region.

    * `identifier` - The identifier provided by the cloud provider for this region.

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.

The Pureport Guide, []()
