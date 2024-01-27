package main

import (
	"fmt"
	"github.com/prog-supdex/mini-project/milestone-code/pkg/cli"
	"os"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
