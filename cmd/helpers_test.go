package cmd

import "bytes"

func executeCommand(args ...string) (string, error) {
	cmd := rootCmd
	cmd.SetArgs(args)

	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := cmd.Execute()

	return buf.String(), err
}
