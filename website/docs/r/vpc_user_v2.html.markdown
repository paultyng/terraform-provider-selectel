---
layout: "selectel"
page_title: "Selectel: selectel_vpc_user_v2"
sidebar_current: "docs-selectel-resource-vpc-user-v2"
description: |-
  Creates and manages a service user for Selectel products using public API v2.
---

# selectel\_vpc\_user_v2

> **WARNING**: This resource is deprecated. Since version 5.0.0, replace the resource with the roles block in the [selectel_iam_serviceuser_v1](https://registry.terraform.io/providers/selectel/selectel/latest/docs/resources/iam_serviceuser_v1) resource. For more information about upgrading to version 5.0.0, see the [upgrading guide](https://registry.terraform.io/providers/selectel/selectel/latest/docs/guides/upgrading_to_version_5).

Creates and manages a service user for Selectel products using public API v2. Selectel products support Identity and Access Management (IAM). For more information about service users, see the [official Selectel documentation](https://docs.selectel.ru/control-panel-actions/users-and-roles/user-types-and-roles/).

When you create users, they do not have any roles. To grant a role, use the [selectel_vpc_role_v2](https://registry.terraform.io/providers/selectel/selectel/latest/docs/resources/vpc_role_v2) resource.

~> **Note:** The user password is stored as raw data in a plain-text file. Learn more about [sensitive data in
state](https://developer.hashicorp.com/terraform/language/state/sensitive-data).

## Example Usage

```hcl
resource "selectel_vpc_user_v2" "user_1" {
  name     = "username"
  password = "verysecret"
  enabled  = true
}
```

## Argument Reference

* `name` - (Required) Name of the service user. Changing this updates the name of the existing user.

* `password` - (Required, Sensitive) Password of the service user. Changing this updates the password of the existing user.

* `enabled` - (Optional) Specifies if you can create a Cloud Platform Keystone token for the user. Boolean flag, the default value is `true`. Learn more about [Cloud Platform Keystone tokens](https://developers.selectel.ru/docs/control-panel/authorization/).
