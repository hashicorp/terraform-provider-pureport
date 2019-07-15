---
layout: "pureport"
page_title: "Pureport: pureport_locations"
sidebar_current: "docs-pureport-datasource-locations"
description: |-
  Provides details about existing Pureport locations.
---

# Data Source: pureport\_locations

## Example Usage

```hcl
data "pureport_locations" "name_regex" {
  name_regex = "^Sea*"
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A filter used to scope the list e.g. by tags.
  * `name` - (Required) The name of the filter. The valid values are defined in the [Pureport SDK Model](https://github.com/pureport/pureport-sdk-go/blob/develop/docs/client/Location.md).
  * `values` - (Required) The value of the filter. Currently only regex strings are supported.

## Attributes

* `locations` - A list of Pureport locations.

    * `id` - The unique identifier for the Pureport locations.

    * `href` - The unique path reference for the Pureport locations. This will be used by other resources to identify the locations in most cases.

    * `name` - The name of the location.

    * `links` - The available links to other Pureport locations.

        * `location_href` - The href of the linked location.

        * `speed` - The connection speed between the locations.

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.

The Pureport Guide, []()
