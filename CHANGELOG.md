## 2.0.1 (Mar 23, 2022)

NOTES:

* Upgrade go to 1.17
* Upgrade go-jsonnet to 0.18.0

## 2.0.0 (Dec 12, 2021)

BREAKING CHANGES:

* Provider migrated to Terraform Plugin SDKv2.
* Type of jsonnet_path changed from list to string. 

IMPROVEMENTS:

* Provider can be configured by `JSONNET_PATH` environment variable.

## 1.0.3 (Dec 25, 2020)

BUG FIXED:

Fixed concurrency issues caused by shared FileImporter.

## 1.0.2 (Dec 23, 2020)

NOTES:

Upgrade go-jsonnet to 0.17.0 

## 1.0.1 (Jun 30, 2020)

NOTES:

Added provider documentation to be published on Terraform Registry.

## 1.0.0 (Jun 14, 2020)

BREAKING CHANGES:

`go-jsonnet` is used for generation json instead of `jsonnet` utility, therefore
the parameter `jsonnet_bin` of provider was removed.

## 0.2.2 (Mar 14, 2020)

BUG FIXED:

Usage of external variables don't distort provider config.

## 0.2.1 (Mar 14, 2020)

NOTES:

* Added tests
* Release script counts checksums of binaries instead of archives 

## 0.2.0 (Mar 12, 2020)

FEATURES:

Supported external variables and top-level arguments

IMPROVEMENTS:

Output full stderr on failure of the command

## 0.1.0 (Mar 10, 2020)

NOTES:

First public release supporting rendering of jsonnet templates using configurable paths to jsonnet binary and libraries.
