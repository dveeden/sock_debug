package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func handleConn(c net.Conn, serverSocket string) {
	srv, err := net.Dial("unix", serverSocket)
	if err != nil {
		log.Fatal("Failed to connect backend: ", err)
	}
	defer srv.Close()
	defer c.Close()

	go connCopy(srv, c, "c<-s")

	connCopy(c, srv, "c->s")
}

func connCopy(c1 net.Conn, c2 net.Conn, tag string) {
	for {
		buf := &bytes.Buffer{}
		header := make([]byte, 4)
		lengthRaw := []byte{0x0, 0x0, 0x0, 0x0}
		_, err := c1.Read(header)
		if err != nil {
			if err != io.EOF {
				log.Print("Failed to read: ", err)
			}
			break
		}

		copy(lengthRaw, header[:3])
		lengthRaw = append(lengthRaw[:3], 0x0)
		length := binary.LittleEndian.Uint32(lengthRaw)

		buf.Write(header)

		pkt := make([]byte, length)
		_, err = c1.Read(pkt)
		if err != nil {
			if err != io.EOF {
				log.Fatal("Failed to read from client: ", err)
			}
			break
		}
		buf.Write(pkt)
		if header[3] == 0x0 {
			fmt.Printf("%15s ", commandName(pkt[0]))
		} else {
			fmt.Printf("%15s ", "")
		}
		fmt.Printf("%s : %q\n", tag, buf)
		_, err = c2.Write(buf.Bytes())
		if err != nil {
			log.Fatal("Failed to write to server: ", err)
		}
	}
}

func commandName(pkt byte) string {
	switch pkt {
	case 0x1:
		return "COM_QUIT"
	case 0x3:
		return "COM_QUERY"
	case 0x9:
		return "COM_STATISTICS"
	default:
		return fmt.Sprintf("Unknown, %x", pkt)
	}
}

func main() {
	var serverSocket string
	flag.StringVar(&serverSocket, "serverSocket", "/tmp/tidb.sock", "socket to connect to")
	flag.Parse()

	os.RemoveAll("/tmp/sock_debug.socket")
	frontSock, err := net.Listen("unix", "/tmp/sock_debug.socket")
	if err != nil {
		log.Fatal("Failed to start listening: ", err)
	}
	defer frontSock.Close()

	for {
		conn, err := frontSock.Accept()
		if err != nil {
			log.Fatal("Failed to accept connection: ", err)
		}
		go handleConn(conn, serverSocket)
	}
}
