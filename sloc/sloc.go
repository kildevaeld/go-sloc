package sloc

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type SlocCounter struct {
	files []string
	info  map[string]*Stats
}

func (s *SlocCounter) Add(path string) error {
	files, err := add(path)
	if err != nil {
		return err
	}
	s.files = append(s.files, files...)
	return nil
}

func (s *SlocCounter) Sloc() map[string]*Stats {
	s.info = make(map[string]*Stats)

	for _, file := range s.files {
		s.handleFile(file)
	}
	return s.info
}

//var info = map[string]*Stats{}

func (s *SlocCounter) handleFileLang(fname string, l Language) {
	i, ok := s.info[l.Name()]
	if !ok {
		i = &Stats{}
		s.info[l.Name()] = i
	}
	c, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ! %s\n", fname)
		return
	}
	l.Update(c, i)
}

func (s *SlocCounter) handleFile(fname string) {
	for _, lang := range languages {
		if lang.Match(fname) {
			s.handleFileLang(fname, lang)
			return
		}
	}
	// TODO No recognized extension - check for hashbang
}

const VERSION = `0.3.1`

type Stats struct {
	FileCount    int
	TotalLines   int
	CodeLines    int
	BlankLines   int
	CommentLines int
}

//var info = map[string]*Stats{}
/*
func handleFileLang(fname string, l Language) {
  i, ok := info[l.Name()]
  if !ok {
    i = &Stats{}
    info[l.Name()] = i
  }
  c, err := ioutil.ReadFile(fname)
  if err != nil {
    fmt.Fprintf(os.Stderr, "  ! %s\n", fname)
    return
  }
  l.Update(c, i)
}

func handleFile(fname string) {
  for _, lang := range languages {
    if lang.Match(fname) {
      handleFileLang(fname, lang)
      return
    }
  }
  // TODO No recognized extension - check for hashbang
}

var files []string*/

func add(n string) ([]string, error) {
	fi, err := os.Stat(n)
	if err != nil {
		return nil, err
	}

	var files []string

	if fi.IsDir() {
		fs, err := ioutil.ReadDir(n)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			if f.Name() == ".nosloc" {
				return nil, err
			}
		}
		for _, f := range fs {
			if f.Name()[0] != '.' {
				fs, e := add(path.Join(n, f.Name()))
				if e != nil {
					return nil, e
				}
				files = append(files, fs...)
			}
		}
		return files, nil
	}
	if fi.Mode()&os.ModeType == 0 {
		files = append(files, n)
		return files, nil
	}

	println(fi.Mode())
	return files, nil

}
