package archtest

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// docs are the files that must name every service registered in layers.
var docs = []string{"AGENTS.md", "docs/architecture.md"}

// TestDocsNameEveryService fails when a service is registered in layers but
// the docs that explain the rules don't mention it, so the docs grow with the
// code.
func TestDocsNameEveryService(t *testing.T) {
	var failures []string
	for _, doc := range docs {
		text, err := os.ReadFile(filepath.Join("..", "..", doc))
		if err != nil {
			t.Fatal(err)
		}
		var missing []string
		for pkg, layer := range layers {
			name, ok := strings.CutPrefix(pkg, "service/")
			if ok && !strings.Contains(string(text), "`"+name+"`") {
				missing = append(missing, fmt.Sprintf("`%s` (%s)", name, layer))
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			failures = append(failures, fmt.Sprintf("%s does not name %s\n"+
				"  Every service in internal/archtest/rules.go must be named where the rules are explained:\n"+
				"  the domain/aggregate lists in %s.\n"+
				"  See \"Adding a package\" in AGENTS.md", doc, strings.Join(missing, ", "), doc))
		}
	}
	if len(failures) > 0 {
		t.Fatalf("docs are behind the code:\n\nFAIL: %s", strings.Join(failures, "\n\nFAIL: "))
	}
}
