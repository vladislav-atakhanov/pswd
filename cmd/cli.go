package main

import (
	"fmt"
	"github.com/vladislav-atakhanov/pswd/cmd/cli"
	"os"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
