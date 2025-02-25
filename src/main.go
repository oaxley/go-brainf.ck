package main

// ----- imports
import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ----- constants
const ProgramVersion = "1.0.0"

// ----- functions
func usage() {
	programName := filepath.Base(os.Args[0])

	fmt.Printf("%s - Brainfuck Interpreter - v%s\n", programName, ProgramVersion)
	fmt.Println("Syntax:")
	fmt.Printf("\t%s <filename>\n", programName)
}

func main() {
	// read the command line arguments
	var filename string
	if len(os.Args) < 2 {
		log.Fatalln("Please specify the name of the file to execute!")
		usage()
		os.Exit(1)
	} else {
		filename = os.Args[1]
	}

	// new VM instance
	var core VMCore

	// load a new ROM
	core.Load(filename)

	// execute the ROM
	core.Execute()
}
