package mapper

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"

	"github.com/loopholelabs/userfaultfd-go/pkg/constants"
	"golang.org/x/sys/unix"
)

var (
	ErrUnexpectedEventType = errors.New("unexpected event type")
)

func Handle(uffd UFFD, start uintptr, src io.ReaderAt) error {
	pagesize := os.Getpagesize()

	for {
		if _, err := unix.Poll(
			[]unix.PollFd{{
				Fd:     int32(uffd),
				Events: unix.POLLIN,
			}},
			-1,
		); err != nil {
			return err
		}

		buf := make([]byte, unsafe.Sizeof(constants.UffdMsg{}))
		if _, err := syscall.Read(int(uffd), buf); err != nil {
			return err
		}

		msg := (*(*constants.UffdMsg)(unsafe.Pointer(&buf[0])))
		if constants.GetMsgEvent(&msg) != constants.UFFD_EVENT_PAGEFAULT {
			return ErrUnexpectedEventType
		}

		arg := constants.GetMsgArg(&msg)
		pagefault := (*(*constants.UffdPagefault)(unsafe.Pointer(&arg[0])))

		addr := constants.GetPagefaultAddress(&pagefault)

		p := make([]byte, pagesize)
		if _, err := src.ReadAt(p, int64(uintptr(addr)-start)); err != nil {
			return err
		}

		cpy := constants.NewUffdioCopy(
			p,
			addr&^constants.CULong(pagesize-1),
			constants.CULong(pagesize),
			0,
			0,
		)

		if _, _, errno := syscall.Syscall(
			syscall.SYS_IOCTL,
			uintptr(uffd),
			constants.UFFDIO_COPY,
			uintptr(unsafe.Pointer(&cpy)),
		); errno != 0 {
			return fmt.Errorf("%v", errno)
		}
	}
}
