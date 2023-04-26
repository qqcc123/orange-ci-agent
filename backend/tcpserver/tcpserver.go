package tcpserver

import (
	"fmt"
	"log"
	"net"
)

func RunSocket() {
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("listen err: ", err)
		return
	}

	var connections []net.Conn
	defer func() {
		for _, conn := range connections {
			conn.Close()
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				log.Printf("accept temp err: %v", ne)
				continue
			}

			fmt.Println("accept err: ", err)
			return
		}

		go handleConnect(conn)
		connections = append(connections, conn)
		if len(connections)%100 == 0 {
			log.Printf("total number of connections: %v", len(connections))
		}
	}
}

func handleConnect(conn net.Conn) {
	defer conn.Close()
	var buf = make([]byte, 1024)

	num, err := conn.Read(buf)
	if err != nil {
		fmt.Println("rad err: ", err)
		fmt.Println("read num: ", num)
		return
	}

	fmt.Println("read byte: ", num)
	fmt.Println("read data is: ", string(buf))
}
