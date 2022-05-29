package main

import (
	"fmt"
	log1 "log"
	"os"
	"plugin.test/main/cmd/logic/solib/sofunc"
	"plugin.test/main/iface"
	"reflect"
	"time"
	"unsafe"
)

var log = log1.New(os.Stderr, "", log1.LstdFlags)

func init() {
	log.SetFlags(log1.Lmicroseconds | log1.LstdFlags | log1.Lmsgprefix | log1.Lshortfile)
	log.SetPrefix(fmt.Sprintf("[so%v]", ModuleVersion))
	log.Println("So is loded")
}

//here should using `go build -buildmode=plugin|c-shared`
//func main() {
//	panic("golang.plugin not work here!!!")
//}

var (
	ModuleVersion = "1001"
	ModuleName    = "game1"
)

type Game struct {
	Version string
	Name    string
}

func Exec(e iface.IEngine, name string) string {
	var game *Game
	v := e.Get("game")
	if v == nil {
		game = &Game{}
		e.Set("game", game)
	} else {
		game = (*Game)(unsafe.Pointer(reflect.ValueOf(v).Pointer()))
	}
	log.Println("so Exec running...")
	ret := fmt.Sprintf(" hello {%s}, this is golang plugin test!, version={%s}, name={%s}, oldversion={%s}, oldName={%s}\n", name, ModuleVersion, ModuleName, game.Version, game.Name)
	game.Version = ModuleVersion
	game.Name = ModuleName

	return ret
}
func NewFactory() iface.IFactory {
	return &Factory{}
}

type Factory struct {
}

func (this *Factory) Exec(e iface.IEngine, name string) string {
	//defer func() {
	//	this.FuncAddNew()
	//}()
	return Exec(e, name)
}
func (this *Factory) FuncInit() {
	log.Println(">> start Func init")
	go func() {
		for i := 0; ; i++ {
			time.Sleep(time.Duration(i) * time.Second)
			sofunc.FuncA(i, ModuleVersion)
		}
	}()
}

func (this *Factory) GetVersion() string {
	return ModuleVersion
}

//func (this *Factory) FuncAddNew() {
//	log.Println("FuncAddNew")
//}
