package tcpserver

import (
	"errors"
	"strings"
)

type ActionFunc func(*ClientSession, []byte) ([]byte, error)

func (server *Server) hookAction(_funcName string, _session *ClientSession, _token []byte) error {
	funcName := strings.ToLower(_funcName)
	actions, exist := server.actions[funcName]
	if !exist {
		return errors.New("action not found")
	}

	var err error
	var token []byte
	if server.middlewareBefore != nil {
		for i, _ := range server.middlewareBefore {
			token, err = server.middlewareBefore[i](_session, _token)
			if err != nil {
				return err
			}
		}
	}

	for i, _ := range actions {
		token, err = actions[i](_session, _token)
		if err != nil {
			return err
		}
	}

	if server.middlewareAfter != nil {
		for i, _ := range server.middlewareAfter {
			token, err = server.middlewareAfter[i](_session, _token)
			if err != nil {
				return err
			}
		}
	}

	if token != nil {
		_session.Send(token)
	}

	return nil
}

func (server *Server) Action(path string, actionFunc ...ActionFunc) error {
	if path == "" {
		return errors.New("action path invaild")
	}

	if _, exist := server.actions[path]; exist {
		return errors.New("action already exist")
	}

	server.actions[path] = actionFunc
	return nil
}
