package cmd

import (
	"testing"
	"bytes"
	"io/ioutil"
)

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	expected := `root.go
root_test.go
`
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected , string(out))
	}
}

func Test_ExecuteCommandWithDirName(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"../assets"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	expected := `dew.txt
for_you.txt
rain.txt
`
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}

func Test_ExecuteCommandWithFileName(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"../assets/dew.txt"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	expected := `dew.txt
`
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}