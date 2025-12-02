package routers

import (
	"portfolio-site/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.PortfolioController{})
}
