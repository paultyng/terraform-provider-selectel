---
layout: "selectel"
page_title: "Selectel: selectel_dbaas_configuration_parameter_v1"
sidebar_current: "docs-selectel-datasource-dbaas-configuration-parameter-v1"
description: |-
  Get information on Selectel DBaaS configuration parameters.
---

# selectel\_dbaas\_configuration_parameter_v1

Use this data source to get all available confguration parameters within Selectel DBaaS API Service

## Example Usage

```hcl
resource "selectel_vpc_project_v2" "project_1" {
  auto_quotas = true
}

data "selectel_dbaas_datastore_type_v1" "dt" {
  project_id   = "${selectel_vpc_project_v2.project_1.id}"
  region       = "ru-3"
  filter {
    engine  = "postgresql"
    version = "12"
  }
}

data "selectel_dbaas_configuration_parameter_v1" "config" {
  project_id   = "${selectel_vpc_project_v2.project_1.id}"
  region       = "ru-3"
  filter {
    datastore_type_id = data.selectel_dbaas_datastore_type_v1.dt.datastore_types[0].id
    name = "work_mem"
  }
}
```

## Argument Reference

The folowing arguments are supported

* `project_id` - (Required) An associated Selectel VPC project.

* `region` - (Required) A Selectel VPC region.

* `filter` - (Optional) One or more values used to look up configuration parameters

**filter**

- `datastore_type_id` - (Optional) Datastore type id to lookup all available parameters for this type.
- `name` - (Optional) Name of the parameter to lookup.

## Attributes Reference

The following attributes are exported:

* `configuration_parameters` - Contains a list of the found configuration parameters.

**configuration_parameters**

- `id` - ID of the configuration parameter.
- `datastore_type_id` - Datastore type id for which the configuration parameter is availabe.
- `name` - Name of the configuration parameter.
- `type` - Type of the configuration parameter.
- `unit` - Unit of the configuration parameter. Might be empty.
- `min` - Min value of the configuration parameter. Might be empty.
- `max` - Max value of the configuration parameter. Might be empty.
- `default_value` - Default value of the configuration parameter. Might be empty.
- `choices` - Available choices for the configuration parameter. Some parameters have list of available options.
- `is_restart_required` - Shows if database needs a restart to apply changes of this parameter.
- `is_changeable` - Shows if parameter can be changed.
