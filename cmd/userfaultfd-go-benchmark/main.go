package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/loopholelabs/userfaultfd-go/pkg/mapper"
)

type dummyReader struct{}

func (a dummyReader) ReadAt(p []byte, off int64) (n int, err error) {
	n = copy(p, make([]byte, len(p)))

	return n, nil
}

func main() {
	length := flag.Int("chunk", os.Getpagesize()*50, "Bytes of memory to allocate")

	flag.Parse()

	before := time.Now()

	b, uffd, start, err := mapper.Register(*length)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := mapper.Handle(uffd, start, dummyReader{}); err != nil {
			panic(err)
		}
	}()

	n := 0
	for _, c := range b {
		_ = append([]byte{}, c)

		n++
	}

	after := time.Since(before)

	rate := (float64(n) / (1024 * 1024)) / float64(after.Seconds())

	fmt.Printf("%.2f MB/s (%.2f Mb/s)\n", rate, rate*8)
}
