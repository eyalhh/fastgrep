package cli

import (
	"os"
	"flag"
	"fmt"
	"errors"
	"regexp"
)

type Config struct {

	Pattern 		*regexp.Regexp
	Paths 			[]string
	IgnoreCase 		bool
	Recursive 		bool
	ShowLineNumbers bool

}

var ignoreCase = flag.Bool("i", false, "ignore case when matching")
var recursive = flag.Bool("r", false, "recursive matching")
var showLineNumbers = flag.Bool("n", false, "show line numbers matching")

func GetConfig() (*Config, error) {
	var startingIndex int
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg[0] != '-' {
			startingIndex = i
			break
		}
	}
	if len(os.Args[startingIndex:]) < 2 {
		return nil, errors.New("not enough args.")
	}

	var config Config
	re, err := regexp.Compile(os.Args[startingIndex])
	if err != nil {
		return nil, err
	}
	config.Pattern = re

	flag.Parse()
	config.IgnoreCase = *ignoreCase
	config.Recursive = *recursive
	config.ShowLineNumbers = *showLineNumbers

	var paths []string
	for _, arg := range os.Args[startingIndex+1:] {
		paths = append(paths, arg)
	}
	config.Paths = paths

	return &config, nil

}

