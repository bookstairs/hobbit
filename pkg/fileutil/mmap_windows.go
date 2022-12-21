//go:build windows

package fileutil

import (
	"errors"
	"os"
	"reflect"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

type MMap []byte

func (m *MMap) header() *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(m))
}

func (m *MMap) addrLen() (data, length uintptr) {
	header := m.header()
	return header.Data, uintptr(header.Len)
}

type mapHandle struct {
	file     windows.Handle
	view     windows.Handle
	writable bool
}

var handleMap = make(map[uintptr]*mapHandle)

var lock4map sync.Mutex

func mmap(fd int, offset int64, size, mode int) ([]byte, error) {
	prot := windows.PAGE_READONLY
	access := windows.FILE_MAP_READ
	writable := false
	if mode&write != 0 {
		prot = windows.PAGE_READWRITE
		access = windows.FILE_MAP_WRITE
		writable = true
	}

	// The maximum size is the area of the file, starting from 0,
	// that we wish to allow to be mappable. It is the sum of
	// the length the user requested, plus the offset where that length
	// is starting from. This does not map the data into memory.
	maxSizeHigh := uint32((offset + int64(size)) >> 32)
	maxSizeLow := uint32((offset + int64(size)) & 0xFFFFFFFF)
	// TODO: Do we need to set some security attributes? It might help portability.
	h, errno := windows.CreateFileMapping(windows.Handle(uintptr(fd)), nil, uint32(prot), maxSizeHigh, maxSizeLow, nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}

	// Actually map a view of the data into memory. The view's size
	// is the length the user requested.
	fileOffsetHigh := uint32(offset >> 32)
	fileOffsetLow := uint32(offset & 0xFFFFFFFF)
	addr, errno := windows.MapViewOfFile(h, uint32(access), fileOffsetHigh, fileOffsetLow, uintptr(size))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	lock4map.Lock()
	defer lock4map.Unlock()
	handleMap[addr] = &mapHandle{
		file:     windows.Handle(uintptr(fd)),
		view:     h,
		writable: writable,
	}

	mmap := MMap{}

	hd := mmap.header()
	hd.Data = addr
	hd.Len = size
	hd.Cap = hd.Len

	return mmap, nil
}

func munmap(f *os.File, bytes []byte) error {
	defer func() {
		// if not close file, when remove file will throw file be used other process.
		_ = f.Close()
	}()
	mmap := MMap(bytes)
	hd := mmap.header()
	addr := hd.Data
	// Lock the UnmapViewOfFile along with the handleMap deletion.
	// As soon as we unmap the view, the OS is free to give the
	// same addr to another new map. We don't want another goroutine
	// to insert and remove the same addr into handleMap while
	// we're trying to remove our old addr/handle pair.
	lock4map.Lock()
	defer lock4map.Unlock()
	err := windows.UnmapViewOfFile(addr)
	if err != nil {
		return err
	}

	handle, ok := handleMap[addr]
	if !ok {
		// should be impossible; we would've errored above
		return errors.New("unknown base address")
	}
	delete(handleMap, addr)

	e := windows.CloseHandle(handle.view)
	return os.NewSyscallError("CloseHandle", e)
}

func msync(bytes []byte) error {
	mmap := MMap(bytes)
	addr, size := mmap.addrLen()
	errno := windows.FlushViewOfFile(addr, size)
	if errno != nil {
		return os.NewSyscallError("FlushViewOfFile", errno)
	}

	lock4map.Lock()
	defer lock4map.Unlock()

	handle, ok := handleMap[addr]
	if !ok {
		// should be impossible; we would've errored above
		return errors.New("unknown base address")
	}

	if handle.writable && handle.file != windows.Handle(^uintptr(0)) {
		if err := windows.FlushFileBuffers(handle.file); err != nil {
			return os.NewSyscallError("FlushFileBuffers", err)
		}
	}

	return nil
}
