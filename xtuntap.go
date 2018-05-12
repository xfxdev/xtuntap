package xtuntap

import (
	"io"
)

type Tun interface {
	io.ReadWriteCloser
}

type TunDevice struct {
	io.ReadWriteCloser
	name string
}

func NewTun(name string) (Tun, error) {
	f, err := tuntapAlloc(name, true)
	if err != nil {
		return nil, err
	}
	return &TunDevice{
		ReadWriteCloser: f,
		name:            name,
	}, nil
}

type DummyTun struct {
}

func (t *DummyTun) Read(p []byte) (n int, err error) {
	return len(p), nil
}
func (t *DummyTun) Write(p []byte) (n int, err error) {
	return len(p), nil
}
func (t *DummyTun) Close() error {
	return nil
}

func NewDummyTun(name string) (Tun, error) {
	return &DummyTun{}, nil
}
