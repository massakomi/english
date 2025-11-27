package main

import (
	"english/cmd"
	"english/test"
	"os"
)

func main() {
	osArgs := os.Args[1:]
	if len(osArgs) > 0 {
		test.TestGo()
	} else {
		cmd.Run()
	}

}
