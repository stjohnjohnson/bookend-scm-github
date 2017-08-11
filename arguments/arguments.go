package arguments

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// CommandArgs is the complete list of arguments we would get back
type CommandArgs struct {
	ScmURL        string
	Host          string
	Repo          string
	CloneURL      string
	Branch        string
	SHA           string
	PullRequest   int
	CloneMethod   string
	TargetDir     string
	GitName       string
	GitEmail      string
	HTTPSUsername string
	HTTPSToken    string
	Version       bool
}

var osGetEnv = os.Getenv

func getFlags(args []string) CommandArgs {
	var config CommandArgs

	f := flag.NewFlagSet(args[0], flag.ExitOnError)

	f.StringVar(&config.Host, "host", "", "Repository Host")
	f.StringVar(&config.Repo, "repo", "", "Repository Org/Repo")
	f.StringVar(&config.Branch, "branch", "master", "Checkout branch")
	f.StringVar(&config.SHA, "sha", "", "Commit SHA1")

	f.IntVar(&config.PullRequest, "pull-request", 0, "Pull Request Number")
	f.StringVar(&config.TargetDir, "target-dir", "", "Checkout directory")
	f.StringVar(&config.CloneMethod, "clone-method", "https", "Git Clone Method (https|ssh)")

	f.StringVar(&config.GitName, "git-name", "sd-buildbot", "Name in Git Config")
	f.StringVar(&config.GitEmail, "git-email", "dev-null@screwdriver.cd", "Email in Git Config")

	f.StringVar(&config.HTTPSUsername, "https-username", osGetEnv("SCM_USERNAME"), "Username to use when authenticating via HTTPS")
	f.StringVar(&config.HTTPSToken, "https-token", osGetEnv("SCM_ACCESS_TOKEN"), "Token to use when authenticating via HTTPS")

	f.BoolVar(&config.Version, "version", false, "Display Version number")

	f.Parse(args[1:])

	config.ScmURL = fmt.Sprintf("%s/%s", config.Host, config.Repo)

	return config
}

func validateConfig(config CommandArgs) error {
	if config.Host == "" {
		return errors.New("--scm-host is required")
	}
	if config.Repo == "" {
		return errors.New("--scm-repo is required")
	}
	if config.SHA == "" {
		return errors.New("--sha is required")
	}
	if config.TargetDir == "" {
		return errors.New("--target-dir is required")
	}
	return nil
}

func addDynamicConfig(config CommandArgs) (CommandArgs, error) {
	switch {
	case config.CloneMethod == "https" && config.HTTPSUsername != "" && config.HTTPSToken != "":
		config.CloneURL = fmt.Sprintf("https://%s:%s@%s/%s.git", config.HTTPSUsername, config.HTTPSToken, config.Host, config.Repo)
	case config.CloneMethod == "https":
		config.CloneURL = fmt.Sprintf("https://%s/%s.git", config.Host, config.Repo)
	case config.CloneMethod == "ssh":
		config.CloneURL = fmt.Sprintf("git@%s:%s.git", config.Host, config.Repo)
	default:
		return config, errors.New("--clone-method must be https or ssh")
	}

	return config, nil
}

// GetArguments returns the flags and options set on the command-line
func GetArguments(args []string) (CommandArgs, error) {
	config := getFlags(args)
	err := validateConfig(config)

	if err != nil {
		return config, err
	}

	return addDynamicConfig(config)
}
