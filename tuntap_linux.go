package xtuntap

import (
	"errors"
	"os"
	"syscall"
	"unsafe"
)

type ifreq struct {
	name  [syscall.IFNAMSIZ]byte
	flags uint16
	_pad  [24 - unsafe.Sizeof(uint16(0))]byte
}

func tuntapAlloc(name string, bTun bool) (*os.File, string, error) {
	if len(name) > syscall.IFNAMSIZ-1 {
		return nil, "", errors.New("device name too long")
	}

	f, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0600)
	if err != nil {
		return nil, "", err
	}

	fd := f.Fd()

	/* Flags: IFF_TUN   - TUN device (no Ethernet headers)
	 *        IFF_TAP   - TAP device
	 *
	 *        IFF_NO_PI - Do not provide packet information
	 */
	flags := syscall.IFF_NO_PI | 0x100 /*syscall.IFF_MULTI_QUEUE*/
	if bTun {
		flags |= syscall.IFF_TUN
	} else {
		flags |= syscall.IFF_TAP
	}
	ifr := ifreq{
		flags: uint16(flags),
	}

	if len(name) > 0 {
		/* if a device name was specified, put it in the structure; otherwise,
		 * the kernel will try to allocate the "next" device of the
		 * specified type */
		copy(ifr.name[:], []byte(name))
	}

	/* try to create the device */
	if _, _, e := syscall.Syscall6(syscall.SYS_IOCTL, fd, syscall.TUNSETIFF, uintptr(unsafe.Pointer(&ifr)), 0, 0, 0); e != 0 {
		f.Close()
		return nil, "", e
	}

	// return kernel allocated name
	return f, string(ifr.name[:]), nil
}
