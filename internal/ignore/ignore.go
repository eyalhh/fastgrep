package ignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func Load(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	return patterns, scanner.Err()

}

func Match(path string, patterns []string) bool {
	if patterns == nil {
		return false
	}
	for _, p := range patterns { 
		matched, err := filepath.Match(p, filepath.Base(path))
		if err == nil && matched {
			return true
		}
	}
	return false
}
