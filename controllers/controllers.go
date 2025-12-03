package controllers

import (
	"github.com/beego/beego/v2/server/web"
)

type PortfolioController struct {
	web.Controller
}

func (c *PortfolioController) Get() {
	// Title and Personal Info
	c.Data["Title"] = "Jake Morgan | Software Engineer"
	c.Data["Name"] = "Jake Morgan"
	c.Data["Career"] = "Software Engineer"

	// Social Links
	c.Data["GithubLink"] = "https://github.com/jmorgan3142001"
	c.Data["LinkedinLink"] = "https://www.linkedin.com/in/jmorgan3142001"
	c.Data["Email"] = "jmorgan3142001@gmail.com"

	// Define Layout and Content
	c.Layout = "layout.html"
	c.TplName = "index.html"
}