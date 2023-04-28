package tcpserver

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type Server struct {
	ip                   string                  // 服务器ip
	port                 int                     //服务器端口
	sessions             *sessionPool            //会话池
	tlsConfig            *tls.Config             //tls配置
	actions              map[string][]ActionFunc //消息处理行为方法
	acceptGoCount        int                     //用于处理连接的协成个数
	onError              func(error)             //错误处理函数
	onNewSessionRegister func(*ClientSession)    //新客户端接入注册
	onSessionClosed      func(*ClientSession)    //客户端断开连接
	resolveAction        ResolveActionFunc       //解析请求action
	spliteRules          bufio.SplitFunc         //拆包规则定义
	middlewareAfter      Middlewares
	middlewareBefore     Middlewares
}

func NewServer(ip string, port int) *Server {
	return newTlsServer(ip, port, nil)
}

func newTlsServer(_ip string, _port int, config *tls.Config) *Server {
	return &Server{
		ip:        _ip,
		port:      _port,
		tlsConfig: config,
		sessions: &sessionPool{
			list: make(chan *sessionHandle, 100),
		},
		acceptGoCount: 1,
		actions:       make(map[string][]ActionFunc),
	}
}

func (server *Server) Start() {
	if server.spliteRules == nil {
		log.Println("use default split rules")
		server.spliteRules = bufio.ScanLines
	}

	if len(server.actions) == 0 {
		log.Println("no actions")
	}

	addr := fmt.Sprintf("%s:%d", server.ip, server.port)
	if server.ip != "" && server.ip != "localhost" {
		ipAddr := net.ParseIP(server.ip)
		if ipAddr == nil {
			log.Println("ip address is not vaild ", server.ip)
			return
		}

		if ipAddr.To4() == nil {
			addr = fmt.Sprintf("[%s]:%d", server.ip, server.port)
		}
	}

	var tcpListener net.Listener
	var err error

	if server.tlsConfig == nil {
		tcpListener, err = net.Listen("tcp", addr)
	} else {
		tcpListener, err = tls.Listen("tcp", addr, server.tlsConfig)
	}
	if err != nil {
		log.Println("server listen err: ", err)
		server.handleError(err)
		return
	}

	//程序返回之前关闭socket
	defer tcpListener.Close()

	go server.sessions.sessionPoolManager()

	var waitGroup sync.WaitGroup
	for i := 0; i < server.acceptGoCount; i++ {
		waitGroup.Add(1)
		go func(acceptIndex int) {
			defer waitGroup.Done()
			for {
				conn, err := tcpListener.Accept()
				if err != nil {
					log.Println("accept client err: ", err)
					continue
				}

				server.handleClient(conn)
			}
		}(i)
	}
	waitGroup.Wait()

}

func (server *Server) handleClient(conn net.Conn) {
	session := ClientSession{
		ID:   uuid.NewV4().String(),
		coon: conn,
	}

	clientAddr := session.coon.RemoteAddr()
	log.Println("client address: ", clientAddr)

	if server.onNewSessionRegister != nil {
		server.onNewSessionRegister(&session)
	}

	//添加会话到会话池
	server.sessions.addSession(&session)

	scanner := bufio.NewScanner(conn)
	scanner.Split(server.spliteRules)
	for scanner.Scan() {
		token := scanner.Bytes()
		if server.resolveAction != nil {
			actionName, resolvedToken, err := server.resolveAction(token)
			if err != nil {
				log.Println("parser action err: ", err)
				break
			} else {
				hookErr := server.hookAction(actionName, &session, resolvedToken)
				if hookErr != nil {
					server.handleError(hookErr)
				}
			}
		}
	}

	scannerErr := scanner.Err()
	if scannerErr != nil {
		server.handleError(scannerErr)
		server.closeSession(&session, scannerErr.Error())
		return
	}
	server.closeSession(&session, "EOF")
}

func (server *Server) closeSession(session *ClientSession, reason string) {
	//关闭client的socket句柄
	go session.Close(reason)

	//会话池中删除当前会话
	go server.sessions.deleteSession(session)
}

func (server *Server) handleError(err error) {
	if server.onError != nil {
		server.onError(err)
	}
}

func (server *Server) SetSessionRegister(sessionRegister func(*ClientSession)) {
	server.onNewSessionRegister = sessionRegister
}

func (server *Server) SetSessionClosed(sessionClosed func(*ClientSession)) {
	server.onSessionClosed = sessionClosed
}
