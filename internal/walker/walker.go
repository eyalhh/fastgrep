package walker


import (
	"fmt"
	"sync"
	"os"
	"path/filepath"
	"github.com/eyalhh/fastgrep/internal/search"
	"github.com/eyalhh/fastgrep/internal/cli"
	"github.com/eyalhh/fastgrep/internal/ignore"
)

func isExecutable(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	mode := fi.Mode()
	return mode.IsRegular() && (mode&0111!=0)

}
type Walker struct {
	paths []string
	tokens chan struct{}
}

func NewWalker(paths []string, maxWalkers int) *Walker {
	return &Walker{paths: paths, tokens: make(chan struct{}, maxWalkers)}
}


func (w *Walker) Walk(conf *cli.Config, matches chan []search.Match) {

	patterns, err := ignore.Load(".grepignore")
	if err != nil {
		fmt.Println(err)
	}

	var n sync.WaitGroup

	for _, root := range w.paths {

		if ignore.Match(root, patterns) {
			continue
		}
		fi, err := os.Stat(root)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("file does not exist")
			}
			panic(err)
		}
		if fi.Mode().IsRegular() {
			if isExecutable(root) {
				continue
			}
			n.Add(1)
			go func(root string) {
				defer n.Done()
				file, err := os.Open(root)
				if err != nil {
					return
				}
				foundMatches, err := search.SearchFile(file, conf)
				if err != nil {
					return
				}
				matches <- foundMatches

			}(root)
		} else {
			n.Add(1)
			go w.walkDir(root, &n, conf, matches, patterns)
		}
	}

	go func() {
		n.Wait()
		close(matches)
	}()
}

func (w *Walker) walkDir(dir string, n *sync.WaitGroup, conf *cli.Config, matches chan []search.Match, patterns []string) error {
	defer n.Done()
	entries, err := w.dirents(dir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, entry := range entries {
		if ignore.Match(entry.Name(), patterns) {
			continue
		}
		fullpath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			n.Add(1)
			go w.walkDir(fullpath, n, conf, matches, patterns)
		} else {
			if isExecutable(entry.Name()) {
				continue
			}
			file, err := os.Open(fullpath)
			if err != nil {
				fmt.Println(err)
				return err
			}
			foundMatches, err := search.SearchFile(file, conf)
			if err != nil {
				fmt.Println(err)
				return err
			}
			matches <- foundMatches
		}

	}
	return nil
}

func (w *Walker) dirents(dir string) ([]os.DirEntry, error) {
	w.tokens <- struct{}{}
	defer func() { <- w.tokens }()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("walker: %v\n", err)
	}
	return entries, nil
}


