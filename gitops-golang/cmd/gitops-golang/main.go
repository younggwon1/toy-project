package main

import (
	"fmt"
	"os"

	"github.com/younggwon1/gitops-golang/cli"
)

func main() {
	err := cli.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
