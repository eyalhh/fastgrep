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
	if len(os.Args) <= 2 {
		return nil, errors.New("not enough args.")
	}

	var config Config
	re, err := regexp.Compile(os.Args[1])
	if err != nil {
		return nil, err
	}
	config.Pattern = re

	flag.Parse()
	config.IgnoreCase = *ignoreCase
	config.Recursive = *recursive
	config.ShowLineNumbers = *showLineNumbers

	var paths []string
	for _, arg := range os.Args[2:] {
		if arg[0] != '-' {
			paths = append(paths, arg)
		}
	}
	config.Paths = paths

	return &config, nil

}


