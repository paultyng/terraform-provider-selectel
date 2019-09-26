## 3.0.0 (Unreleased)

BREAKING CHANGES:

* Removed `selectel_vpc_crossregion_subnet_v2` resource because it has been depreacted in the Selectel VPC V2 API

IMPROVEMENTS:

* Added ability to revoke tokens in API while deleting `selvpc_resell_token_v2` resource [GH-66]
* Added ability to import `selvpc_resell_user_v2` resource [GH-65]

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
