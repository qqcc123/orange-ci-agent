package tcpserver

import (
	"log"
	"sync"
)

// 会话管理池
type sessionPool struct {
	pool    sync.Map            //会话池
	list    chan *sessionHandle //注册会话通道
	counter int                 //计数器
}

type sessionHandle struct {
	session   *ClientSession
	isAddPool bool //是否添加到会话池
}

func (p *sessionPool) sessionPoolManager() {
	for {
		handle, ok := <-p.list

		if !ok {
			log.Println("session channel closed")
			return
		}

		if handle.isAddPool {
			p.pool.Store(handle.session.ID, handle.session)
			p.counter++
		} else {
			p.pool.Delete(handle.session.ID)
			p.counter--
		}
	}
}

func (p *sessionPool) addSession(session *ClientSession) {
	p.list <- &sessionHandle{
		session:   session,
		isAddPool: true,
	}
}

func (p *sessionPool) deleteSession(session *ClientSession) {
	p.list <- &sessionHandle{
		session:   session,
		isAddPool: false,
	}
}
