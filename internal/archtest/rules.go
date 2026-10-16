// Package archtest enforces the import rules in docs/architecture.md.
package archtest

// layers assigns every package under internal/ to a layer.
var layers = map[string]string{
	"handler":              "handler",
	"service/booking":      "aggregate",
	"service/room":         "domain",
	"service/availability": "domain",
	"service/pricing":      "domain",
	"service/payment":      "domain",
	"service/report":       "domain",
	"repository":           "repository",
	"model":                "model",
	"jobs":                 "jobs",
}

// allowed lists which layers each layer may import.
var allowed = map[string][]string{
	"handler":    {"aggregate", "domain", "model"},
	"jobs":       {"aggregate", "domain", "model"},
	"aggregate":  {"domain", "repository", "model"},
	"domain":     {"repository", "model"},
	"repository": {"model"},
	"model":      {},
}
