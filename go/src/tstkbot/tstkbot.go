package main

import (
	"os"
)

func main() {
	cliMode := false
	if len(os.Args) > 1 {
		if os.Args[1] == "--cli" {
			cliMode = true
		}
	}

	//InitDatabase(databaseName)

	if cliMode {
		InitCLI()
	} else {
		//InitServer()
	}
}
