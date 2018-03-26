package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var listen string
var connect string

var ln net.Listener
var conn *net.TCPAddr

func setup() {
	flag.StringVar(&listen, "listen", ":8000", "listen on ip and port")
	flag.StringVar(&connect, "connect", "", "forward to ip and port")
	flag.Parse()

	var err error

	// check and parse address
	conn, err = net.ResolveTCPAddr("tcp", connect)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	// listen on address
	ln, err = net.Listen("tcp", listen)
	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	log.Printf("listening on %v", ln.Addr())
	log.Printf("will connect to %v", conn)
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

		cn, err := net.DialTCP("tcp", nil, conn)
		if err != nil {
			c.Close()
			log.Print(err)
			continue
		}

		go pipe(c, cn, i)
		go pipe(cn, c, i)
	}
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
