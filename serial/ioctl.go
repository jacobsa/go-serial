package serial

import "golang.org/x/sys/unix"

func ioctl(command int, fd, ret uintptr) error {
	_, _, err := unix.Syscall(unix.SYS_IOCTL, fd, uintptr(command), ret)
	if err != 0 {
		return unix.Errno(err)
	}
	return nil
}
