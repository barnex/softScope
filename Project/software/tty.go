package softscope

// Access to serial device

import (
	"errors"
	"log"
	"strconv"
	"unsafe"
)

////cgo CFLAGS: -Wall -Werror
//#include "tty.h"
import "C"

type TTY int

func InitTTY(ttyDev, baudrate string) TTY {
	baud, err := strconv.Atoi(baudrate)
	if err != nil {
		log.Fatal("invalid baud rate:", baudrate)
	}
	log.Println("Using", ttyDev, "@", baud, "baud")

	tty, err2 := OpenTTY(ttyDev, baud)
	if err2 != nil {
		log.Fatal(err2)
	}
	return tty
}

func OpenTTY(fname string, baud int) (TTY, error) {
	fd := C.openTTY(C.CString(fname), C.int(baud))
	if fd < 0 {
		return 0, errors.New(C.GoString(C.TTYerr))
	}
	return TTY(fd), nil
}

// TODO (proto): read errors never reported, even when, e.g., disconnecting usb to serial.
func (t TTY) Read(p []byte) (n int, err error) {
	n = int(C.readTTY(C.int(t), unsafe.Pointer(&p[0]), C.int(len(p))))
	if n <= 0 {
		return 0, errors.New(C.GoString(C.TTYerr))
	} else {
		return n, nil
	}
}

func (t TTY) Write(p []byte) (n int, err error) {
	n = int(C.writeTTY(C.int(t), unsafe.Pointer(&p[0]), C.int(len(p))))
	if n != len(p) {
		return n, errors.New(C.GoString(C.TTYerr))
	} else {
		return n, nil
	}
}

func (t TTY) ReadFull(p []byte) {
	N, err := t.Read(p)
	for N < len(p) {
		if err != nil {
			log.Fatal(err)
		}
		var n int
		n, err = t.Read(p[N:])
		N += n
	}
}

//// little endian
//func (t TTY) writeInt(I uint32) {
//	i := uint32(I)
//	bytes := []byte{
//		byte((i & 0x000000FF) >> 0),
//		byte((i & 0x0000FF00) >> 8),
//		byte((i & 0x00FF0000) >> 16),
//		byte((i & 0xFF000000) >> 24)}
//
//	_, err := t.Write(bytes)
//	log.Println("write", bytes)
//	if err != nil {
//		log.Fatal(err)
//	}
//}
//
//func (t TTY) readInt() int {
//	var bytes [4]byte
//	t.ReadFull(bytes[:])
//	return (int(bytes[0]) << 0) | (int(bytes[1]) << 8) | (int(bytes[2]) << 16) | (int(bytes[3]) << 24)
//}
