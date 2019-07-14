package main

import (
	"account-sync/core"
	"account-sync/initialize/environment"
	"fmt"
)

func main() {
	fmt.Printf("Press Ctrl+C to exit\n")
	if environment.IsListener() {
		core.RunListener()
	} else {
		core.RunWorker()
	}
}