---
layout: "selectel"
page_title: "Selectel: selectel_dbaas_prometheus_metric_token_v1"
sidebar_current: "docs-selectel-resource-dbaas-prometheus-metric-token-v1"
description: |-
  Manages a V1 prometheus metrics tokens resource within Selectel Managed Databases Service.
---

# selectel\_dbaas\_prometheus_metric_token_v1

Manages a V1 prometheus metrics tokens resource within Selectel Managed Databases Service.

## Example Usage

```hcl
resource "selectel_vpc_project_v2" "project_1" {
  auto_quotas = true
}

resource "selectel_dbaas_prometheus_metric_token_v1" "token" {
  project_id   = "${selectel_vpc_project_v2.project_1.id}"
  region       = "ru-3"
  name         = "token"
}
```

## Argument Reference

The folowing arguments are supported

* `project_id` - (Required) An associated Selectel VPC project.

* `region` - (Required) A Selectel VPC region.

* `name` - (Required) A name of the token.

## Attributes Reference

The following attributes are exported:

* `value` - value of the token.

## Import

Prometheus metrics token can be imported using the `id`, e.g.

```shell
$ env SEL_TOKEN=SELECTEL_API_TOKEN SEL_PROJECT_ID=SELECTEL_VPC_PROJECT_ID SEL_REGION=SELECTEL_VPC_REGION terraform import selectel_dbaas_prometheus_metric_token_v1.token b311ce58-2658-46b5-b733-7a0f418703f2
```
