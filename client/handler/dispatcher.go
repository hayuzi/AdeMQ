package handler

type handlerFunc func(cmd string, params ...string)

type Dispatcher struct {
	History  []string
	Handlers map[string]handlerFunc
}

var Dpt = &Dispatcher{
	History:  make([]string, 256),
	Handlers: make(map[string]handlerFunc),
}
