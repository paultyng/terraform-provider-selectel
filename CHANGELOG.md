## 3.7.0 (Unreleased)

FEATURES:
* Added `feature_gates` and `admission_controllers` to cluster kubernetes options ([#166](https://github.com/selectel/terraform-provider-selectel/issues/166))
* __New Data Source:__ `selectel_dbaas_configuration_parameter_v1` ([#162](https://github.com/selectel/terraform-provider-selectel/issues/162))
* __New Data Source:__ `selectel_mks_available_feature_gates_v1`
* __New Data Source:__ `selectel_mks_available_admission_controllers_v1`

IMPROVEMENTS:
* Fixed docs for DBaaS data sources ([#160](https://github.com/selectel/terraform-provider-selectel/pull/160))
* Updated Go version to `1.16` ([#164](https://github.com/selectel/terraform-provider-selectel/pull/164))
* Added `config` argument for the `selectel_dbaas_datastore_v1` resource ([#162](https://github.com/selectel/terraform-provider-selectel/issues/162))

## 3.6.2 (June 11, 2021)

IMPROVEMENTS:

* Change `connection` attribute type from TypeSet to TypeMap for the `selectel_dbaas_datastore_v1` resource ([#159](https://github.com/selectel/terraform-provider-selectel/pull/159))

## 3.6.1 (June 08, 2021)

IMPROVEMENTS:

* Add `filter` argument for the `selectel_dbaas_flavor_v1` data source ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* Add `connections` attribute for the `selectel_dbaas_datastore_v1` resource ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))

## 3.6.0 (June 03, 2021)

FEATURES:

* __New Resource:__ `selectel_dbaas_datastore_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Resource:__ `selectel_dbaas_user_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Resource:__ `selectel_dbaas_database_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Resource:__ `selectel_dbaas_grant_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Resource:__ `selectel_dbaas_extension_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Data Source:__ `selectel_dbaas_datastore_type_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Data Source:__ `selectel_dbaas_available_extension_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))
* __New Data Source:__ `selectel_dbaas_flavor_v1` ([#150](https://github.com/selectel/terraform-provider-selectel/issues/150))

IMPROVEMENTS:

* Updated Go version to `1.15` ([#148](https://github.com/selectel/terraform-provider-selectel/pull/148))
* Updated Terraform SDK to `v2.6.1` ([#149](https://github.com/selectel/terraform-provider-selectel/pull/149))

BUG FIXES:

* Fixed an issue with failing MKS acceptance tests when cluster maintenance could start during the tests ([#146](https://github.com/selectel/terraform-provider-selectel/issues/146))

## 3.5.0 (Nov 19, 2020)

FEATURES:

* __New Data Source:__ `selectel_domains_domain_v1` ([#128](https://github.com/selectel/terraform-provider-selectel/issues/128))

IMPROVEMENTS:

* Added `taints` argument for the `selectel_mks_nodegroup_v1` resource ([#130](https://github.com/selectel/terraform-provider-selectel/issues/130))
* Allowed using `ru-9` region as `region` argument for the `selectel_mks_cluster_v1` resource ([#141](https://github.com/selectel/terraform-provider-selectel/pull/141))
* Updated `domains-go` dependency to `v0.3.0` ([#142](https://github.com/selectel/terraform-provider-selectel/pull/142))

## 3.4.0 (Aug 28, 2020)

IMPROVEMENTS:

* Added `enable_pod_security_policy` argument for the `selectel_mks_cluster_v1` resource ([#115](https://github.com/selectel/terraform-provider-selectel/pull/115))
* Added `zonal` argument for the `selectel_mks_cluster_v1` resource ([#125](https://github.com/selectel/terraform-provider-selectel/pull/125))
* Upgraded Terraform SDK to `v2.0.1` ([#129](https://github.com/selectel/terraform-provider-selectel/pull/129))

## 3.3.3 (Aug 20, 2020)

IMPROVEMENTS:

* Updated Go version to `1.14.7` ([#117](https://github.com/selectel/terraform-provider-selectel/pull/117))
* Updated `golangci-lint` in CI to `v1.30.0` ([#118](https://github.com/selectel/terraform-provider-selectel/pull/118))
* Updated `go-selvpcclient` in CI to `v1.12.0` ([#118](https://github.com/selectel/terraform-provider-selectel/pull/118))

## 3.3.2 (July 13, 2020)

BUG FIXES:

* Fixed an issue when an emtpy value in the `quotas.resource_quotas.zone` argument of the `selectel_vpc_project_v2` resource caused Resell V2 API errors ([#111](https://github.com/terraform-providers/terraform-provider-selectel/issues/111))

## 3.3.1 (June 25, 2020)

BUG FIXES:

* Fixed an issue when the `selectel_mks_cluster_v1` resource was recreated on every Terraform run because of upper case usage in the `name` argument ([#103](https://github.com/terraform-providers/terraform-provider-selectel/issues/103))
* Fixed an issue when the `selectel_vpc_keypair_v2` resource was recreated on every Terraform run because the `regions` argument was computed from API responses ([#104](https://github.com/terraform-providers/terraform-provider-selectel/issues/104))
* Fixed unreadable error output for `selectel_mks_nodegroup_v1` resource creation failures ([#100](https://github.com/terraform-providers/terraform-provider-selectel/issues/100))
* Fixed an issue when an emtpy value in the `quotas.resource_quotas.region` argument of the `selectel_vpc_project_v2` resource caused Resell V2 API errors ([#107](https://github.com/terraform-providers/terraform-provider-selectel/issues/107))

## 3.3.0 (May 26, 2020)

FEATURES:

* __New Resource:__ `selectel_domains_domain_v1` ([#86](https://github.com/terraform-providers/terraform-provider-selectel/issues/86))
* __New Resource:__ `selectel_domains_record_v1` ([#91](https://github.com/terraform-providers/terraform-provider-selectel/issues/91))

IMPROVEMENTS:

* Added `labels` argument for `selectel_mks_nodegroup_v1` resource ([#88](https://github.com/terraform-providers/terraform-provider-selectel/issues/88))
* Added support to upgrade a minor version of Kubernetes cluster for `selectel_mks_cluster_v1` resource ([#93](https://github.com/terraform-providers/terraform-provider-selectel/issues/93))
* Tuned default timeouts for `selectel_mks_cluster_v1`, `selectel_mks_nodegroup_v1` resources ([#95](https://github.com/terraform-providers/terraform-provider-selectel/issues/95))
* Added retryable HTTP client to use with Selectel Domains API V1 client to make provider more resilient to 5xx errors ([#98](https://github.com/terraform-providers/terraform-provider-selectel/issues/98))
* Updated `terraform-plugin-sdk` to `v1.13.0` ([#96](https://github.com/terraform-providers/terraform-provider-selectel/issues/96))

BUG FIXES:

* Fixed an issue when `selectel_mks_cluster_v1` resource tried to downgrade its `kube_version` in case it was automatically upgraded in the MKS backend ([#85](https://github.com/terraform-providers/terraform-provider-selectel/issues/85))
* Fixed an issue when `nodes_count` attribute of `selectel_mks_nodegroup_v1` resource couldn't be imported ([#89](https://github.com/terraform-providers/terraform-provider-selectel/issues/89))

## 3.2.0 (April 29, 2020)

FEATURES:

* __New Resource:__ `selectel_mks_cluster_v1` ([#79](https://github.com/terraform-providers/terraform-provider-selectel/issues/79))
* __New Resource:__ `selectel_mks_nodegroup_v1` ([#80](https://github.com/terraform-providers/terraform-provider-selectel/issues/80))

IMPROVEMENTS:

* Updated Go version to `1.14.2` ([#83](https://github.com/terraform-providers/terraform-provider-selectel/issues/83))
* Updated `terraform-plugin-sdk` to `v1.10.0` ([#83](https://github.com/terraform-providers/terraform-provider-selectel/issues/83))
* Updated `golangci-lint` in CI to `v1.25.1` ([#83](https://github.com/terraform-providers/terraform-provider-selectel/issues/83))

## 3.1.0 (March 11, 2020)

IMPROVEMENTS:

* Added `network_id`, `subnet_id`, `port_id` attributes into `selectel_vpc_license_v2` resource ([#78](https://github.com/terraform-providers/terraform-provider-selectel/issues/78))
* Updated `terraform-plugin-sdk` to `v1.7.0` ([#76](https://github.com/terraform-providers/terraform-provider-selectel/issues/76))
* Updated `golangci-lint` in CI to `v1.23.8` ([#77](https://github.com/terraform-providers/terraform-provider-selectel/issues/77))

## 3.0.0 (September 26, 2019)

BREAKING CHANGES:

* Removed `selectel_vpc_crossregion_subnet_v2` resource because it has been deprecated in the Selectel VPC V2 API ([#69](https://github.com/terraform-providers/terraform-provider-selectel/issues/69))

IMPROVEMENTS:

* Added ability to revoke tokens in API while deleting `selectel_vpc_project_v2` resource ([#66](https://github.com/terraform-providers/terraform-provider-selectel/issues/66))
* Added ability to import `selectel_vpc_user_v2` resource ([#65](https://github.com/terraform-providers/terraform-provider-selectel/issues/65))
* Added extended error messages to all resources ([#73](https://github.com/terraform-providers/terraform-provider-selectel/issues/73))
* Migrated from Terraform in-tree `helper/*` SDK to the separate `terraform-plugin-sdk` v1.0.0 ([#72](https://github.com/terraform-providers/terraform-provider-selectel/issues/72))

BUG FIXES:

* Fixed an issue where checks for 404 errors could cause panics ([#73](https://github.com/terraform-providers/terraform-provider-selectel/issues/73))

## 2.3.0 (July 09, 2019)

BUG FIXES:

* Fixed an issue with `selectel_vpc_project_v2` when `quotas` argument has been updated incorrectly ([#64](https://github.com/terraform-providers/terraform-provider-selectel/issues/64))

IMPROVEMENTS:

* Updated Terraform SDK to `v1.12.2` from `v1.12.0` ([#61](https://github.com/terraform-providers/terraform-provider-selectel/issues/61))
* Updated `golangci-lint` in CI to `v1.17.1` ([#63](https://github.com/terraform-providers/terraform-provider-selectel/issues/63))
* Fixed Terraform and Go versions in documentation ([#63](https://github.com/terraform-providers/terraform-provider-selectel/issues/63))

## 2.2.0 (May 23, 2019)

IMPROVEMENTS:

* Updated Terraform SDK to `v1.12.0` from `v1.12.0-beta1` ([#58](https://github.com/terraform-providers/terraform-provider-selectel/issues/58))
* Updated `golangci-lint` in CI to `v1.16.0` ([#55](https://github.com/terraform-providers/terraform-provider-selectel/issues/55))

## 2.1.0 (March 14, 2019)

BUG FIXES:

* Fixed an issue with empty `project_id` argument of the `selectel_vpc_crossregion_subnet_v2` resource ([#52](https://github.com/terraform-providers/terraform-provider-selectel/issues/52))

IMPROVEMENTS:

* Migrated to Go Modules ([#47](https://github.com/terraform-providers/terraform-provider-selectel/issues/47))
* Updated Terraform SDK to `v1.12.0-beta1` ([#51](https://github.com/terraform-providers/terraform-provider-selectel/issues/51))
* Updated `golangci-lint` in CI to `v1.15.0` ([#54](https://github.com/terraform-providers/terraform-provider-selectel/issues/54))

## 2.0.0 (February 04, 2019)

BREAKING CHANGES:

* All `selvpc_resell_*` resources were renamed to `selectel_vpc_*` resources ([#45](https://github.com/terraform-providers/terraform-provider-selectel/issues/45))

FEATURES:

* __New Resource:__ `selectel_vpc_crossregion_subnet_v2` ([#43](https://github.com/terraform-providers/terraform-provider-selectel/issues/43))

BUG FIXES:

* Fixed VPC V2 Token Account acceptance test ([#41](https://github.com/terraform-providers/terraform-provider-selectel/issues/41))

## 1.1.0 (January 08, 2019)

FEATURES:

* __New Resource:__ `selvpc_resell_keypair_v2` ([#29](https://github.com/terraform-providers/terraform-provider-selectel/issues/29))
* __New Resource:__ `selvpc_resell_vrrp_subnet_v2` ([#35](https://github.com/terraform-providers/terraform-provider-selectel/issues/35))

IMPROVEMENTS:

* Added tuned HTTP client to prevent errors when making call to the Resell API ([#30](https://github.com/terraform-providers/terraform-provider-selectel/issues/30))
* Added the same format for all debug messages ([#32](https://github.com/terraform-providers/terraform-provider-selectel/issues/32))
* Remove the `type` argument of the `selvpc_resell_subnet_v2` from the documentation as it doesn't exist ([#36](https://github.com/terraform-providers/terraform-provider-selectel/issues/36))
* Updated Go-selvpcclient dependency to `v1.6.0` ([#33](https://github.com/terraform-providers/terraform-provider-selectel/issues/33))
* Used `v1.11.x` Go version in Travis CI ([#40](https://github.com/terraform-providers/terraform-provider-selectel/issues/40))
* Updated GolangCI-Lint in Travis CI to `v1.12.5` ([#37](https://github.com/terraform-providers/terraform-provider-selectel/issues/37))

## 1.0.0 (December 19, 2018)

FEATURES:

* __New Resource:__ `selvpc_resell_role_v2` ([#4](https://github.com/terraform-providers/terraform-provider-selectel/issues/4))
* __New Resource:__ `selvpc_resell_subnet_v2` ([#1](https://github.com/terraform-providers/terraform-provider-selectel/issues/1))
* __New Resource:__ `selvpc_resell_token_v2` ([#2](https://github.com/terraform-providers/terraform-provider-selectel/issues/2))
* __New Resource:__ `selvpc_resell_user_v2` ([#3](https://github.com/terraform-providers/terraform-provider-selectel/issues/3))

IMPROVEMENTS:

* Updated `Building The Provider` and `Using the provider` sections in the Readme ([#6](https://github.com/terraform-providers/terraform-provider-selectel/issues/6))
* Added `GolangCI-Lint` in the `TravisCI`, removed separated linters scripts and cleaned up `GNUmakefile` ([#12](https://github.com/terraform-providers/terraform-provider-selectel/issues/12))
* Added more context into error messages ([#17](https://github.com/terraform-providers/terraform-provider-selectel/issues/17))
* Added tuned HTTP timeouts instead of the default ones from Go's `net/http` package ([#14](https://github.com/terraform-providers/terraform-provider-selectel/issues/14))
* Updated `go-selvpcclient` dependency to `v1.5.0` ([#14](https://github.com/terraform-providers/terraform-provider-selectel/issues/14))

## 0.3.0 (November 26, 2018)

IMPROVEMENTS:

* Updated `go-selvpcclient` dependency to `v1.4.0` ([#51](https://github.com/selectel/terraform-provider-selvpc/issues/51))
* Updated documentation for `floatingip_v2`, `license_v2` and `project_v2` resources ([#50](https://github.com/selectel/terraform-provider-selvpc/issues/50))
* Changed `TypeList` to `TypeSet` for the `servers`, `quotas`, `all_quotas`, `resource_quotas` attributes ([#48](https://github.com/selectel/terraform-provider-selvpc/issues/48))
* Added a check for error on setting non-scalars ([#52](https://github.com/selectel/terraform-provider-selvpc/issues/52))
* Added a check for if resources don’t exist during read with unsetting the ID ([#53](https://github.com/selectel/terraform-provider-selvpc/issues/53))
* Grouped attributes at the top of resources followed by the optional attributes ([#54](https://github.com/selectel/terraform-provider-selvpc/issues/54)) 

BUG FIXES: 

* Fixed `golint` URL in the TravisCI configuration ([#49](https://github.com/selectel/terraform-provider-selvpc/issues/49))
* Fixed `all_quotas` attribute checking in the `TestAccResellV2ProjectAutoQuotas` ([#57](https://github.com/selectel/terraform-provider-selvpc/issues/57)), ([#62](https://github.com/selectel/terraform-provider-selvpc/issues/62))
* Fixed quotas in the created project of the `selvpc_resell_floatingip_v2` resource ([#58](https://github.com/selectel/terraform-provider-selvpc/issues/58))
* Fixed `structLitKeyOrder` errors in the CI ([#60](https://github.com/selectel/terraform-provider-selvpc/issues/60))

## 0.2.0 (Oct 3, 2018)

FEATURES:

* Added `auto_quotas` attribute for the `selvpc_resell_project_v` resource ([#41](https://github.com/selectel/terraform-provider-selvpc/issues/41))

IMPROVEMENTS:

* Added `critic` target in the `GNUmakefile` that will run `gocritic` linter. This target will be called by the Travis CI ([#43](https://github.com/selectel/terraform-provider-selvpc/issues/43))
* Updated Go version to the `1.11.1` in the Travis CI configuration ([#44](https://github.com/selectel/terraform-provider-selvpc/issues/44))

## 0.1.0 (May 13, 2018)

FEATURES:

* __New Resource:__ `selvpc_resell_project_v2` ([#3](https://github.com/selectel/terraform-provider-selvpc/issues/3))
* __New Resource:__ `selvpc_resell_floatingip_v2` ([#34](https://github.com/selectel/terraform-provider-selvpc/issues/34))
* __New Resource:__ `selvpc_resell_license_v2` ([#33](https://github.com/selectel/terraform-provider-selvpc/issues/33))
