package tcpserver

func StartSocket() {
	tcpServer := NewServer("172.18.3.68", 90)
	tcpServer.Start()
}
