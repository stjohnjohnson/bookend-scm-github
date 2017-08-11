package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stjohnjohnson/bookend-scm-github/arguments"
)

func mockPrint(messages []string, t *testing.T) func(...interface{}) (int, error) {
	index := 0

	return func(input ...interface{}) (int, error) {
		message := fmt.Sprint(input...)
		if index >= len(messages) {
			t.Errorf("Received an unexpected message: %v", message)
		} else if messages[index] != message {
			t.Errorf("Received the wrong message: %v, want %v", message, messages[index])
		}
		index++

		return 0, nil
	}
}

func mockExec(commands []string, t *testing.T, fail bool) func(...string) error {
	index := 0

	return func(input ...string) error {
		command := strings.Join(input, " ")
		if index >= len(commands) {
			t.Errorf("Received an unexpected command: %v", command)
		} else if commands[index] != command {
			t.Errorf("Received the wrong command: %v, want %v", command, commands[index])
		}
		index++

		if fail {
			return errors.New("Failed command")
		}
		return nil
	}
}

func mockExit(want int, t *testing.T) func(int) {
	return func(code int) {
		if code != want {
			t.Errorf("Received the wrong exit code: %v, want %v", code, want)
		}
	}
}

func TestMain(m *testing.M) {
	VERSION = "1.0.0"
	os.Exit(m.Run())
}

func TestMainNonPR(t *testing.T) {
	executeStream = mockExec([]string{
		"clone --quiet --progress --branch master https://github.com/testOrg/testRepo.git /tmp/foo",
		"config user.name sd-buildbot",
		"config user.email dev-null@screwdriver.cd",
		"reset --hard 302f5f5b48b9feee797a66c88811f1770bcb2dcf",
	}, t, false)
	fmtPrint = mockPrint([]string{
		"Bookend:\tv1.0.0\n",
		"Git Client:\tv1.2.3\n",
		"\n☛ Cloning github.com/testOrg/testRepo, on branch master\n",
		"\n☛ Saving local git config\n",
		"\n☛ Resetting to 302f5f5b48b9feee797a66c88811f1770bcb2dcf\n",
		"\n✓ Done\n",
	}, t)
	osExit = mockExit(0, t)
	getGitVersion = func() (string, error) { return "v1.2.3", nil }
	getGitSha = func() (string, error) { return "ace893fb2c9553a38a873fb03d0e21a406b351a1", nil }

	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{
			Host:        "github.com",
			Repo:        "testOrg/testRepo",
			Branch:      "master",
			ScmURL:      "github.com/testOrg/testRepo",
			SHA:         "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
			PullRequest: 0,
			TargetDir:   "/tmp/foo",
			CloneMethod: "https",
			CloneURL:    "https://github.com/testOrg/testRepo.git",
			GitName:     "sd-buildbot",
			GitEmail:    "dev-null@screwdriver.cd",
			Version:     false,
		}, nil
	}
	main()
}

func TestMainPR(t *testing.T) {
	executeStream = mockExec([]string{
		"clone --quiet --progress --branch master https://github.com/testOrg/testRepo.git /tmp/foo",
		"config user.name sd-buildbot",
		"config user.email dev-null@screwdriver.cd",
		"fetch origin pull/15/head:pr",
		"merge --no-edit ace893fb2c9553a38a873fb03d0e21a406b351a1",
	}, t, false)
	fmtPrint = mockPrint([]string{
		"Bookend:\tv1.0.0\n",
		"Git Client:\tv1.2.3\n",
		"\n☛ Cloning github.com/testOrg/testRepo, on branch master\n",
		"\n☛ Saving local git config\n",
		"\n☛ Fetching PR 15\n",
		"\n☛ Merging with master\n",
		"\n☛ Checked out 302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"\n✓ Done\n",
	}, t)
	osExit = mockExit(0, t)
	getGitVersion = func() (string, error) { return "v1.2.3", nil }
	getGitSha = func() (string, error) { return "302f5f5b48b9feee797a66c88811f1770bcb2dcf", nil }

	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{
			Host:        "github.com",
			Repo:        "testOrg/testRepo",
			Branch:      "master",
			ScmURL:      "github.com/testOrg/testRepo",
			SHA:         "ace893fb2c9553a38a873fb03d0e21a406b351a1",
			PullRequest: 15,
			TargetDir:   "/tmp/foo",
			CloneMethod: "https",
			CloneURL:    "https://github.com/testOrg/testRepo.git",
			GitName:     "sd-buildbot",
			GitEmail:    "dev-null@screwdriver.cd",
			Version:     false,
		}, nil
	}
	main()
}

func TestMainFailCommand(t *testing.T) {
	executeStream = mockExec([]string{
		"clone --quiet --progress --branch master https://github.com/testOrg/testRepo.git /tmp/foo",
		"config user.name sd-buildbot",
		"config user.email dev-null@screwdriver.cd",
		"reset --hard 302f5f5b48b9feee797a66c88811f1770bcb2dcf",
	}, t, true)
	fmtPrint = mockPrint([]string{
		"Bookend:\tv1.0.0\n",
		"Git Client:\tv1.2.3\n",
		"\n☛ Cloning github.com/testOrg/testRepo, on branch master\n",
		"Failed command\n",
		"\n☛ Saving local git config\n",
		"Failed command\n",
		"Failed command\n",
		"\n☛ Resetting to 302f5f5b48b9feee797a66c88811f1770bcb2dcf\n",
		"Failed command\n",
		"\n✓ Done\n",
	}, t)
	osExit = mockExit(1, t)
	getGitVersion = func() (string, error) { return "v1.2.3", nil }
	getGitSha = func() (string, error) { return "ace893fb2c9553a38a873fb03d0e21a406b351a1", nil }

	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{
			Host:        "github.com",
			Repo:        "testOrg/testRepo",
			Branch:      "master",
			ScmURL:      "github.com/testOrg/testRepo",
			SHA:         "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
			PullRequest: 0,
			TargetDir:   "/tmp/foo",
			CloneMethod: "https",
			CloneURL:    "https://github.com/testOrg/testRepo.git",
			GitName:     "sd-buildbot",
			GitEmail:    "dev-null@screwdriver.cd",
			Version:     false,
		}, nil
	}
	main()
}

func TestMainNoVersion(t *testing.T) {
	executeStream = mockExec([]string{}, t, false)
	fmtPrint = mockPrint([]string{
		"Unable to get Git version: Bad Version\n",
	}, t)
	osExit = mockExit(1, t)
	getGitVersion = func() (string, error) { return "", errors.New("Bad Version") }
	getGitSha = func() (string, error) { return "302f5f5b48b9feee797a66c88811f1770bcb2dcf", nil }

	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{
			Host:        "github.com",
			Repo:        "testOrg/testRepo",
			Branch:      "master",
			ScmURL:      "github.com/testOrg/testRepo",
			SHA:         "ace893fb2c9553a38a873fb03d0e21a406b351a1",
			PullRequest: 15,
			TargetDir:   "/tmp/foo",
			CloneMethod: "https",
			CloneURL:    "https://github.com/testOrg/testRepo.git",
			GitName:     "sd-buildbot",
			GitEmail:    "dev-null@screwdriver.cd",
			Version:     false,
		}, nil
	}
	main()
}

func TestMainBadArgs(t *testing.T) {
	executeStream = mockExec([]string{}, t, false)
	fmtPrint = mockPrint([]string{
		"CLI flags invalid: --foo is required\n",
	}, t)
	osExit = mockExit(1, t)
	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{}, errors.New("--foo is required")
	}
	main()
}

func TestMainNoSha(t *testing.T) {
	executeStream = mockExec([]string{
		"clone --quiet --progress --branch master https://github.com/testOrg/testRepo.git /tmp/foo",
		"config user.name sd-buildbot",
		"config user.email dev-null@screwdriver.cd",
		"fetch origin pull/15/head:pr",
		"merge --no-edit ace893fb2c9553a38a873fb03d0e21a406b351a1",
	}, t, false)
	fmtPrint = mockPrint([]string{
		"Bookend:\tv1.0.0\n",
		"Git Client:\tv1.2.3\n",
		"\n☛ Cloning github.com/testOrg/testRepo, on branch master\n",
		"\n☛ Saving local git config\n",
		"\n☛ Fetching PR 15\n",
		"\n☛ Merging with master\n",
		"Unable to get current Git revision: Bad Revision\n",
	}, t)
	osExit = mockExit(1, t)
	getGitVersion = func() (string, error) { return "v1.2.3", nil }
	getGitSha = func() (string, error) { return "", errors.New("Bad Revision") }

	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{
			Host:        "github.com",
			Repo:        "testOrg/testRepo",
			Branch:      "master",
			ScmURL:      "github.com/testOrg/testRepo",
			SHA:         "ace893fb2c9553a38a873fb03d0e21a406b351a1",
			PullRequest: 15,
			TargetDir:   "/tmp/foo",
			CloneMethod: "https",
			CloneURL:    "https://github.com/testOrg/testRepo.git",
			GitName:     "sd-buildbot",
			GitEmail:    "dev-null@screwdriver.cd",
			Version:     false,
		}, nil
	}
	main()
}

func TestMainVersion(t *testing.T) {
	executeStream = mockExec([]string{}, t, false)
	fmtPrint = mockPrint([]string{
		"1.0.0",
	}, t)
	osExit = mockExit(0, t)

	getArguments = func(_ []string) (arguments.CommandArgs, error) {
		return arguments.CommandArgs{
			Version: true,
		}, nil
	}
	main()
}
