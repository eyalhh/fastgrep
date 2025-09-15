package main

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"github.com/eyalhh/fastgrep/internal/search"
	"fmt"
	"os"
)

func main() {
	config, err := cli.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Printf("ignore-case: %t, recursive: %t, show-line-numbers: %t, paths: %v, needle: %s\n\n\nMathces are:\n\n", config.IgnoreCase, config.Recursive, config.ShowLineNumbers, config.Paths, config.Needle)
	file, err := os.Open(config.Paths[0])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	matches, err := search.SearchFile(file, config)
	if err != nil {
		panic(err)
	}

	for index, match := range matches {
		fmt.Printf("\n\n\n=================%d=================\n", index)
		fmt.Printf("line number: %d\n", match.Number)
		for i, r := range match.Ranges {
			fmt.Printf("match number %d: %s\n", i, match.Line[r[0]:r[1]])
		}
	}
}

