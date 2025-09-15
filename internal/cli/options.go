package cli

import (
	"os"
	"flag"
	"errors"
	"regexp"
)

type Config struct {

	Pattern 		*regexp.Regexp
	Needle 			string
	Paths 			[]string
	IgnoreCase 		bool
	Recursive 		bool
	ShowLineNumbers bool

}

var ignoreCase = flag.Bool("i", false, "ignore case when matching")
var recursive = flag.Bool("r", false, "recursive matching")
var showLineNumbers = flag.Bool("n", false, "show line numbers matching")
var enableRegex = flag.Bool("regex", false, "enable regex")

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
	if startingIndex == 0 {
		startingIndex = len(os.Args) - 1
	}
	if len(os.Args[startingIndex:]) < 2 {
		return nil, errors.New("not enough args.")
	}

	var config Config
	if *enableRegex {
		var re *regexp.Regexp
		var err error
		if *ignoreCase {
			re, err = regexp.Compile("(?i)" + os.Args[startingIndex])

		} else {
			re, err = regexp.Compile(os.Args[startingIndex])
		}
		if err != nil {
			return nil, err
		}
		config.Pattern = re
	} else {
		config.Needle = os.Args[startingIndex]
	}

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

