package diagrams

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"system-design/internal/registry"
)

var diagramMacroRE = regexp.MustCompile(`\{\{diagram\s+"([^"]+)"`)

func TestProblemAndFundamentalContentFilesExist(t *testing.T) {
	root := repoRoot(t)
	reg, err := registry.Load(os.DirFS(root), "content/_registry.yaml")
	if err != nil {
		t.Fatalf("load registry: %v", err)
	}

	for _, p := range reg.Problems {
		p := filepath.Join(root, "content", p.Path, "index.html")
		if _, err := os.Stat(p); err != nil {
			t.Errorf("missing problem content file: %s", p)
		}
	}

	for _, f := range reg.AllFundamentals() {
		p := filepath.Join(root, "content", f.Path, "index.html")
		if _, err := os.Stat(p); err != nil {
			t.Errorf("missing fundamental content file: %s", p)
		}
	}
}

func TestProblemAndFundamentalDiagramReferencesAreValid(t *testing.T) {
	reg := BuildDefault()
	contentToSlugs := map[string]map[string]bool{}

	root := repoRoot(t)
	for _, dir := range []string{filepath.Join(root, "content", "problems"), filepath.Join(root, "content", "fundamentals")} {
		err := walkIndexFiles(dir, func(filePath string) {
			contentRoot := filepath.Join(root, "content") + string(os.PathSeparator)
			rel := strings.TrimSuffix(strings.TrimPrefix(filePath, contentRoot), string(os.PathSeparator)+"index.html")
			rel = filepath.ToSlash(rel)
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Errorf("read %s: %v", filePath, err)
				return
			}

			matches := diagramMacroRE.FindAllStringSubmatch(string(data), -1)
			if _, ok := contentToSlugs[rel]; !ok {
				contentToSlugs[rel] = map[string]bool{}
			}
			for _, m := range matches {
				slug := m[1]
				contentToSlugs[rel][slug] = true
				d := reg.Get(slug)
				if d == nil {
					t.Errorf("unknown diagram slug %q in %s", slug, filePath)
					continue
				}
				if d.ContentFile != rel {
					t.Errorf("diagram %q used in %s but registered for %s", slug, rel, d.ContentFile)
				}
			}
		})
		if err != nil {
			t.Fatalf("walk %s: %v", dir, err)
		}
	}

}

func walkIndexFiles(root string, onFile func(filePath string)) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	for _, e := range entries {
		p := filepath.Join(root, e.Name())
		if e.IsDir() {
			if err := walkIndexFiles(p, onFile); err != nil {
				return err
			}
			continue
		}
		if e.Name() == "index.html" {
			onFile(p)
		}
	}
	return nil
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("unable to resolve caller path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}
