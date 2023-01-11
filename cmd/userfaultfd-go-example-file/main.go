package main

import (
	"flag"
	"log"
	"os"

	"github.com/loopholelabs/userfaultfd-go/pkg/mapper"
)

func main() {
	length := flag.Int("length", os.Getpagesize()*2, "Amount of bytes to allocate")
	file := flag.String("file", "LICENSE", "File to map")

	flag.Parse()

	f, err := os.OpenFile(*file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	b, uffd, start, err := mapper.Register(*length)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := mapper.Handle(uffd, start, f); err != nil {
			panic(err)
		}
	}()

	for i, c := range b {
		log.Printf("%v %c", i, rune(c))
	}
}
