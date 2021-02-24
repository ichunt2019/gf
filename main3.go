package main

import (
	"fmt"
	"github.com/ichunt2019/gf/os/gcfg"
	"github.com/ichunt2019/gf/container/gmap"
)

func Instance(name ...string) *gcfg.Config {
	key := "config"
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	instances := gmap.NewStrAnyMap(true)
	supportedFileTypes := []string{"toml", "yaml", "json", "ini", "xml"}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		fmt.Println("8899")
		c := gcfg.New()
		c.SetPath("./config/dev/")
		for _, fileType := range supportedFileTypes {
			if file := fmt.Sprintf(`%s.%s`, key, fileType); c.Available(file) {
				c.SetFileName(file)
				break
			}
		}
		return c
	}).(*gcfg.Config)
}

func main() {

	fmt.Println(Instance("proxy").GetString("http.addr"))
	fmt.Println(Instance("proxy").GetString("http.max_header_bytes"))
	fmt.Println(Instance("config").GetString("database.default.0.host"))
	fmt.Println(Instance("config").GetArray("owner.name"))
	fmt.Println(Instance("config").GetString("viewpath"))
	fmt.Println(Instance().GetString("viewpath"))


}