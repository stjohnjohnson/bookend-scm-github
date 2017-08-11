package git

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

var osGetEnv = os.Getenv
var execCommand = exec.Command

// GetGitVersion returns the version of Git that we're using
func GetGitVersion() (string, error) {
	out, err := ExecuteReturn("--version")
	if err != nil {
		return "", fmt.Errorf("Unable to get Git version: %v", err)
	}
	re := regexp.MustCompile("git version (.*)")
	match := re.FindStringSubmatch(out)

	return fmt.Sprintf("v%s", match[1]), nil
}

// GetGitSha returns the current SHA
func GetGitSha() (string, error) {
	out, err := ExecuteReturn("rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("Unable to get current Git revision: %v", err)
	}
	return out, nil
}

// ExecuteStream will stream the input/output from a Git call
func ExecuteStream(arguments ...string) error {
	cmd := execCommand(osGetEnv("GIT_PATH"), arguments...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Command wouldn't start: %v", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Command failed: %v", err)
	}

	return nil
}

// ExecuteReturn will return the output from a Git call
func ExecuteReturn(arguments ...string) (string, error) {
	bytes, err := execCommand(osGetEnv("GIT_PATH"), arguments...).Output()

	return string(bytes), err
}
