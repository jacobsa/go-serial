package serial

import "golang.org/x/sys/unix"

// ioctl provides a wrapper around the unix.Syscall, returning errors that general go code can deal with
func ioctl(command int, fd, ret uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(command), ret)
	if err != 0 {
		return unix.Errno(err)
	}
	return nil
}
