package iface

type IEngine interface {
	Get(string) interface{}
	Set(string, interface{})
	Del(string)
}

type IFactory interface {
	Exec(e IEngine, name string) string
	FuncInit()
	GetVersion() string
}
