package main

import (
	"fmt"

	"github.com/younggwon1/k8s-resource-manager/cli"
)

func main() {
	err := cli.RootCmd.Execute()
	if err != nil {
		fmt.Println(err)

	}
}
