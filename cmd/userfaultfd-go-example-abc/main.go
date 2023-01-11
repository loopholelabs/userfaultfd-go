package main

import (
	"bytes"
	"errors"
	"flag"
	"log"
	"net"
	"os"
	"strings"

	"github.com/loopholelabs/userfaultfd-go/pkg/mapper"
	"github.com/loopholelabs/userfaultfd-go/pkg/transfer"
)

type abcReader struct{}

func (a abcReader) ReadAt(p []byte, off int64) (n int, err error) {
	log.Println("Reading at offset", off)

	n = copy(p, bytes.Repeat([]byte{'A' + byte(off%20)}, len(p)))

	return n, nil
}

func main() {
	socket := flag.String("socket", "userfaultd.sock", "Socket to share the file descriptor over")
	server := flag.Bool("server", false, "Whether to serve as the server instead of the client")
	length := flag.Int("length", os.Getpagesize()*2, "Amount of bytes to allocate")

	flag.Parse()

	if strings.TrimSpace(*socket) == "" {
		panic(errors.New("could not work with empty socket path"))
	}

	addr, err := net.ResolveUnixAddr("unix", *socket)
	if err != nil {
		panic(err)
	}

	if *server {
		lis, err := net.ListenUnix("unix", addr)
		if err != nil {
			panic(err)
		}

		log.Println("Listening on", addr.String())

		for {
			conn, err := lis.AcceptUnix()
			if err != nil {
				panic(err)
			}

			go func() {
				defer func() {
					if err := recover(); err != nil {
						log.Println("Could not handle connection, stopping:", err)
					}

					_ = conn.Close()
				}()

				uffd, start, err := transfer.ReceiveUFFD(conn)
				if err != nil {
					panic(err)
				}

				if err := mapper.Handle(uffd, start, abcReader{}); err != nil {
					panic(err)
				}
			}()
		}
	} else {
		conn, err := net.DialUnix("unix", nil, addr)
		if err != nil {
			panic(err)
		}

		log.Println("Connected to", conn.RemoteAddr())

		b, uffd, start, err := mapper.Register(*length)
		if err != nil {
			panic(err)
		}

		if err := transfer.SendUFFD(conn, uffd, start); err != nil {
			panic(err)
		}

		for i, c := range b {
			log.Printf("%v %c", i, rune(c))
		}
	}
}
