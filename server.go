package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net"
)

var listen string
var connect string
var connectSSH string

var ln net.Listener
var conn *net.TCPAddr
var connSSH *net.TCPAddr

func setup() {
	flag.StringVar(&listen, "listen", ":8000", "listen on address")
	flag.StringVar(&connect, "connect", "", "forward to address")
	flag.StringVar(&connectSSH, "ssh", "", "if set, will do basic introspection to forward SSH traffic to this address")
	flag.Parse()

	var err error

	// check and parse address
	conn, err = net.ResolveTCPAddr("tcp", connect)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	// check and parse SSH address
	connSSH, _ = net.ResolveTCPAddr("tcp", connectSSH)
	if connectSSH == "" {
		connSSH = nil
	}

	// listen on address
	ln, err = net.Listen("tcp", listen)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	log.Printf("listening on %v", ln.Addr())
	log.Printf("will connect to %v", conn)
	if connSSH != nil {
		log.Printf("will connect SSH to %v", connSSH)
	}
}

func serve() {
	for i := 0; ; i++ {
		// accept new connection
		c, err := ln.Accept()
		if err != nil {
			log.Print(err)
			break
		}

		log.Printf("connection %v from %v", i, c.RemoteAddr())
		go handle(c, i)
	}
}

var magic = []byte{'S', 'S', 'H', '-'}

var magicLen = len(magic)

func handle(c net.Conn, count int) {
	// read first four characters
	readMagic := make([]byte, magicLen, magicLen)
	n, err := c.Read(readMagic)
	if n != magicLen {
		log.Printf("warning! could not read header")
		return
	}
	opError, ok := err.(*net.OpError)
	if err != nil && (!ok || opError.Op != "readfrom") {
		log.Printf("warning! %v", err)
		return
	}

	connTo := conn
	// if the header looks like SSH, forward to SSH connection
	if bytes.Equal(readMagic, magic) {
		connTo = connSSH
	}

	cn, err := net.DialTCP("tcp", nil, connTo)
	if err != nil {
		c.Close()
		log.Print(err)
		return
	}

	// write the first four characters
	cn.Write(readMagic)

	go pipe(c, cn, count)
	go pipe(cn, c, count)
}

func pipe(w io.WriteCloser, r io.ReadCloser, count int) {
	n, err := io.Copy(w, r)

	r.Close()
	w.Close()

	log.Printf("connection %v closed, %v bytes", count, n)

	opError, ok := err.(*net.OpError)
	if err != nil && (!ok || opError.Op != "readfrom") {
		log.Printf("warning! %v", err)
	}
}
