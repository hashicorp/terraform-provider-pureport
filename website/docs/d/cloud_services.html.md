---
layout: "pureport"
page_title: "Pureport: pureport_cloud_services"
sidebar_current: "docs-pureport-datasource-cloud_services"
description: |-
  Provides details about an existing Pureport cloud_services.
---

# Data Source: pureport\_cloud\_services

## Example Usage

```hcl
data "pureport_cloud_services" "name_regex" {
  name_regex = ".*S3 us-west-2"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to appy to the cloud services list returned by Pureport. This
  filtering is done locally on what Pureport returns, and could have a performance impact if the
  result is large.

## Attributes

* `services` - The found list of cloud provider services.

    * `id` - The unique identifier for the cloud service.

    * `name` - The display name for the cloud service.

    * `provider` - The cloud provider for the cloud service.

    * `href` - The unique path reference to the cloud service. This will be used by other resources to identify the service in most cases.

    * `ipv4_prefix_count` - The number of IPv4 prefixes associated with this cloud service.

    * `ipv6_prefix_count` - The number of IPv6 prefixes associated with this cloud service.

    * `cloud_region_id` - The identifier for the cloud service where this service is located.

The Pureport Guide, []()
