package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func executeCommand(args ...string) (out string, err string) {
	cmd := NewRootCmd()
	b := new(bytes.Buffer)
	cmd.SetOut(b)
	cmd.SetArgs(args)
	e := new(bytes.Buffer)
	cmd.SetErr(e)
	cmd.Execute()
	out = b.String()
	err = e.String()
	return out, err
}

func Test_ExecuteCommandCatchErrors(t *testing.T) {
	_, err := executeCommand("h")
	expectedError := `no such file or directory`
	if !strings.Contains(err, expectedError) {
		t.Fatalf("expected \"%s\" got \"%s\"", expectedError, err)
	}
}

func Test_ExecuteCommandWithNoArgs(t *testing.T) {
	out, _ := executeCommand()
	expected := `root.go
root_test.go
`
	if !strings.Contains(out, expected) {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, out)
	}
}

func Test_ExecuteCommandWithDirName(t *testing.T) {
	out, _ := executeCommand("../assets")
	expected := `dew.txt
for_you.txt
rain.txt
`
	if !strings.Contains(out, expected) {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, out)
	}
}

func Test_ExecuteCommandWithFileName(t *testing.T) {
	out, _ := executeCommand("../assets/dew.txt")
	expected := `../assets/dew.txt`
	if !strings.Contains(out, expected) {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, out)
	}
}

func Test_ExecuteCommandWithTwoFileNames(t *testing.T) {
	out, _ := executeCommand("../assets/dew.txt", "../assets/rain.txt")
	expected := `../assets/dew.txt
../assets/rain.txt
`
	if out != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, out)
	}
}

func Test_ExecuteCommandWithFileNameAndDir(t *testing.T) {
	out, _ := executeCommand("../go.mod", "../assets")
	expected := `../go.mod
dew.txt
for_you.txt
rain.txt
`
	if out != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, out)
	}
}