package tcpserver

import "bufio"

type ResolveActionFunc func(token []byte) (actionName string, msg []byte, err error)

type MsgDecorder interface {
	SplitFunc() bufio.SplitFunc
	ResolveFunc() ResolveActionFunc
}

func (server *Server) setDecorder(decorder MsgDecorder) {
	server.resolveAction = decorder.ResolveFunc()
	server.spliteRules = decorder.SplitFunc()
}

// 定义数据包的拆分结构
// @Begin 数据包开始标记
// @End 数据包结束标记
type NetDataSplit struct {
	Begin []byte
	End   []byte
}

//func (d *NetDataSplit) spliteRules() bufio.SplitFunc {
//	beginLength := len(d.Begin)
//	endLength := len(d.End)
//	return func(data []byte, atEOF bool) (int, []byte, error) {
//		if atEOF {
//			return 0, nil, nil
//		}
//		start, end := 0, 0
//
//	}
//}
