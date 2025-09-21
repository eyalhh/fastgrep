package main

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"time"
	"strings"
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

	fmt.Printf("\nthe options of the search:\nignore-case: %t, recursive: %t, show-line-numbers: %t, paths: %v\n\n\nMathces are:\n\n", config.IgnoreCase, config.Recursive, config.ShowLineNumbers, config.Paths)


	walker := walker.NewWalker(config.Paths, 20)
	matches := make(chan []search.Match)
	start := time.Now()
	walker.Walk(config, matches)
	
	
	for fileMatches := range matches {
		for _, match := range fileMatches {
			highlightedLine := strings.TrimSpace(highlight.HighlightGreen(match.Line, match.Ranges))
			header := fmt.Sprintf("%s: line %d: ", match.FileName, match.Number)
			highlightedHeader := highlight.HighlightRed(header, []search.Range{{0, len(header)}})
			fmt.Printf("%s\n%s\n", highlightedHeader, highlightedLine)
		}
	}

	fmt.Printf("\n\ntime it took for the search: %v\n", time.Since(start).Seconds())


}

