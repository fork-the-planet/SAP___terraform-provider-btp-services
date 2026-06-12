// internal/version/version.go

package version

// ProviderVersion is replaced with the actual release tag during the build process.
// See the Makefile ld-flags target: -ldflags "-X github.com/.../version.ProviderVersion=<tag>"
var ProviderVersion string = "dev"
