package arguments

import (
	"reflect"
	"testing"
)

func TestGetArguments(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "" }
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
	}

	args, err := GetArguments(osArgs)
	want := CommandArgs{
		Host:        "github.com",
		Repo:        "testOrg/testRepo",
		Branch:      "master",
		ScmURL:      "github.com/testOrg/testRepo",
		CloneURL:    "https://github.com/testOrg/testRepo.git",
		SHA:         "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		PullRequest: 0,
		TargetDir:   "/tmp/foo",
		CloneMethod: "https",
		GitName:     "sd-buildbot",
		GitEmail:    "dev-null@screwdriver.cd",
		Version:     false,
	}

	if !reflect.DeepEqual(args, want) {
		t.Errorf("Received the wrong arguments: %v, want %v", args, want)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetArgumentsEnvironment(t *testing.T) {
	osGetEnv = func(variable string) (out string) {
		if variable == "SCM_USERNAME" {
			return "stjohn"
		}
		if variable == "SCM_ACCESS_TOKEN" {
			return "875fc3f0c3613de2a999295616af7db0fced4056"
		}
		return ""
	}
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
	}

	args, err := GetArguments(osArgs)
	want := CommandArgs{
		Host:          "github.com",
		Repo:          "testOrg/testRepo",
		Branch:        "master",
		ScmURL:        "github.com/testOrg/testRepo",
		CloneURL:      "https://stjohn:875fc3f0c3613de2a999295616af7db0fced4056@github.com/testOrg/testRepo.git",
		SHA:           "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		PullRequest:   0,
		TargetDir:     "/tmp/foo",
		CloneMethod:   "https",
		GitName:       "sd-buildbot",
		GitEmail:      "dev-null@screwdriver.cd",
		Version:       false,
		HTTPSUsername: "stjohn",
		HTTPSToken:    "875fc3f0c3613de2a999295616af7db0fced4056",
	}

	if !reflect.DeepEqual(args, want) {
		t.Errorf("Received the wrong arguments: %v, want %v", args, want)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDynamicArgumentsHttps(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "" }
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
		"--clone-method=https",
	}

	args, err := GetArguments(osArgs)
	want := CommandArgs{
		Host:        "github.com",
		Repo:        "testOrg/testRepo",
		Branch:      "master",
		ScmURL:      "github.com/testOrg/testRepo",
		CloneURL:    "https://github.com/testOrg/testRepo.git",
		SHA:         "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		PullRequest: 0,
		TargetDir:   "/tmp/foo",
		CloneMethod: "https",
		GitName:     "sd-buildbot",
		GitEmail:    "dev-null@screwdriver.cd",
		Version:     false,
	}

	if !reflect.DeepEqual(args, want) {
		t.Errorf("Received the wrong arguments: %v, want %v", args, want)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDynamicArgumentsHttpsAuth(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "" }
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
		"--clone-method=https",
		"--https-username=stjohn",
		"--https-token=875fc3f0c3613de2a999295616af7db0fced4056",
	}

	args, err := GetArguments(osArgs)
	want := CommandArgs{
		Host:          "github.com",
		Repo:          "testOrg/testRepo",
		Branch:        "master",
		ScmURL:        "github.com/testOrg/testRepo",
		CloneURL:      "https://stjohn:875fc3f0c3613de2a999295616af7db0fced4056@github.com/testOrg/testRepo.git",
		SHA:           "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		PullRequest:   0,
		TargetDir:     "/tmp/foo",
		CloneMethod:   "https",
		GitName:       "sd-buildbot",
		GitEmail:      "dev-null@screwdriver.cd",
		Version:       false,
		HTTPSUsername: "stjohn",
		HTTPSToken:    "875fc3f0c3613de2a999295616af7db0fced4056",
	}

	if !reflect.DeepEqual(args, want) {
		t.Errorf("Received the wrong arguments: %v, want %v", args, want)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDynamicArgumentsSsh(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "" }
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
		"--clone-method=ssh",
	}

	args, err := GetArguments(osArgs)
	want := CommandArgs{
		Host:        "github.com",
		Repo:        "testOrg/testRepo",
		Branch:      "master",
		ScmURL:      "github.com/testOrg/testRepo",
		CloneURL:    "git@github.com:testOrg/testRepo.git",
		SHA:         "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		PullRequest: 0,
		TargetDir:   "/tmp/foo",
		CloneMethod: "ssh",
		GitName:     "sd-buildbot",
		GitEmail:    "dev-null@screwdriver.cd",
		Version:     false,
	}

	if !reflect.DeepEqual(args, want) {
		t.Errorf("Received the wrong arguments: %v, want %v", args, want)
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestDynamicArgumentsFail(t *testing.T) {
	osGetEnv = func(variable string) (out string) { return "" }
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
		"--clone-method=foobar",
	}

	args, err := GetArguments(osArgs)
	want := CommandArgs{
		Host:        "github.com",
		Repo:        "testOrg/testRepo",
		Branch:      "master",
		ScmURL:      "github.com/testOrg/testRepo",
		CloneURL:    "",
		SHA:         "302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		PullRequest: 0,
		TargetDir:   "/tmp/foo",
		CloneMethod: "foobar",
		GitName:     "sd-buildbot",
		GitEmail:    "dev-null@screwdriver.cd",
		Version:     false,
	}

	if !reflect.DeepEqual(args, want) {
		t.Errorf("Received the wrong arguments: %v, want %v", args, want)
	}

	wantErr := "--clone-method must be https or ssh"
	if err == nil || err.Error() != wantErr {
		t.Errorf("Received the wrong error: %v, want %v", err, wantErr)
	}
}

func TestValidateConfigHost(t *testing.T) {
	osArgs := []string{
		"fakeapp",
		"--repo=testOrg/testRepo",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
	}
	_, err := GetArguments(osArgs)

	wantErr := "--host is required"
	if err == nil || err.Error() != wantErr {
		t.Errorf("Received the wrong error: %v, want %v", err, wantErr)
	}
}

func TestValidateConfigRepo(t *testing.T) {
	osArgs := []string{
		"fakeapp",
		"--host=github.com",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
		"--target-dir=/tmp/foo",
	}
	_, err := GetArguments(osArgs)

	wantErr := "--repo is required"
	if err == nil || err.Error() != wantErr {
		t.Errorf("Received the wrong error: %v, want %v", err, wantErr)
	}
}

func TestValidateConfigSha(t *testing.T) {
	osArgs := []string{
		"fakeapp",
		"--host=github.com",
		"--repo=testOrg/testRepo",
		"--target-dir=/tmp/foo",
	}
	_, err := GetArguments(osArgs)

	wantErr := "--sha is required"
	if err == nil || err.Error() != wantErr {
		t.Errorf("Received the wrong error: %v, want %v", err, wantErr)
	}
}

func TestValidateConfigTargetDir(t *testing.T) {
	osArgs := []string{
		"fakeapp",
		"--host=github.com",
		"--repo=testOrg/testRepo",
		"--sha=302f5f5b48b9feee797a66c88811f1770bcb2dcf",
	}
	_, err := GetArguments(osArgs)

	wantErr := "--target-dir is required"
	if err == nil || err.Error() != wantErr {
		t.Errorf("Received the wrong error: %v, want %v", err, wantErr)
	}
}
