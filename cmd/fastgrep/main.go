package main

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"github.com/eyalhh/fastgrep/internal/search"
	"github.com/eyalhh/fastgrep/internal/walker"
	"github.com/eyalhh/fastgrep/pkg/highlight"
	"fmt"
)

func main() {
	config, err := cli.GetConfig()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	walker := walker.NewWalker(config.Paths, 20)
	matches := make(chan []search.Match)
	walker.Walk(config, matches)
	
	
	fmt.Printf("the options of the search:\n\nignore-case: %t, recursive: %t, show-line-numbers: %t, paths: %v\n\n\nMathces are:\n\n", config.IgnoreCase, config.Recursive, config.ShowLineNumbers, config.Paths)

	for fileMatches := range matches {
		for _, match := range fileMatches {
			highlightedLine := highlight.HighlightGreen(match.Line, match.Ranges)
			fmt.Printf("%s: line %d: %s\n", match.FileName, match.Number, highlightedLine)
		}
	}


}

