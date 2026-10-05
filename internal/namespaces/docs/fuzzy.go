package docs

import (
	"strings"
	"unicode"
)

// fuzzyMatch performs subsequence-based fuzzy matching (similar to fzf).
// It checks whether all characters of the query appear in the target string
// in order (case-insensitive). Returns a match score (higher = better) and
// true if matched, or 0 and false if not matched.
//
// Scoring heuristics:
//   - Consecutive matches get a bonus (adjacency).
//   - Matches at word boundaries (start of string, after space, after camelCase
//     transitions) get a bonus.
//   - Compact matches (smaller gap between first and last matched char) score higher.
func FuzzyMatch(target, query string) (int, bool) {
	if query == "" {
		return 0, true
	}

	targetLower := strings.ToLower(target)
	queryLower := strings.ToLower(query)

	queryRunes := []rune(queryLower)
	targetRunes := []rune(targetLower)

	qi := 0
	score := 0
	consecutive := 0
	firstMatch := -1
	lastMatch := -1

	for ti := 0; ti < len(targetRunes) && qi < len(queryRunes); ti++ {
		if targetRunes[ti] == queryRunes[qi] {
			if firstMatch == -1 {
				firstMatch = ti
			}
			lastMatch = ti

			if consecutive > 0 {
				score += 10 + consecutive*2
			} else {
				score += 1
			}
			consecutive++

			if isWordBoundary(targetRunes, ti) {
				score += 15
			}

			qi++
		} else {
			consecutive = 0
		}
	}

	if qi < len(queryRunes) {
		return 0, false
	}

	if firstMatch == 0 {
		score += 20
	}

	if lastMatch >= 0 && firstMatch >= 0 {
		span := lastMatch - firstMatch
		if span > 0 {
			score -= span
		}
	}

	return score, true
}

func isWordBoundary(runes []rune, i int) bool {
	if i == 0 {
		return true
	}

	prev := runes[i-1]
	curr := runes[i]

	if prev == ' ' || prev == '-' || prev == '.' || prev == '/' || prev == '_' {
		return true
	}

	if unicode.IsLower(prev) && unicode.IsUpper(curr) {
		return true
	}

	return false
}

// fuzzySearch returns matched items sorted by fuzzy score (best first).
type fuzzyItem struct {
	value string
	score int
}

func fuzzyFilter(items []string, query string) []string {
	if query == "" {
		return items
	}

	var matched []fuzzyItem
	for _, item := range items {
		score, ok := FuzzyMatch(item, query)
		if ok {
			matched = append(matched, fuzzyItem{value: item, score: score})
		}
	}

	// Sort by score descending, then alphabetically for ties.
	for i := 0; i < len(matched); i++ {
		for j := i + 1; j < len(matched); j++ {
			if matched[j].score > matched[i].score ||
				(matched[j].score == matched[i].score && matched[j].value < matched[i].value) {
				matched[i], matched[j] = matched[j], matched[i]
			}
		}
	}

	result := make([]string, len(matched))
	for i, m := range matched {
		result[i] = m.value
	}

	return result
}
