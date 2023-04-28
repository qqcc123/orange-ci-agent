package tcpserver

type Middlewares []ActionFunc

type MiddlewareBeforeAction interface {
	MiddlewareBeforeAction() Middlewares
}

type MiddlewareAfterAction interface {
	MiddlewareAfterAction() Middlewares
}
