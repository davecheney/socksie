package main

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"sync"
)

var connections = new(sync.WaitGroup)

func handleConn(local *net.TCPConn, dialer Dialer) {
	connections.Add(1)
	defer local.Close()
	defer connections.Done()

	in := bufio.NewReader(local)
	//out := local

	buf := make([]byte, 8)
	if _, err := io.ReadFull(in, buf); err != nil {
		log.Println("unable to read SOCKS request header", err)
		return
	}
	version := buf[0]
	command := buf[1]
	port := binary.BigEndian.Uint16(buf[2:4])
	ip := net.IP(buf[4:])
	addr := &net.TCPAddr{ip, int(port)}
	user, err := in.ReadString(0)
	if err != nil {
		log.Println("unable to read remote SOCKS user", err)
		return
	}
	log.Printf("incoming SOCKS request, version=%d, command=%v, raddr=%v, user=%s\n", version, command, addr, user)

	if version != 4 {
		log.Println("unknown version, closing connection")
		return
	}

	if command != 1 {
		log.Println("unsupported command, closing connection")
	}

	remote, err := dialer.DialTCP("tcp", local.RemoteAddr().(*net.TCPAddr), addr)
	if err != nil {
		log.Println("unable to connect to remote host", err)
		// TODO(dfc) send proper disconnection
		return
	}
	log.Println("connection succeeded to", remote.RemoteAddr())
	local.Write([]byte{0, 0x5a, 0, 0, 0, 0, 0, 0})
	go func() {
		n, err := io.Copy(remote, local)
		log.Printf("local to remote done, sent=%d, err=%v\n", n, err)
	}()
	n, err := io.Copy(local, remote)
	log.Printf("remote to local done, sent=%d, err=%v\n", n, err)
}
