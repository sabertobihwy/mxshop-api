package utils

import (
	"fmt"
	"net"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	fmt.Printf("%v\n", *addr)
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	fmt.Printf("%+v\n", *l)
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
