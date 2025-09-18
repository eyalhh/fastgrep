package main

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"github.com/eyalhh/fastgrep/internal/search"
	"github.com/eyalhh/fastgrep/pkg/highlight"
	"fmt"
	"os"
)

func main() {
	config, err := cli.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Printf("the options of the search:\n\nignore-case: %t, recursive: %t, show-line-numbers: %t, paths: %v, needle: %s\n\n\nMathces are:\n\n", config.IgnoreCase, config.Recursive, config.ShowLineNumbers, config.Paths, config.Needle)
	file, err := os.Open(config.Paths[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	matches, err := search.SearchFile(file, config)
	if err != nil {
		panic(err)
	}

	for _, match := range matches {
		highlightedLine := highlight.HighlightGreen(match.Line, match.Ranges)
		fmt.Printf("%s:line %d: %s\n", file.Name(), match.Number, highlightedLine)
	}


}

