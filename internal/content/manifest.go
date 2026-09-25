// Package content loads the curated per-project metadata (title, description,
// backend address, design-doc/survey filenames) that ratchet-gallery renders.
// Every read hits the filesystem fresh — onboarding a project via
// scripts/onboard.sh needs no gallery restart, just a browser refresh.
package content

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

// Manifest describes one onboarded Ratchet project.
type Manifest struct {
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Added       string `json:"added"`
	// Backend is the container's host:port on the shared "web" Docker
	// network, e.g. "fractalviz:8080".
	Backend   string `json:"backend"`
	DesignDoc string `json:"design_doc"`
	Survey    string `json:"survey"`
}

type Store struct {
	dir string
}

func NewStore(dir string) *Store {
	return &Store{dir: dir}
}

// List returns every onboarded project's manifest, sorted by slug.
func (s *Store) List() ([]Manifest, error) {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out []Manifest
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		m, err := s.Load(e.Name())
		if err != nil {
			continue // skip malformed/incomplete entries rather than fail the whole list
		}
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Slug < out[j].Slug })
	return out, nil
}

// Load reads content/<slug>/manifest.json.
func (s *Store) Load(slug string) (Manifest, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, slug, "manifest.json"))
	if err != nil {
		return Manifest{}, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return Manifest{}, err
	}
	if m.Slug == "" {
		m.Slug = slug
	}
	return m, nil
}

// DocPath resolves a filename from a manifest (e.g. DesignDoc) to its path
// on disk under content/<slug>/.
func (s *Store) DocPath(slug, filename string) string {
	return filepath.Join(s.dir, slug, filename)
}

var beadNum = regexp.MustCompile(`bead-(\d+)`)

// Traces lists the per-bead build reports under content/<slug>/traces/
// (Ratchet's own mechanically-written, no-model-call post-execution
// reports), ordered by bead number. The full project-report.md is
// deliberately not onboarded here — it can run to tens of megabytes,
// unlike the per-bead reports.
func (s *Store) Traces(slug string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(s.dir, slug, "traces"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Slice(names, func(i, j int) bool {
		ni, nj := beadNumOf(names[i]), beadNumOf(names[j])
		if ni != nj {
			return ni < nj
		}
		return names[i] < names[j]
	})
	return names, nil
}

func beadNumOf(name string) int {
	m := beadNum.FindStringSubmatch(name)
	if m == nil {
		return -1
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

// TracePath resolves a trace filename to its path on disk, guarding against
// path traversal via the filename coming from a URL path segment.
func (s *Store) TracePath(slug, filename string) string {
	return filepath.Join(s.dir, slug, "traces", filepath.Base(filename))
}
