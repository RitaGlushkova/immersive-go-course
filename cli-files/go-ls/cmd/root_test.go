package cmd

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
	"strings"
	"github.com/spf13/cobra"
)
func execute(t *testing.T, c *cobra.Command, args ...string) (string, error) {
  t.Helper()

  buf := new(bytes.Buffer)
  c.SetOut(buf)
  c.SetErr(buf)
  c.SetArgs(args)

  err := c.Execute()
  return strings.TrimSpace(buf.String()), err
}

// func Test_ExecuteCommand(t *testing.T) {
// 	cmd := NewRootCmd()
// 	b := bytes.NewBufferString("")
// 	cmd.SetOut(b)
// 	cmd.SetArgs([]string{})
// 	err := cmd.Execute()
// 	expected := `root.go
// root_test.go
// `
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	require.Equal(t, expected, b.String())
// }

// func Test_ExecuteCommandWithDirName(t *testing.T) {
// 	cmd := NewRootCmd()
// 	b := bytes.NewBufferString("")
// 	cmd.SetOut(b)
// 	cmd.SetArgs([]string{"../assets"})
// 	cmd.Execute()
// 	_, err := cmd.ExecuteC()
// 	expected := `dew.txt
// for_you.txt
// rain.txt
// `
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	require.Equal(t, expected, b.String())
// }

func Test_ExecuteCommand(t *testing.T) {
  tt := []struct {
    args []string
    err  error
    out  string
  }{
    {
      args: nil,
      err:  nil,
	  out:  `root.go
root_test.go
`,
    },
    {
      args: []string{"../assets"},
      err:  nil,
      out: `dew.txt
for_you.txt
rain.txt
`,
    },
    {
      args: []string{"../assets/dew.txt"},
      err:  nil,
      out:  `dew.txt
`,
    },
  }
  cmd := NewRootCmd()
  for _, tc := range tt {
    out, err := execute(t, cmd, tc.args...)

    require.Equal(t, tc.err, err)

    if tc.err == nil {
      require.Equal(t, tc.out, out)
    }
  }
// 	cmd := NewRootCmd()
// 	b := bytes.NewBufferString("")
// 	cmd.SetOut(b)
// 	cmd.SetArgs([]string{"../assets/dew.txt"})
// 	err := cmd.Execute()
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
// 	expected := `dew.txt
// `
// 	require.Equal(t, expected, b.String())
}
