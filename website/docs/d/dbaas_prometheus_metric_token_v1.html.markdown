---
layout: "selectel"
page_title: "Selectel: selectel_dbaas_prometheus_metric_token_v1"
sidebar_current: "docs-selectel-datasource-dbaas-prometheus-metric-token-v1"
description: |-
  Get information on Selectel DBaaS prometheus metrics tokens.
---

# selectel\_dbaas\_prometheus_metric_token_v1

Use this data source to get all available prometheus metrics tokens within Selectel DBaaS API Service

## Example Usage

```hcl
resource "selectel_vpc_project_v2" "project_1" {
  auto_quotas = true
}

data "selectel_dbaas_prometheus_metric_token_v1" "token" {
  project_id   = "${selectel_vpc_project_v2.project_1.id}"
  region       = "ru-3"
}
```

## Argument Reference

The folowing arguments are supported

* `project_id` - (Required) An associated Selectel VPC project.

* `region` - (Required) A Selectel VPC region.

## Attributes Reference

The following attributes are exported:

* `prometheus_metrics_tokens` - Contains a list of the found prometheus metrics tokens.

**prometheus_metrics_tokens**

- `id` - ID of the token.
- `created_at` - Create datetime of the token.
- `updated_at` - Update datetime of the token.
* `project_id` - Project ID associated with the token.
- `name` - Name of the token.
- `value` - Token's value.
