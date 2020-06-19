// +build windows

package readpassword

import (
	"io"
	"syscall"
	"unsafe"
)

const (
	enableLineInput       = 2
	enableEchoInput       = 4
	enableProcessedInput  = 1
	enableProcessedOutput = 1
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = kernel32.NewProc("SetConsoleMode")
)

// https://github.com/Nerdmaster/terminal/blob/master/util_linux.go

func readPassword(fd int) ([]byte, error) {
	var st uint32
	_, _, e := syscall.Syscall(procGetConsoleMode.Addr(), 2, uintptr(fd), uintptr(unsafe.Pointer(&st)), 0)
	if e != 0 {
		return nil, error(e)
	}
	old := st

	st &^= (enableEchoInput)
	st |= (enableProcessedInput | enableLineInput | enableProcessedOutput)
	_, _, e = syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(st), 0)
	if e != 0 {
		return nil, error(e)
	}

	defer func() {
		syscall.Syscall(procSetConsoleMode.Addr(), 2, uintptr(fd), uintptr(old), 0)
	}()

	var buf [16]byte
	var ret []byte
	for {
		n, err := syscall.Read(syscall.Handle(fd), buf[:])
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
		if n > 0 && buf[n-1] == '\r' {
			n--
		}
		ret = append(ret, buf[:n]...)
		if n < len(buf) {
			break
		}
	}

	return ret, nil
}
