package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/stjohnjohnson/bookend-scm-github/arguments"
	"github.com/stjohnjohnson/bookend-scm-github/git"
)

// VERSION gets set by the build script via the LDFLAGS
var VERSION string

var fmtPrint = fmt.Print
var osExit = os.Exit
var getArguments = arguments.GetArguments
var getGitVersion = git.GetGitVersion
var getGitSha = git.GetGitSha
var executeStream = git.ExecuteStream
var blackColor = color.New(color.FgHiBlack).SprintFunc()
var redColor = color.New(color.FgHiRed).SprintFunc()
var greenColor = color.New(color.FgHiGreen).SprintFunc()

func executeStreamFail(args ...string) {
	err := executeStream(args...)
	if err != nil {
		fmtPrint(redColor(fmt.Sprintf("%v\n", err)))
		osExit(1)
		return
	}
}

func main() {
	args, err := getArguments(os.Args)
	if args.Version {
		fmtPrint(VERSION)
		osExit(0)
		return
	}
	if err != nil {
		fmtPrint(redColor(fmt.Sprintf("CLI flags invalid: %v\n", err)))
		osExit(1)
		return
	}

	clientVersion, err := getGitVersion()
	if err != nil {
		fmtPrint(redColor(fmt.Sprintf("Unable to get Git version: %v\n", err)))
		osExit(1)
		return
	}

	fmtPrint(fmt.Sprintf("%s\tv%s\n", blackColor("Bookend:"), VERSION))
	fmtPrint(fmt.Sprintf("%s\t%s\n", blackColor("Git Client:"), clientVersion))

	fmtPrint(greenColor(fmt.Sprintf("\n☛ Cloning %s, on branch %s\n", args.ScmURL, args.Branch)))
	executeStreamFail("clone", "--quiet", "--progress", "--branch", args.Branch, args.CloneURL, args.TargetDir)
	os.Chdir(args.TargetDir)

	fmtPrint(greenColor("\n☛ Saving local git config\n"))
	executeStreamFail("config", "user.name", args.GitName)
	executeStreamFail("config", "user.email", args.GitEmail)

	if args.PullRequest != 0 {
		fmtPrint(greenColor(fmt.Sprintf("\n☛ Fetching PR %d\n", args.PullRequest)))
		executeStreamFail("fetch", "origin", fmt.Sprintf("pull/%d/head:pr", args.PullRequest))

		fmtPrint(greenColor(fmt.Sprintf("\n☛ Merging with %s\n", args.Branch)))
		executeStreamFail("merge", "--no-edit", args.SHA)

		gitSha, err := getGitSha()
		if err != nil {
			fmtPrint(redColor(fmt.Sprintf("Unable to get current Git revision: %v\n", err)))
			osExit(1)
			return
		}
		fmtPrint(greenColor(fmt.Sprintf("\n☛ Checked out %s", gitSha)))
	} else {
		fmtPrint(greenColor(fmt.Sprintf("\n☛ Resetting to %s\n", args.SHA)))
		executeStreamFail("reset", "--hard", args.SHA)
	}

	fmtPrint(greenColor("\n✓ Done\n"))
}
