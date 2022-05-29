package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"plugin"
	"plugin.test/main/iface"
	"runtime/debug"
	"sync"
)

func init() {
	log.SetFlags(log.Lmicroseconds | log.LstdFlags | log.Lmsgprefix | log.Lshortfile)
	log.SetPrefix("[engine]")
}

var (
	engine *Engine
	FnExec func(iface.IEngine, string) string
)

type Engine struct {
	mutex sync.Mutex
	vars  map[string]interface{}
}

func (e *Engine) Set(key string, val interface{}) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.vars[key] = val
}

func (e *Engine) Get(key string) interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.vars[key]
}

func (e *Engine) Del(key string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.vars, key)
}

func handleLoad(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	names := req.Form["name"]
	if len(names) > 0 {
		if e := load(names[0], false); e != nil {
			io.WriteString(w, e.Error())
			return
		}
	} else {
		log.Println("filename error")
	}
	io.WriteString(w, "done")
}

func handleHello(w http.ResponseWriter, req *http.Request) {
	if FnExec == nil {
		io.WriteString(w, "not execute the plugin!!!!!")
		return
	}
	//str := FnExec(engine, "test")
	str := SoFactory.Exec(engine, "test")
	io.WriteString(w, str)
}

func main() {
	engine = &Engine{
		vars: make(map[string]interface{}),
	}

	if e := load("plugin1.so", true); e != nil {
		panic(e)
	}

	http.HandleFunc("/load", handleLoad)
	http.HandleFunc("/hello", handleHello)
	http.HandleFunc("/gc", func(writer http.ResponseWriter, request *http.Request) {
		debug.FreeOSMemory()
	})
	log.Fatal(http.ListenAndServe(":12345", nil))
}

var SoFactory iface.IFactory

func load(filename string, isFirst bool) error {
	defer func() {
		if SoFactory != nil {
			log.Printf("SoVersion=%v", SoFactory.GetVersion())
			SoFactory.FuncInit()
		}
	}()
	p, err := plugin.Open(filename)
	if err != nil {
		log.Println("open plugin err:", err, filename)
		return err
	}
	log.Printf("plugin loaded success, p=%+v", p)
	if isFirst || true {
		f, e := p.Lookup("NewFactory")
		if e != nil {
			log.Println("not found symbol FuncA", err)
			return err
		}
		if v, ok := f.(func() iface.IFactory); ok {
			SoFactory = v()
		} else {
			panic("NewFactory Not type of func() iface.IFactory")
		}
	}
	fn, err := p.Lookup("Exec")
	if err != nil {
		log.Println("not found symbol Exec", err)
		return err
	}
	var ok bool
	FnExec, ok = fn.(func(iface.IEngine, string) string)
	if !ok {
		log.Printf("loaded plugin success, but not the correct func:%T %v", fn, fn)
		FnExec(nil, "") //为了panic
		return fmt.Errorf("loaded plugin success, but not the correct func:%T", fn)
	}
	log.Println("loaded plugin successed! file=", filename)
	return nil
}
