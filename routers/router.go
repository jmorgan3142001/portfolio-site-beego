package routers

import (
	"portfolio-site/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.PortfolioController{}, "get:Get")
    beego.Router("/about", &controllers.PortfolioController{}, "get:About")
    beego.Router("/directory", &controllers.PortfolioController{}, "get:Directory")
    beego.Router("/terminal", &controllers.PortfolioController{}, "get:Terminal")
    beego.Router("/challenges", &controllers.PortfolioController{}, "get:Challenge")
    beego.Router("/challenges/run", &controllers.PortfolioController{}, "post:RunCode")
    beego.Router("/logs/submit", &controllers.PortfolioController{}, "post:SubmitLog")
}