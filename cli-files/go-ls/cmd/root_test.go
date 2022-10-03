package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
	
)

// rename err if it is not type is error
func executeCommand(args ...string) (out string, stderr string, err error) {
	cmd := NewRootCmd()
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetArgs(args)	
	e := new(bytes.Buffer)
	cmd.SetErr(e)
	err = cmd.Execute()
	out = b.String()
	stderr = e.String()
	return out, stderr, err
}
func assertContains(t *testing.T, str, expected string) {
	t.Helper()
	if !strings.Contains(str, expected) {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, str)
	}
}

// func assertError(t *testing.T, err error) {
// 	t.Helper()
// 	if (!os.IsNotExist(err)) && (err != nil) {
// 		t.Fatalf("Unexpected error, want nil, got %v", err)
// 	}
// }

func Test_ExecuteCommandCatchErrors(t *testing.T) {
	_, stderr, err := executeCommand("h")
	if !os.IsNotExist(err) {
		t.Fatalf("Wrong error, expected no such file or directory, got %v", err)
	}
	expectedError := `no such file or directory`
	assertContains(t, stderr, expectedError)

}

func Test_ExecuteCommandWithNoArgs(t *testing.T) {
	out, stderr, err := executeCommand()
	if err != nil {
		t.Fatalf("Could not execute command %v", stderr)
	}
	expected := `root.go
root_test.go
`
	assertContains(t, out, expected)
}

func Test_ExecuteCommandWithDirName(t *testing.T) {
	out, _, err := executeCommand("../assets")
	if err != nil {
		t.Fatalf("Could not execute command %v", err)
	}
	expected := `dew.txt
for_you.txt
rain.txt
`
	assertContains(t, out, expected)
}

func Test_ExecuteCommandWithFileName(t *testing.T) {
	out, _, err := executeCommand("../assets/dew.txt")
	if err != nil {
		t.Fatalf("Could not execute command %v", err)
	}
	expected := `../assets/dew.txt`
	assertContains(t, out, expected)
}

func Test_ExecuteCommandWithTwoFileNames(t *testing.T) {
	out, _, err := executeCommand("../assets/dew.txt", "../assets/rain.txt")
	if err != nil {
		t.Fatalf("Could not execute command %v", err)
	}
	expected := `../assets/dew.txt
../assets/rain.txt
`
	assertContains(t, out, expected)
}

func Test_ExecuteCommandWithFileNameAndDir(t *testing.T) {
	out, _, _ := executeCommand("../go.mod", "../assets")
	expected := `../go.mod
dew.txt
for_you.txt
rain.txt
`
	assertContains(t, out, expected)
}
