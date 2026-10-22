package tutorial

type Difficulty string

const (
	DifficultyBeginner     Difficulty = "beginner"
	DifficultyIntermediate Difficulty = "intermediate"
	DifficultyAdvanced     Difficulty = "advanced"
)

type Tutorial struct {
	ID                string
	Title             string
	Description       string
	Difficulty        Difficulty
	EstimatedDuration string
	Prerequisites     []string
	Steps             []TutorialStep
	Recommended       bool
}

type TutorialStep struct {
	Title          string
	Concept        string
	Command        string
	ExpectedOutput string
	ErrorHint      string
}

var gettingStarted = Tutorial{
	ID:                "getting-started",
	Title:             "Getting Started",
	Description:       "Learn the basics of the Scaleway CLI: configuration, commands, and output formats.",
	Difficulty:        DifficultyBeginner,
	EstimatedDuration: "~5 minutes",
	Recommended:       true,
	Prerequisites:     []string{},
	Steps: []TutorialStep{
		{
			Title:   "Welcome",
			Concept: "Welcome to the Scaleway CLI! This tutorial will guide you through the basics of using scw to manage your Scaleway cloud infrastructure.\n\nThe CLI follows a simple pattern: scw <namespace> <resource> <verb>\n\nFor example: scw instance server list",
		},
		{
			Title:          "Check Configuration",
			Concept:        "Let's verify your CLI is properly configured. This command shows your current settings including credentials, default zone, and project.",
			Command:        "scw config show",
			ExpectedOutput: "A table showing your CLI settings including access key, secret key (masked), default zone, region, project ID, and organization ID.",
			ErrorHint:      "If your config is empty, you need to run 'scw init' first to set up your credentials.",
		},
		{
			Title:          "Your First Command",
			Concept:        "Now let's list your servers. This is a read-only command that shows all compute instances in your default zone.",
			Command:        "scw instance server list",
			ExpectedOutput: "A table listing your servers with their ID, name, type, state, IP, and zone. If you have no servers, the table will be empty.",
			ErrorHint:      "If you get a permission error, make sure your API key has the necessary IAM permissions.",
		},
		{
			Title:          "Filtering Output",
			Concept:        "The CLI supports multiple output formats. You can use JSON output to pipe results to tools like jq for advanced filtering.",
			Command:        "scw instance server list -o json",
			ExpectedOutput: "JSON array of server objects. You can pipe this to jq: scw instance server list -o json | jq '.[0].name'",
			ErrorHint:      "Make sure you have servers in your default zone. Use --zone=all to list across all zones.",
		},
		{
			Title:   "Getting Help",
			Concept: "You can always get help for any command by adding --help. You can also browse the full documentation interactively by running 'scw docs'.\n\nTry it now: scw instance server list --help",
		},
		{
			Title:   "Congratulations!",
			Concept: "You've completed the Getting Started tutorial! You now know how to:\n\n- Check your CLI configuration\n- List resources using the namespace.resource.verb pattern\n- Use different output formats\n- Get help for any command\n\nTo learn more, run 'scw docs' to browse the full documentation, or 'scw tutorial' to see other available tutorials.",
		},
	},
}

var allTutorials = []Tutorial{
	gettingStarted,
}

func ListTutorials() []Tutorial {
	result := make([]Tutorial, len(allTutorials))
	copy(result, allTutorials)

	for i := range result {
		for j := i + 1; j < len(result); j++ {
			if result[j].Recommended && !result[i].Recommended {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

func GetTutorial(id string) (Tutorial, bool) {
	for _, t := range allTutorials {
		if t.ID == id {
			return t, true
		}
	}

	return Tutorial{}, false
}
