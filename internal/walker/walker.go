package walker


import (
	"fmt"
	"sync"
	"os"
	"path/filepath"
	"github.com/eyalhh/fastgrep/internal/search"
	"github.com/eyalhh/fastgrep/internal/cli"
)
type Walker struct {
	paths []string
	tokens chan struct{}
}

func NewWalker(paths []string, maxWalkers int) *Walker {
	return &Walker{paths: paths, tokens: make(chan struct{}, maxWalkers)}
}


func (w *Walker) Walk(conf *cli.Config, matches chan []search.Match) {

	var n sync.WaitGroup

	for _, root := range w.paths {
		fi, err := os.Stat(root)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("file does not exist")
			}
			panic(err)
		}
		if fi.Mode().IsRegular() {
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
			go w.walkDir(root, &n, conf, matches)
		}
	}

	go func() {
		n.Wait()
		close(matches)
	}()
}

func (w *Walker) walkDir(dir string, n *sync.WaitGroup, conf *cli.Config, matches chan []search.Match) error {
	defer n.Done()
	entries, err := w.dirents(dir)
	if err != nil {
		panic(err)
		return err
	}
	for _, entry := range entries {
		fullpath := filepath.Join(dir, entry.Name())
		if entry.IsDir() {
			n.Add(1)
			go w.walkDir(fullpath, n, conf, matches)
		} else {
			file, err := os.Open(fullpath)
			if err != nil {
				panic(err)
				return err
			}
			foundMatches, err := search.SearchFile(file, conf)
			if err != nil {
				panic(err)
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
		panic(err)
		return nil, fmt.Errorf("walker: %v\n", err)
	}
	return entries, nil
}


