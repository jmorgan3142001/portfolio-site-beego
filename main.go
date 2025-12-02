package main

import (
	_ "portfolio-site/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	beego.Run()
}

