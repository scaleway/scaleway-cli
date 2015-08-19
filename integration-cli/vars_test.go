package integrationcli

var (
	scwcli         string   = "../scw"
	publicCommands []string = []string{
		"help", "attach", "commit", "cp", "create",
		"events", "exec", "history", "images", "info",
		"inspect", "kill", "login", "logout", "logs",
		"port", "ps", "rename", "restart", "rm", "rmi",
		"run", "search", "start", "stop", "tag", "top",
		"version", "wait",
	}
	secretCommands []string = []string{
		"_patch", "_completion", "_flush-cache",
	}
	publicOptions []string = []string{
		"-h, --help=false",
		"-D, --debug=false",
		"-V, --verbose=false",
		"-q, --quiet=false",
		"--api-endpoint=APIEndPoint",
		"--sensitive=false",
		"-v, --version=false",
	}
)
