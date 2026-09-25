package content

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
)

// BeadOutcome is what Headline extracts from one bead-N-report.md.
type BeadOutcome struct {
	Number    int
	Title     string
	Status    string
	Attempts  int
	Revisions int
	WallSecs  int
}

// Headline summarizes a run's per-bead reports into "did anything
// interesting happen" — the magazine-headline view above the full reports,
// which stay one click away for anyone who wants the whole story.
type Headline struct {
	TotalBeads       int
	CleanBeads       int // succeeded, one attempt, one spec revision
	Escalated        []BeadOutcome
	MultiAttempt     []BeadOutcome
	Revised          []BeadOutcome
	TotalWallSeconds int
}

func (h Headline) WallTimeStr() string {
	m := h.TotalWallSeconds / 60
	if m < 60 {
		return fmt.Sprintf("%dm", m)
	}
	return fmt.Sprintf("%dh %dm", m/60, m%60)
}

var (
	beadTitleRe = regexp.MustCompile(`(?m)^# Bead (\d+): (.+)$`)
	statusRe    = regexp.MustCompile(`\*\*Status:\*\*\s*(\S+)`)
	attemptsRe  = regexp.MustCompile(`\*\*Attempts:\*\*\s*(\d+)`)
	wallRe      = regexp.MustCompile(`\*\*Wall time:\*\*\s*(\d+)s`)
	revisionRe  = regexp.MustCompile(`(?m)^### Revision (\d+)`)
)

// Headline reads every bead report for slug and summarizes them. Reports
// that don't match the expected shape are skipped rather than failing the
// whole summary — this is best-effort scanning of Ratchet's own output, not
// a contract either side is obligated to keep stable.
func (s *Store) Headline(slug string) (Headline, error) {
	names, err := s.Traces(slug)
	if err != nil {
		return Headline{}, err
	}
	var h Headline
	for _, name := range names {
		data, err := os.ReadFile(s.TracePath(slug, name))
		if err != nil {
			continue
		}
		src := string(data)

		outcome := BeadOutcome{Attempts: 1, Revisions: 1}
		if m := beadTitleRe.FindStringSubmatch(src); m != nil {
			outcome.Number, _ = strconv.Atoi(m[1])
			outcome.Title = m[2]
		} else {
			continue // not a bead report we recognize
		}
		if m := statusRe.FindStringSubmatch(src); m != nil {
			outcome.Status = m[1]
		}
		if m := attemptsRe.FindStringSubmatch(src); m != nil {
			outcome.Attempts, _ = strconv.Atoi(m[1])
		}
		if m := wallRe.FindStringSubmatch(src); m != nil {
			secs, _ := strconv.Atoi(m[1])
			outcome.WallSecs = secs
			h.TotalWallSeconds += secs
		}
		for _, m := range revisionRe.FindAllStringSubmatch(src, -1) {
			n, _ := strconv.Atoi(m[1])
			if n > outcome.Revisions {
				outcome.Revisions = n
			}
		}

		h.TotalBeads++
		switch {
		case outcome.Status != "" && outcome.Status != "succeeded":
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
