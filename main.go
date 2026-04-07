package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "/?" {
		printHelp()
		return
	}
	initDynamicThings()
	fmt.Printf("Count of characters and pool info with errors: %d\n", checkIntegrity())
	showMainMenu()
}
