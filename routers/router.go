package routers

import (
	"portfolio-site/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
    beego.Router("/", &controllers.PortfolioController{}, "get:Get")
    beego.Router("/about", &controllers.PortfolioController{}, "get:About")
    beego.Router("/library", &controllers.PortfolioController{}, "get:Library")
    beego.Router("/research", &controllers.PortfolioController{}, "get:Research")
	beego.Router("/terminal", &controllers.PortfolioController{}, "get:Terminal")
    beego.Router("/logs/submit", &controllers.PortfolioController{}, "post:SubmitLog")
}