package cmd

import (
	"testing"
	"bytes"

	"io/ioutil"

)

func Test_ExecuteCommand(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("test")
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "test" {
		t.Fatalf("expected \"%s\" got \"%s\"", "test", string(out))
	}
}

func Test_ExecuteCommandWithFile(t *testing.T) {
	cmd := NewRootCmd()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{"../assets/rain.txt"})
	cmd.Execute()
	out, err := ioutil.ReadAll(b)
	expected := `“The Taste of Rain” by Jack Kerouac

The taste
Of rain
—Why kneel?`

	if err != nil {
		t.Fatal(err)
	}
	if string(out) != expected {
		t.Fatalf("expected \"%s\" got \"%s\"", expected, string(out))
	}
}