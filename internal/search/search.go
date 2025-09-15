package search

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"fmt"
	"io"
	"strings"
	"bufio"
)

type Range [2]int

type Match struct {
	Line string
	Ranges []Range
	Number int
}

func SearchFile(r io.Reader, conf *cli.Config) ([]Match, error) {
	var result []Match
	scanner := bufio.NewScanner(r)

	counter := 0
	for scanner.Scan() {
		var currentMatch Match
		counter++
		line := scanner.Text()
		if conf.Pattern != nil {
			// need to also recompile pattern with ?i if case insensitive enabled
			matches := conf.Pattern.FindAllStringIndex(line, -1)
			if len(matches) > 0{
				currentMatch.Number = counter
				for _, match := range matches {
					currentRange := Range([2]int{match[0], match[1]})
					currentMatch.Ranges = append(currentMatch.Ranges, currentRange)
				}
				currentMatch.Line = line
				result = append(result, currentMatch)
			}

		} else {
			if conf.IgnoreCase {
				if strings.Contains(strings.ToLower(line), strings.ToLower(conf.Needle)) {
					index := strings.Index(strings.ToLower(line), strings.ToLower(conf.Needle))
					currentMatch.Number = counter
					currentMatch.Ranges = []Range{Range([2]int{index, index + len(conf.Needle)})}
					currentMatch.Line = line
					result = append(result, currentMatch)
				}				
			} else {

				if strings.Contains(line, conf.Needle) {
					index := strings.Index(line, conf.Needle)
					currentMatch.Number = counter
					currentMatch.Ranges = []Range{Range([2]int{index, index + len(conf.Needle)})}
					currentMatch.Line = line
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
