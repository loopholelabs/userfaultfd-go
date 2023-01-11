package transfer

import (
	"encoding/binary"
	"errors"
	"log"
	"net"
	"syscall"

	"github.com/loopholelabs/userfaultfd-go/pkg/mapper"
)

var (
	ErrInvalidFDsLength = errors.New("invalid file descriptors length")
)

func ReceiveUFFD(conn *net.UnixConn) (mapper.UFFD, uintptr, error) {
	var start int64
	if err := binary.Read(conn, binary.BigEndian, &start); err != nil {
		return 0, 0, err
	}

	fds, err := ReceiveFds(conn, 1)
	if err != nil {
		return 0, 0, err
	}

	if len(fds) != 2 {
		return 0, 0, ErrInvalidFDsLength
	}

	return mapper.UFFD(fds[1]), uintptr(start), nil
}

func ReceiveFds(conn *net.UnixConn, num int) ([]uintptr, error) {
	if num <= 0 {
		return []uintptr{}, nil
	}

	f, err := conn.File()
	if err != nil {
		return []uintptr{}, err
	}
	defer f.Close()

	buf := make([]byte, syscall.CmsgSpace(num*4)) // See https://github.com/ftrvxmtrx/fd/blob/master/fd.go#L51
	if _, _, _, _, err = syscall.Recvmsg(int(f.Fd()), nil, buf, 0); err != nil {
		return []uintptr{}, err
	}

	msgs, err := syscall.ParseSocketControlMessage(buf)
	if err != nil {
		return []uintptr{}, err
	}

	fds := []uintptr{}
	for _, msg := range msgs {
		newFds, err := syscall.ParseUnixRights(&msg)
		if err != nil {
			return []uintptr{}, err
		}

		for _, newFd := range newFds {
			fds = append(fds, uintptr(newFd))
		}
	}

	return fds, nil
}

func SendUFFD(conn *net.UnixConn, uffd mapper.UFFD, start uintptr) error {
	log.Println(start, uffd)

	if err := binary.Write(conn, binary.BigEndian, int64(start)); err != nil {
		return err
	}

	return SendFds(conn, []uintptr{uintptr(uffd)})
}

func SendFds(conn *net.UnixConn, fds []uintptr) error {
	if len(fds) <= 0 {
		return nil
	}

	f, err := conn.File()
	if err != nil {
		return err
	}
	defer f.Close()

	b := make([]int, len(fds))
	for _, fd := range fds {
		b = append(b, int(fd))
	}

	return syscall.Sendmsg(int(f.Fd()), nil, syscall.UnixRights(b...), nil, 0)
}
