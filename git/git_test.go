package git

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

type execFunc func(command string, args ...string) *exec.Cmd

func getFakeExecCommand(validator func(string, ...string)) execFunc {
	return func(command string, args ...string) *exec.Cmd {
		validator(command, args...)
		return fakeExecCommand(command, args...)
	}
}

func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestExecuteSuccess(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "/bin/git" }

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	err := ExecuteStream("foo", "bar")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestExecuteFailedCommand(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "/bin/git" }

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	err := ExecuteStream("baz", "bar")

	want := "Command failed: exit status 150"
	if err == nil || err.Error() != want {
		t.Errorf("Expected '%v', got '%v'", want, err)
	}
}

func TestGetGitVersion(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "/bin/git" }

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	version, err := GetGitVersion()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	want := "v1.2.3"
	if version != want {
		t.Errorf("Received the wrong version: %q, want %q", version, want)
	}
}

func TestGetGitVersionFail(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "/bin/fail" }

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	version, err := GetGitVersion()

	want := "Unable to get Git version: exit status 100"
	if err == nil || err.Error() != want {
		t.Errorf("Expected '%v', got '%v'", want, err)
	}

	want = ""
	if version != want {
		t.Errorf("Received the wrong version: %q, want %q", version, want)
	}
}

func TestGetGitSha(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "/bin/git" }

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	sha, err := GetGitSha()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	want := "302f5f5b48b9feee797a66c88811f1770bcb2dcf"
	if sha != want {
		t.Errorf("Received the wrong sha: %q, want %q", sha, want)
	}
}

func TestGetGitShaFail(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "/bin/false" }

	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	sha, err := GetGitSha()

	wantErr := "Unable to get current Git revision: exit status 100"
	if err == nil || err.Error() != wantErr {
		t.Errorf("Expected '%v', got '%v'", wantErr, err)
	}

	wantSha := ""
	if sha != wantSha {
		t.Errorf("Received the wrong sha: %v, wants %v", sha, wantSha)
	}
}

// This is a fake test for mocking out exec calls.
// See https://golang.org/src/os/exec/exec_test.go and
// https://npf.io/2015/06/testing-exec-command/ for more info
func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args[:]
	for i, val := range os.Args { // Should become something lke ["git", "tag"]
		args = os.Args[i:]
		if val == "--" {
			args = args[1:]
			break
		}
	}

	if len(args) >= 2 && args[0] == "/bin/git" {
		switch args[1] {
		default:
			os.Exit(150)
		case "foo":
			fmt.Println("OK")
			return
		case "rev-parse":
			fmt.Print("302f5f5b48b9feee797a66c88811f1770bcb2dcf")
			return
		case "--version":
			fmt.Print("git version 1.2.3")
			return
		}
	}

	fmt.Println(args[0])
	os.Exit(100)
}
