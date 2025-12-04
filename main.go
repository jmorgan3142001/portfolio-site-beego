package main

import (
	_ "portfolio-site/routers"

	"github.com/beego/beego/v2/server/web"
)

func main() {
	web.SetStaticPath("/static", "static")
	web.Run()
}

