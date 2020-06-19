// +build aix darwin dragonfly freebsd linux,!appengine netbsd openbsd

package readpassword

import (
	"io"
	"syscall"
	"unsafe"
)

func readPassword(fd int) ([]byte, error) {
	var oldState syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlReadTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	newState := oldState
	newState.Lflag &^= syscall.ECHO
	newState.Lflag |= syscall.ICANON | syscall.ISIG
	newState.Iflag |= syscall.ICRNL
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&newState)), 0, 0, 0); err != 0 {
		return nil, err
	}

	defer func() {
		syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), ioctlWriteTermios, uintptr(unsafe.Pointer(&oldState)), 0, 0, 0)
	}()

	var buf [16]byte
	var ret []byte
	for {
		n, err := syscall.Read(fd, buf[:])
		if err != nil {
			return nil, err
		}
		if n == 0 {
			if len(ret) == 0 {
				return nil, io.EOF
			}
			break
		}
		if buf[n-1] == '\n' {
			n--
		}
		ret = append(ret, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	return ret, nil
}
