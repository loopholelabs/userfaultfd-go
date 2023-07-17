package mapper

import (
	"fmt"
	"math"
	"os"
	"syscall"
	"unsafe"

	"github.com/loopholelabs/userfaultfd-go/pkg/constants"
)

type UFFD uintptr

func Register(length int) ([]byte, UFFD, uintptr, error) {
	pagesize := os.Getpagesize()

	uffd, _, errno := syscall.Syscall(constants.NR_userfaultfd, 0, 0, 0)
	if int(uffd) == -1 {
		return []byte{}, 0, 0, fmt.Errorf("%v", errno)
	}

	uffdioAPI := constants.NewUffdioAPI(
		constants.UFFD_API,
		0,
	)

	if _, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		uffd,
		constants.UFFDIO_API,
		uintptr(unsafe.Pointer(&uffdioAPI)),
	); errno != 0 {
		return []byte{}, 0, 0, fmt.Errorf("%v", errno)
	}

	l := int(math.Ceil(float64(length)/float64(pagesize)) * float64(pagesize))
	b, err := syscall.Mmap(
		-1,
		0,
		l,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS,
	)
	if err != nil {
		return []byte{}, 0, 0, fmt.Errorf("%v", errno)
	}

	start := uintptr(unsafe.Pointer(&b[0]))

	uffdioRegister := constants.NewUffdioRegister(
		constants.CULong(start),
		constants.CULong(l),
		constants.UFFDIO_REGISTER_MODE_MISSING,
	)

	if _, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		uffd,
		constants.UFFDIO_REGISTER,
		uintptr(unsafe.Pointer(&uffdioRegister)),
	); errno != 0 {
		return []byte{}, 0, 0, fmt.Errorf("%v", errno)
	}

	return b[:length], UFFD(uffd), start, nil
}
