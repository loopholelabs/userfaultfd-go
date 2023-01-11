package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/loopholelabs/userfaultfd-go/pkg/mapper"
)

func main() {
	file := flag.String("file", "LICENSE", "File to map")
	dst := flag.String("dst", filepath.Join(os.TempDir(), "LICENSE"), "Destination to write to")

	flag.Parse()

	f, err := os.OpenFile(*file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	s, err := f.Stat()
	if err != nil {
		panic(err)
	}

	b, uffd, start, err := mapper.Register(int(s.Size()))
	if err != nil {
		panic(err)
	}

	go func() {
		if err := mapper.Handle(uffd, start, f); err != nil {
			panic(err)
		}
	}()

	if err := os.WriteFile(*dst, b, os.ModePerm); err != nil {
		panic(err)
	}
}
