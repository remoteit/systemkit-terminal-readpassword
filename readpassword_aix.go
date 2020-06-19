// +build aix

package readpassword

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TCGETS
const ioctlWriteTermios = unix.TCSETS
