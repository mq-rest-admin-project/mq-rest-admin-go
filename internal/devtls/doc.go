// Package devtls provides TLS-trust helpers for the integration test suites.
//
// TLS certificate verification is always enabled. The development queue
// managers use self-signed certificates, so the integration suites trust them
// explicitly by extracting their certificate chains at runtime into a
// temporary CA bundle -- never by disabling verification. The implementation
// is build-tagged "integration"; this file keeps the package buildable in the
// default (unit) build.
package devtls
