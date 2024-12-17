## 2.4.1 (Dec 17, 2024)

IMPROVEMENTS:

* Append outputs of std.trace() function to the error diagnostics.

## 2.4.0 (Dec 16, 2024)

IMPROVEMENTS:

* Added the `trace` attribute containing the outputs of std.trace() function (fixed #13).

## 2.3.2 (Oct 4, 2024)

IMPROVEMENTS:

* Added ability to import libraries from jsonnet snippets provided by `content` attribute (fixes #12).

NOTES:

* Used go-1.23.1 to build provider.

## 2.3.1 (Nov 8, 2023)

BUG FIXES:

* Fixed processing of `jsonnet_path` provider attribute (fixes #9).

## 2.3.0 (Nov 4, 2023)

IMPROVEMENTS:

* Added `jsonnet_path` attribute of data source.
* Upgraded go-jsonnet to 0.20.0

NOTES:

* Used go-1.21 to build provider.

## 2.2.0 (Mar 28, 2023)

IMPROVEMENTS:

* Added `string_output` attribute
* Upgraded go-jsonnet to 0.19.1

NOTES:

* Used go-1.19 to build provider

## 2.1.0 (Mar 23, 2022)

IMPROVEMENTS:

* Upgraded go-jsonnet to 0.17.0
* Used go-1.17 to build provider

## 2.0.0 (Dec 12, 2021)

BREAKING CHANGES:

* Provider migrated to Terraform Plugin SDKv2.
* Type of jsonnet_path changed from list to string. 

IMPROVEMENTS:

* Provider can be configured by `JSONNET_PATH` environment variable.

## 1.0.3 (Dec 25, 2020)

BUG FIXES:

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

BUG FIXES:

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
