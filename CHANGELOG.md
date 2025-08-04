# Changelog

## [2.2.1](https://github.com/promhippie/scw_exporter/compare/v2.2.0...v2.2.1) (2025-08-04)


### Bugfixes

* **deps:** update module github.com/golangci/golangci-lint/v2 to v2.3.1 ([d659ac3](https://github.com/promhippie/scw_exporter/commit/d659ac37a310826e4a269e3a71a48487684dbd7e))
* **deps:** update module github.com/prometheus/client_golang to v1.23.0 ([f95115f](https://github.com/promhippie/scw_exporter/commit/f95115f00dc9d9e88dcf50831a910d6911ec6560))

## [2.2.0](https://github.com/promhippie/scw_exporter/compare/v2.1.0...v2.2.0) (2025-07-28)


### Features

* **deps:** update module github.com/golangci/golangci-lint/v2 to v2.3.0 ([eb92e5c](https://github.com/promhippie/scw_exporter/commit/eb92e5c87d6e0d0781e0d5d6104ba60dc9c1790a))


### Miscellaneous

* **deps:** update docker digests ([9c9c112](https://github.com/promhippie/scw_exporter/commit/9c9c1124f0484e3538f44c2123d8f5fe090097cb))
* **deps:** update golang:1.24.5-alpine3.21 docker digest to 6edc205 ([ecdea64](https://github.com/promhippie/scw_exporter/commit/ecdea6450354fd4b07a34dd0903dfdd780812bc4))
* **deps:** update golang:1.24.5-alpine3.21 docker digest to 72ff633 ([8613da5](https://github.com/promhippie/scw_exporter/commit/8613da547657986dca7d77fbf0cdabde93477967))

## [2.1.0](https://github.com/promhippie/scw_exporter/compare/v2.0.0...v2.1.0) (2025-07-14)


### Features

* **deps:** update module github.com/mgechev/revive to v1.11.0 ([c932188](https://github.com/promhippie/scw_exporter/commit/c93218820d748373496d95eaca7807728a953e7b))


### Bugfixes

* **deps:** update golang docker tag to v1.24.5 ([fb09cc5](https://github.com/promhippie/scw_exporter/commit/fb09cc5f0d545b4104ee51732be4f3b5b16f8245))
* **deps:** update module github.com/golangci/golangci-lint/v2 to v2.2.2 ([fc72fd7](https://github.com/promhippie/scw_exporter/commit/fc72fd769b3a6f1dba18fc8370a7942ab967227e))

## [2.0.0](https://github.com/promhippie/scw_exporter/compare/v1.2.1...v2.0.0) (2025-07-09)


### ⚠ BREAKING CHANGES

* rename project label and add project name labels
* restructure build and release process

### Features

* **deps:** update module github.com/oklog/run to v1.2.0 ([#142](https://github.com/promhippie/scw_exporter/issues/142)) ([d12d784](https://github.com/promhippie/scw_exporter/commit/d12d78436d8033312e00d1c5842ec8b56cba70f7))
* rename project label and add project name labels ([b4228db](https://github.com/promhippie/scw_exporter/commit/b4228db64a9001064713ff5d0146abe992a02f5d))
* restructure build and release process ([3dc26ed](https://github.com/promhippie/scw_exporter/commit/3dc26ed88a59e6ccf8cfd62ed2d7eca493becfbd))


### Bugfixes

* **deps:** update module github.com/go-chi/chi/v5 to v5.2.2 ([#141](https://github.com/promhippie/scw_exporter/issues/141)) ([5d65abc](https://github.com/promhippie/scw_exporter/commit/5d65abc0e9261f10728f93bf373d33fd47ca2d7a))
* **deps:** update module github.com/scaleway/scaleway-sdk-go to v1.0.0-beta.34 ([47691fa](https://github.com/promhippie/scw_exporter/commit/47691fa65e9adcb3921407f8c331b03a4c6d28b2))


### Miscellaneous

* **deps:** pin dependencies ([4c8f517](https://github.com/promhippie/scw_exporter/commit/4c8f5175e6230b755abe372f249f2b9ba195a85e))
* enable logging for watch task ([44c9e4b](https://github.com/promhippie/scw_exporter/commit/44c9e4bde75e9f33a840d62beb3627564af5fc4b))
* **flake:** updated lockfile [skip ci] ([e6fc8f7](https://github.com/promhippie/scw_exporter/commit/e6fc8f7d53cc412fe39ab81ca0cd3a1cf7934ba1))
* **flake:** updated lockfile [skip ci] ([74b5922](https://github.com/promhippie/scw_exporter/commit/74b5922d4eb73cd205fdb0c3c1950035ff6b8ae5))
* **flake:** updated lockfile [skip ci] ([14b2657](https://github.com/promhippie/scw_exporter/commit/14b2657255536accaf87db66087a83cf942a7af3))

## 1.2.1

The following sections list the changes for 1.2.1.

## Summary

 * Fix #138: Add missing value for consumption metrics

## Details

 * Bugfix #138: Add missing value for consumption metrics

   Until now the value for the consumption metrics was missing any assignment of a
   real value as the variable have only been defined with the default zero value.
   Beside that you are also able to add a currency label to make sure you can
   calculate correct currencies.

   https://github.com/promhippie/scw_exporter/issues/138


## 1.2.0

The following sections list the changes for 1.2.0.

## Summary

 * Chg #104: Switch to official logging library
 * Enh #114: Add metrics for consumption statistics

## Details

 * Change #104: Switch to official logging library

   Since there have been a structured logger part of the Go standard library we
   thought it's time to replace the library with that. Be aware that log messages
   should change a little bit.

   https://github.com/promhippie/scw_exporter/issues/104

 * Enhancement #114: Add metrics for consumption statistics

   We've added new metrics for the consumption API endpoints to give an overview
   about consumed and billed resources.

   https://github.com/promhippie/scw_exporter/pull/114


## 1.1.0

The following sections list the changes for 1.1.0.

## Summary

 * Chg #53: Read secrets form files
 * Chg #53: Integrate standard web config
 * Enh #53: Integrate option pprof profiling

## Details

 * Change #53: Read secrets form files

   We have added proper support to load secrets like the password from files or
   from base64-encoded strings. Just provide the flags or environment variables for
   token or private key with a DSN formatted string like `file://path/to/file` or
   `base64://Zm9vYmFy`.

   https://github.com/promhippie/scw_exporter/pull/53

 * Change #53: Integrate standard web config

   We integrated the new web config from the Prometheus toolkit which provides a
   configuration for TLS support and also some basic builtin authentication. For
   the detailed configuration you can check out the documentation.

   https://github.com/promhippie/scw_exporter/pull/53

 * Enhancement #53: Integrate option pprof profiling

   We have added an option to enable a pprof endpoint for proper profiling support
   with the help of tools like Parca. The endpoint `/debug/pprof` can now
   optionally be enabled to get the profiling details for catching potential memory
   leaks.

   https://github.com/promhippie/scw_exporter/pull/53


## 1.0.0

The following sections list the changes for 1.0.0.

## Summary

 * Chg #12: Refactor build tools and project structure
 * Chg #14: Drop darwin/386 release builds

## Details

 * Change #12: Refactor build tools and project structure

   To have a unified project structure and build tooling we have integrated the
   same structure we already got within our GitHub exporter.

   https://github.com/promhippie/scw_exporter/issues/12

 * Change #14: Drop darwin/386 release builds

   We dropped the build of 386 builds on Darwin as this architecture is not
   supported by current Go versions anymore.

   https://github.com/promhippie/scw_exporter/issues/14


## 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #11: Initial release of basic version

## Details

 * Change #11: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/scw_exporter/issues/11
