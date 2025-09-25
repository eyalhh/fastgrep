package search

import (
	"github.com/eyalhh/fastgrep/internal/cli"
	"fmt"
	"os"
	"io"
	"strings"
	"bufio"
)

type Range [2]int

type Match struct {
	FileName string
	Line string
	Ranges []Range
	Number int
}

func SearchFile(r io.Reader, conf *cli.Config) ([]Match, error) {
	var result []Match
	scanner := bufio.NewScanner(r)

	const maxCapacity = 30 * 1024 * 1024
	buf := make([]byte, 0, 64 * 1024)
	scanner.Buffer(buf, maxCapacity)

	counter := 0
	for scanner.Scan() {
		// define the currentmatch here for zeroing out the attributes every loop
		var currentMatch Match
		if res, ok := r.(*os.File); ok {
			currentMatch.FileName = res.Name()
		}
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
		if res, ok := r.(*os.File); ok {
			fmt.Println(res.Name())
		}
		return nil, err
	}

	return result, nil

}
