package archtest

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"testing"
)

const internalPrefix = "github.com/pangaunn/harness-engineering-workshop/internal/"

// TestArchitecture checks the real imports of every package under internal/
// (test files excluded) against the rules in rules.go.
func TestArchitecture(t *testing.T) {
	if len(layers) == 0 {
		t.Skip("no rules yet: add layers to internal/archtest/rules.go")
	}
	imports, err := internalImports("..")
	if err != nil {
		t.Fatal(err)
	}
	var pkgs []string
	for pkg := range imports {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	var failures []string
	for _, pkg := range pkgs {
		layer, ok := layers[pkg]
		if !ok {
			failures = append(failures, fmt.Sprintf("internal/%s is not registered\n"+
				"  Every package must be listed in `layers` in internal/archtest/rules.go.", pkg))
			continue
		}
		for _, imp := range imports[pkg] {
			impLayer, ok := layers[imp]
			if !ok || slices.Contains(allowed[layer], impLayer) {
				continue // unregistered imports are reported on their own
			}
			failures = append(failures, violation(pkg, layer, imp, impLayer))
		}
	}
	if len(failures) > 0 {
		t.Fatalf("%d architecture violation(s):\n\nFAIL: %s", len(failures), strings.Join(failures, "\n\nFAIL: "))
	}
}

func violation(pkg, layer, imp, impLayer string) string {
	may := "nothing"
	if len(allowed[layer]) > 0 {
		may = strings.Join(allowed[layer], ", ")
	}
	return fmt.Sprintf("internal/%s imports internal/%s\n"+
		"  Layer %q may not import layer %q. It may import: %s.", pkg, imp, layer, impLayer, may)
}

// internalImports maps each package directory under root (relative, slash
// separated) to the internal packages its non-test files import.
func internalImports(root string) (map[string][]string, error) {
	out := map[string][]string{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return err
		}
		rel, err := filepath.Rel(root, filepath.Dir(path))
		if err != nil {
			return err
		}
		pkg := filepath.ToSlash(rel)
		if pkg == "archtest" {
			return nil
		}
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}
		if _, ok := out[pkg]; !ok {
			out[pkg] = nil
		}
		for _, spec := range f.Imports {
			p, _ := strconv.Unquote(spec.Path.Value)
			if dep, ok := strings.CutPrefix(p, internalPrefix); ok && dep != pkg && !slices.Contains(out[pkg], dep) {
				out[pkg] = append(out[pkg], dep)
			}
		}
		return nil
	})
	for _, deps := range out {
		sort.Strings(deps)
	}
	return out, err
}
