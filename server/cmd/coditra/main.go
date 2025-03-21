package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/fauu/coditra"
	"github.com/fauu/coditra/prepare"
)

func main() {
	if len(os.Args) == 1 {
		err := coditra.RunServer()
		if err != nil {
			printError(err)
			if runtime.GOOS == "windows" {
				_, err := fmt.Scanln()
				if err != nil {
					printError(err)
				}
			}
		}
		return
	}

	if len(os.Args) == 4 && os.Args[1] == "--prepare" {
		inPath := os.Args[2]
		outPath := os.Args[3]
		err := prepare.Do(inPath, outPath)
		if err != nil {
			printError(fmt.Errorf("preparing file %s: %v", inPath, err))
		}
		return
	}

	fmt.Printf(`USAGE:
  coditra                              run the translation companion
  coditra --prepare file.md file.html  convert file.md to file.html prepared for use with the translation companion`)
}

func printError(err error) {
	fmt.Printf("ERROR: %v\n", err)
}
