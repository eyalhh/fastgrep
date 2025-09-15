package search

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"io"
	"strings"
	"bufio"
)

type Match struct {
	line string
	number int
	offset int
}

func SearchFile(r io.Reader, conf *cli.Config) ([]Match, error) {
	var result []Match
	var currentMatch Match
	scanner := bufio.NewScanner(r)

	counter := 0
	for scanner.Scan() {
		counter++
		line := scanner.Text()
		if conf.Pattern != nil {
			// need to also recompile pattern with ?i if case insensitive enabled
			match := conf.Pattern.MatchString(line)
			if match {
				currentMatch.line = line
				currentMatch.number = counter
				currentMatch.offset = 0
				result = append(result, currentMatch)
			}

		} else {
			if conf.IgnoreCase {
				if strings.Contains(strings.ToLower(line), strings.ToLower(conf.Needle)) {
					currentMatch.line = line
					currentMatch.number = counter
					currentMatch.offset = 0
					result = append(result, currentMatch)
				}				
			} else {

				if strings.Contains(line, conf.Needle) {
					currentMatch.line = line
					currentMatch.number = counter
					currentMatch.offset = 0
					result = append(result, currentMatch)

				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil

}
