// Package archtest checks the imports between packages under internal/
// against the rules below. Run it with: go test ./internal/archtest/
package archtest

// layers puts every package under internal/ into a layer.
// The key is the package path below internal/, the value is a layer name you choose.
// While layers is empty the check is skipped. Once it has an entry, every
// package must be listed.
var layers = map[string]string{
	// Example from a blog app (not this repo):
	//   "web":          "web",
	//   "service/post": "service",
	//   "service/user": "service",
	//   "store":        "store",
}

// allowed lists, for each layer, the layers it may import.
var allowed = map[string][]string{
	// Example from the same blog app: web may use services, services may use
	// the store, and the store imports nothing from internal/.
	//   "web":     {"service"},
	//   "service": {"store"},
	//   "store":   {},
}
