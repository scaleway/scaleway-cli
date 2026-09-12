package docs

import (
	"sort"

	"github.com/scaleway/scaleway-cli/v2/core"
)

type SearchResult struct {
	Path  string
	Short string
}

func SearchCommands(commands *core.Commands, keyword string) []SearchResult {
	type scoredResult struct {
		result SearchResult
		score  int
	}

	var scored []scoredResult

	for _, cmd := range commands.GetAll() {
		if cmd.Hidden {
			continue
		}
		if cmd.Run == nil {
			continue
		}

		path := cmd.GetCommandLine("scw")
		short := cmd.Short
		long := cmd.Long

		pathScore, pathMatch := FuzzyMatch(path, keyword)
		shortScore, shortMatch := FuzzyMatch(short, keyword)
		longScore, longMatch := FuzzyMatch(long, keyword)

		if !pathMatch && !shortMatch && !longMatch {
			continue
		}

		score := 0
		if pathMatch {
			score += pathScore
		}
		if shortMatch {
			score += shortScore
		}
		if longMatch {
			score += longScore / 2
		}

		scored = append(scored, scoredResult{
			result: SearchResult{Path: path, Short: short},
			score:  score,
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		if scored[i].score == scored[j].score {
			return scored[i].result.Path < scored[j].result.Path
		}

		return scored[i].score > scored[j].score
	})

	results := make([]SearchResult, len(scored))
	for i, s := range scored {
		results[i] = s.result
	}

	return results
}
