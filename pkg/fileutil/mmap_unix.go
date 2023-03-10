//go:build !windows

package fileutil

import (
	"os"

	"golang.org/x/sys/unix"
)

func mmap(fd int, offset int64, size, mode int) ([]byte, error) {
	var prot int
	if mode&read != 0 {
		prot |= unix.PROT_READ
	}

	if mode&write != 0 {
		prot |= unix.PROT_WRITE
	}

	data, err := unix.Mmap(fd, offset, size, prot, unix.MAP_SHARED)
	return data, err
}

func munmap(_ *os.File, data []byte) error {
	return unix.Munmap(data)
}

func msync(data []byte) error {
	return unix.Msync(data, unix.MS_SYNC)
}
