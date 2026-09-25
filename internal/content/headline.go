package content

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// BeadOutcome is one row of project-summary.md's Bead Summary table.
type BeadOutcome struct {
	Number    int
	Title     string
	Status    string
	Attempts  int
	Revisions int
	WallSecs  int
}

// Headline summarizes a run into "did anything interesting happen" — the
// magazine-headline view above the full reports, which stay one click away
// for anyone who wants the whole story.
type Headline struct {
	TotalBeads       int
	CleanBeads       int // succeeded, one attempt, one spec revision
	Escalated        []BeadOutcome
	MultiAttempt     []BeadOutcome
	Revised          []BeadOutcome
	TotalWallSeconds int // real created_at -> updated_at elapsed time, not summed per-bead attempt time
}

func (h Headline) WallTimeStr() string {
	m := h.TotalWallSeconds / 60
	if m < 60 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh %dm", m/60, m%60)
}

var (
	projectWallTimeRe = regexp.MustCompile(`\*\*Wall time:\*\*\s*(\d+)s`)
	beadRowRe         = regexp.MustCompile(`(?m)^\|\s*(\d+)\s*\|\s*(.+?)\s*\|\s*(\S+)\s*\|\s*(\d+)\s*\|\s*(\d+)\s*\|\s*(\d+)s\s*\|$`)
)

// Headline reads content/<slug>/project-summary.md — the small header
// onboard.sh extracts from Ratchet's own project-report.md before its
// multi-megabyte final-source-files dump — and summarizes it. A missing or
// unrecognized file yields a zero Headline (hidden by the template) rather
// than an error; this is best-effort scanning of Ratchet's own report
// format, not a contract either side is obligated to keep stable.
func (s *Store) Headline(slug string) (Headline, error) {
	data, err := os.ReadFile(filepath.Join(s.dir, slug, "project-summary.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return Headline{}, nil
		}
		return Headline{}, err
	}
	src := string(data)

	var h Headline
	if m := projectWallTimeRe.FindStringSubmatch(src); m != nil {
		h.TotalWallSeconds, _ = strconv.Atoi(m[1])
	}

	for _, m := range beadRowRe.FindAllStringSubmatch(src, -1) {
		outcome := BeadOutcome{}
		outcome.Number, _ = strconv.Atoi(m[1])
		outcome.Title = m[2]
		outcome.Status = m[3]
		outcome.Attempts, _ = strconv.Atoi(m[4])
		outcome.Revisions, _ = strconv.Atoi(m[5])
		outcome.WallSecs, _ = strconv.Atoi(m[6])

		h.TotalBeads++
		switch {
		case outcome.Status != "succeeded":
			h.Escalated = append(h.Escalated, outcome)
		case outcome.Attempts > 1:
			h.MultiAttempt = append(h.MultiAttempt, outcome)
		case outcome.Revisions > 1:
			h.Revised = append(h.Revised, outcome)
		default:
			h.CleanBeads++
		}
	}
	return h, nil
}
