package readpassword

// const ioctlReadTermios = 0x5401  // syscall.TCGETS
// const ioctlWriteTermios = 0x5402 // syscall.TCSETS

import "golang.org/x/sys/unix"

const ioctlReadTermios = unix.TCGETS
const ioctlWriteTermios = unix.TCSETS
