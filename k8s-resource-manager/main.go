package main

import (
	"fmt"
	"os"

	"github.com/younggwon1/k8s-resource-manager/cli"
)

func main() {
	err := cli.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
