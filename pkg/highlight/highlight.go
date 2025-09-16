package highlight

import (
	"github.com/eyalhh/fastgrep/internal/search"
	"sort"
	"strings"
)


const (
	Red = 	"\033[31m"
	Green = "\033[32m"
	Reset = "\033[0m"
)

func HighlightRed(line string, ranges []search.Range) string{
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] <= ranges[j][0];
	})
	var builder strings.Builder
	prev := 0
	for _, r := range ranges {
		builder.WriteString(line[prev:r[0]])
		builder.WriteString(Red)
		builder.WriteString(line[r[0]:r[1]])
		builder.WriteString(Reset)
		prev = r[1]

	}
	builder.WriteString(line[prev:])

	return builder.String()

}

func HighlightGreen(line string, ranges []search.Range) string{
	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i][0] <= ranges[j][0];
	})
	var builder strings.Builder
	prev := 0
	for _, r := range ranges {
		builder.WriteString(line[prev:r[0]])
		builder.WriteString(Green)
		builder.WriteString(line[r[0]:r[1]])
		builder.WriteString(Reset)
		prev = r[1]

	}
	builder.WriteString(line[prev:])

	return builder.String()

}
