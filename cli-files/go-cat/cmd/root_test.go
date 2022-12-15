package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	errB := bytes.NewBufferString("")
	cmd.SetErr(errB)
	cmd.SetArgs([]string{})
	cmd.Execute()
	err := errB.String()
	wantError := "Error: missing parameter, provide file name"
	if !strings.Contains(err, wantError) {
		t.Fatalf("expected \"%s\" got \"%s\"", wantError, err)
	}
}

func Test_ExecuteCommandWithFile(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	errB := bytes.NewBufferString("")
	cmd.SetErr(errB)
	err := errB.String()
	cmd.SetOut(b)
	cmd.SetArgs([]string{"../assets/rain.txt"})
	cmd.Execute()
	out := b.String()
	expected := `“The Taste of Rain” by Jack Kerouac

The taste
Of rain
—Why kneel?`
	if out != expected {
		t.Fatalf("expected \"%s\" got \"%s\" stderr = %s", expected, out, err)
	}
}
