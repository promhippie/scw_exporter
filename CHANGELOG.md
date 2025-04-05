# Changelog for unreleased

The following sections list the changes for unreleased.

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


# Changelog for 1.1.0

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


# Changelog for 1.0.0

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


# Changelog for 0.1.0

The following sections list the changes for 0.1.0.

## Summary

 * Chg #11: Initial release of basic version

## Details

 * Change #11: Initial release of basic version

   Just prepared an initial basic version which could be released to the public.

   https://github.com/promhippie/scw_exporter/issues/11


