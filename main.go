package main

import (
	_ "docker-ops-server/boot"
	_ "docker-ops-server/router"

	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
