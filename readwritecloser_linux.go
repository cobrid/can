package can

import (
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
)

const (
	_SOL_CAN_RAW        = 101
	_CAN_RAW_ERR_FILTER = 2
)

func NewReadWriteCloserForInterface(i *net.Interface) (ReadWriteCloser, error) {
	s, err := syscall.Socket(syscall.AF_CAN, syscall.SOCK_RAW, unix.CAN_RAW)
	if err != nil {
		return nil, err
	}
	allowAllIDsMask := [4]byte{0xff, 0xff, 0xff, 0xff}
	err = syscall.SetsockoptInet4Addr(s, _SOL_CAN_RAW, _CAN_RAW_ERR_FILTER, allowAllIDsMask)
	if err != nil {
		return nil, err
	}

	addr := &unix.SockaddrCAN{Ifindex: i.Index}
	if err := unix.Bind(s, addr); err != nil {
		return nil, err
	}

	f := os.NewFile(uintptr(s), fmt.Sprintf("fd %d", s))

	return &readWriteCloser{f}, nil
}
