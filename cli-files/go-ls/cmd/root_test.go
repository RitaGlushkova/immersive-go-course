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
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "testisawesome" {
		t.Fatalf("expected \"%s\" got \"%s\"", "testisawesome", string(out))
	}
}