---
layout: "pureport"
page_title: "Pureport: pureport_cloud_services"
sidebar_current: "docs-pureport-datasource-cloud_services"
description: |-
  Provides details about existing Pureport cloud services.
---

# Data Source: pureport\_cloud\_services

## Example Usage

```hcl
data "pureport_cloud_services" "name_regex" {
  filter {
    name = "Name"
    values = [".*S3 us-west-2"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A filter used to scope the list e.g. by tags.
  * `name` - (Required) The name of the filter. The valid values are defined in the [Pureport SDK Model](https://github.com/pureport/pureport-sdk-go/blob/develop/docs/client/CloudService.md).
  * `values` - (Required) The value of the filter. Currently only regex strings are supported.

## Attributes

* `services` - The found list of cloud provider services.

    * `id` - The unique identifier for the cloud service.

    * `name` - The display name for the cloud service.

    * `provider` - The cloud provider for the cloud service.

    * `href` - The unique path reference to the cloud service. This will be used by other resources to identify the service in most cases.

    * `ipv4_prefix_count` - The number of IPv4 prefixes associated with this cloud service.

    * `ipv6_prefix_count` - The number of IPv6 prefixes associated with this cloud service.

    * `cloud_region_id` - The identifier for the cloud service where this service is located.

    * `tags` - A dictionary of user defined key/value pairs associated with this resource.

The Pureport Guide, []()
