package tcpserver

import (
	"log"
	"net"
)

type ClientSession struct {
	ID       string   //客户端连接唯一标识
	coon     net.Conn //socket连接句柄
	IsClosed bool     // 标记会话是否关闭
}

func (session *ClientSession) Send(data []byte) {
	coon := session.coon
	coon.Write(data)
}

func (session *ClientSession) Close(reason string) {
	log.Println("client socket handler closed, reason: ", reason)
	session.IsClosed = true
	err := session.coon.Close()
	if err != nil {
		log.Println("client socket handler close err: ", err)
	}
}
