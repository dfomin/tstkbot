package main

import (
	"bufio"
	"os"
	"tstkbot/commands"
)

// InitCLI inits command line interface mode. Bot works locally and ready
// messages from command line instead of network.
func InitCLI() {
	reader := bufio.NewReader(os.Stdin)
	for {
		command, _ := reader.ReadString('\n')
		commands.ProcessMessage(command)
	}
}
